package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/pterm/pterm"
)

func main() {
	flags := ParseFlags()
	if flags.Version {
		fmt.Println(FormatVersion())
		return
	}

	if flags.Help {
		flags.Usage()
		return
	}

	pterm.Info.Println(
		fmt.Sprintf("Scanning %s for node_modules...", pterm.Magenta(flags.Directory)),
	)

	sp, err := pterm.DefaultSpinner.Start("Searching for node_modules...")
	if err != nil {
		log.Fatal(err)
	}

	result, err := FindAndMeasureNodeModules(flags.Directory)
	if err != nil {
		sp.Fail("Scan failed")
		return
	}

	if result.FoldersCount() == 0 {
		sp.Info("No node_modules folders found")
		return
	}

	msg := fmt.Sprintf("Found %d node_modules folders (%s total)\n", result.FoldersCount(), byteCountIEC(result.TotalSize()))
	sp.Success(msg)

	renderFolders(result.Folders)

	var quit bool
	for !quit {
		selection, _ := pterm.DefaultInteractiveTextInput.Show("Enter a number, a list (e.g. 1,3,5), 'all' or 'q' to quit")
		selection = strings.TrimSpace(selection)

		switch strings.ToLower(selection) {
		case "q", "quit":
			quit = true
		case "all":
			var indices []int
			for i, f := range result.Folders {
				if !f.Deleted {
					indices = append(indices, i)
				}
			}

			deleteWithProgress(len(indices), func(i int) error {
				return result.Folders[indices[i]].Delete()
			})
			renderFolders(result.Folders)
		default:
			var indices []int
			var parseErr error

			if strings.Contains(selection, ",") {
				indices, parseErr = parseList(selection, result.FoldersCount())
			} else {
				num, err := strconv.Atoi(selection)
				if err != nil {
					parseErr = errors.New("not a number")
				} else if num < 1 || num > len(result.Folders) {
					parseErr = fmt.Errorf("folder %d not found (valid: 1-%d)", num, len(result.Folders))
				} else {
					if result.Folders[num-1].Deleted {
						renderFolders(result.Folders)
						continue
					}

					indices = append(indices, num-1)
				}
			}

			if parseErr != nil {
				pterm.Error.Print(parseErr.Error())
				pterm.Println()
				continue
			}

			deleteWithProgress(len(indices), func(i int) error {
				return result.Folders[indices[i]].Delete()
			})
			renderFolders(result.Folders)
		}
	}
}

func renderFolders(folders []Folder) {
	pterm.Println()

	for i, f := range folders {
		var prefix string
		if f.Deleted {
			prefix = pterm.Red("[deleted]")
		} else {
			prefix = pterm.Blue(fmt.Sprintf("[%d]", i+1))
		}

		size := byteCountIEC(f.Size)
		age := relativeTime(f.ModTime)

		pterm.Printf("%s %s\n", prefix, f.Path)

		if age == "" {
			pterm.Print(pterm.Gray(fmt.Sprintf("    %s\n", size)))
		} else {
			pterm.Print(pterm.Gray(fmt.Sprintf("    %s • %s\n", size, age)))
		}

		pterm.Println()
	}
}

func parseList(input string, max int) ([]int, error) {
	parts := strings.SplitSeq(input, ",")

	var indices []int
	seen := make(map[int]bool)

	for part := range parts {
		part = strings.TrimSpace(part)

		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("not a number: %s", part)
		}

		if num < 0 || num > max {
			return nil, fmt.Errorf("folder %d not found (valid: 1-%d)", num, max)
		}

		if !seen[num] {
			indices = append(indices, num-1)
			seen[num] = true
		}
	}

	return indices, nil
}

type DeleteFolderFn func(int) error

func deleteWithProgress(total int, fn DeleteFolderFn) {
	if total == 0 {
		pterm.Info.Println("No folders to delete.")
		return
	}

	msg := fmt.Sprintf("Delete %d folder(s)?", total)
	if confirmed, _ := pterm.DefaultInteractiveConfirm.Show(msg); !confirmed {
		pterm.Warning.Println("Operation cancelled")
		return
	}

	sp, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Deleting 0/%d folders", total))
	deleted := 0

	for i := range total {
		sp.UpdateText(fmt.Sprintf("Deleting %d/%d folders", deleted+1, total))
		if err := fn(i); err != nil {
			pterm.Error.Printf("❌ %v\n", err)
		} else {
			deleted++
		}
	}

	sp.Success(fmt.Sprintf("Completed! %d/%d deleted\n", deleted, total))
}
