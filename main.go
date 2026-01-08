package main

import (
	"fmt"
	"log"
	"log/slog"
	"strings"

	"github.com/pterm/pterm"
)

const (
	allKey  = "all"
	quitKey = "q"
)

func main() {
	flags := ParseFlags()
	if flags.Version {
		fmt.Println(FormatVersion())
		return
	}

	if flags.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
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

	var quit bool
	for !quit {
		renderFolders(result.Folders)

		selection, _ := pterm.DefaultInteractiveTextInput.Show("Enter a range (e.g. 1-3), a list (e.g. 1,3,5), 'all', or 'q' to quit")
		selection = strings.Trim(selection, "")

		if selection == "" || selection == quitKey {
			quit = true
		}

		if selection == allKey {
			confirm, err := pterm.DefaultInteractiveConfirm.Show()
			if err != nil || !confirm {
				pterm.Warning.Println("Operation cancelled")
				return
			}

			totalCount := result.PendingCount()
			sp, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Deleting 0/%d folders", totalCount))

			deletedCount := 0
			for i := 0; i <= len(result.Folders)-1; i++ {
				if result.Folders[i].Deleted {
					continue
				}

				sp.UpdateText(fmt.Sprintf("Deleting %d/%d folders", deletedCount+1, totalCount))

				if err := result.Folders[i].Delete(); err != nil {
					pterm.Error.Printf("Failed to delete %s\n", result.Folders[i].Path)
				} else {
					deletedCount++
				}
			}

			sp.Success(fmt.Sprintf("Completed! %d/%d folders deleted", deletedCount, totalCount))
		}

		pterm.Println()
	}
}

func renderFolders(folders []Folder) {
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
