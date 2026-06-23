package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/imvalerio/winskan/pkg/registers/userassist"
)

func printEntries(w io.Writer, categoryName, guid string, lastRunOnly, focusTimeOnly bool) {
	fmt.Fprintf(w, "\n>>> CATEGORY: %s (%s) <<<\n", categoryName, guid)
	fmt.Fprintln(w, strings.Repeat("=", 60))

	entries, err := userassist.Parse(guid)
	if err != nil {
		fmt.Fprintf(w, "[-] Failed to parse registry key for %s. It might not exist on this system or there is an error: %v\n", categoryName, err)
		return
	}

	for _, entry := range entries {
		if lastRunOnly && entry.LastRun.IsZero() {
			continue
		}
		if focusTimeOnly && entry.FocusTimeMs <= 0 {
			continue
		}

		fmt.Fprintf(w, "Entry: %s\n", entry.Name)
		fmt.Fprintf(w, "  Run Count:   %d\n", entry.RunCount)
		fmt.Fprintf(w, "  Focus Count: %d\n", entry.FocusCount)
		fmt.Fprintf(w, "  Focus Time:  %d ms (%.2f seconds)\n", entry.FocusTimeMs, float64(entry.FocusTimeMs)/1000.0)

		if !entry.LastRun.IsZero() {
			fmt.Fprintf(w, "  Last Run:    %s\n", entry.LastRun.Format(time.RFC1123))
		} else {
			fmt.Fprintf(w, "  Last Run:    Never\n")
		}
		fmt.Fprintln(w, strings.Repeat("-", 50))
	}
}

func main() {
	lastRunOnly := flag.Bool("last-run", false, "Print only entries with a non-zero last execution time")
	focusTimeOnly := flag.Bool("focus-time", false, "Print only entries with a non-zero focus time")
	category := flag.String("category", "all", "Select category to display: 'exe', 'lnk', or 'all'")
	outputFile := flag.String("o", "", "Write output to the specified .txt file")
	flag.Parse()

	var out io.Writer = os.Stdout

	if *outputFile != "" {
		f, err := os.Create(*outputFile)
		if err != nil {
			log.Fatalf("Failed to create output file: %v", err)
		}
		defer f.Close()
		out = io.MultiWriter(os.Stdout, f)
	}

	cat := strings.ToLower(*category)

	if cat == "all" || cat == "exe" {
		printEntries(out, "Executables", userassist.GUIDExecutables, *lastRunOnly, *focusTimeOnly)
	}

	if cat == "all" || cat == "lnk" {
		printEntries(out, "Shortcuts", userassist.GUIDShortcuts, *lastRunOnly, *focusTimeOnly)
	}
}
