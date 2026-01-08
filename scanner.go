package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type Folder struct {
	Path    string
	Size    int64
	ModTime time.Time
	Deleted bool
}

func (f *Folder) Delete() error {
	if err := os.RemoveAll(f.Path); err != nil {
		return fmt.Errorf("cannot delete %s: %w", f.Path, err)
	}

	f.Deleted = true
	return nil
}

type SkippedFolder struct {
	Path   string
	Reason error
}

func (s SkippedFolder) Error() string {
	return fmt.Sprintf("cannot read info for %s: %s", s.Path, s.Reason.Error())
}

type NodeModulesScanResult struct {
	Folders        []Folder
	Errors         []error
	SkippedFolders []SkippedFolder
}

func (n *NodeModulesScanResult) TotalSize() (total int64) {
	for _, f := range n.Folders {
		total += f.Size
	}

	return
}

func (n *NodeModulesScanResult) FoldersCount() int {
	return len(n.Folders)
}

func (n *NodeModulesScanResult) PendingCount() (total int) {
	for _, f := range n.Folders {
		if !f.Deleted {
			total++
		}
	}

	return
}

func findNodeModules(rootDir string) (NodeModulesScanResult, error) {
	result := NodeModulesScanResult{}

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			result.Errors = append(result.Errors, err)
			return nil
		}

		if d.IsDir() && d.Name() == "node_modules" {
			fileInfo, err := d.Info()
			if err != nil {
				result.SkippedFolders = append(result.SkippedFolders, SkippedFolder{Path: path, Reason: err})
				return filepath.SkipDir
			}

			folder := Folder{
				Path:    path,
				Size:    0,
				ModTime: fileInfo.ModTime(),
				Deleted: false,
			}

			result.Folders = append(result.Folders, folder)
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return NodeModulesScanResult{}, err
	}

	return result, nil
}

func calculateDirSize(path string) (int64, error) {
	var size int64

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		fileInfo, err := d.Info()
		if err != nil {
			return err
		}

		size += fileInfo.Size()
		return nil
	})

	if err != nil {
		return 0, err
	}

	return size, nil
}

type sizeResult struct {
	index int
	size  int64
	err   error
}

func measureSizesConcurrently(scanResult *NodeModulesScanResult) {
	c := make(chan sizeResult, len(scanResult.Folders))

	for i := range scanResult.Folders {
		go func(idx int) {
			size, err := calculateDirSize(scanResult.Folders[idx].Path)
			if err != nil {
				c <- sizeResult{index: idx, size: -1, err: err}
			} else {
				c <- sizeResult{index: idx, size: size, err: nil}
			}
		}(i)
	}

	for range len(scanResult.Folders) {
		data := <-c
		if data.err == nil {
			scanResult.Folders[data.index].Size = data.size
		}
	}
}

func FindAndMeasureNodeModules(root string) (NodeModulesScanResult, error) {
	result, err := findNodeModules(root)
	if err != nil {
		return NodeModulesScanResult{}, err
	}

	measureSizesConcurrently(&result)
	return result, nil
}
