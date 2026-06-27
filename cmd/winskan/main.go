package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/imvalerio/winskan/pkg/registers/mru"
	"github.com/imvalerio/winskan/pkg/registers/system"
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

func printUSBHistory(w io.Writer) {
	fmt.Fprintf(w, "\n>>> CATEGORY: USB History (SYSTEM Hive) <<<\n")
	fmt.Fprintln(w, strings.Repeat("=", 60))

	devices, err := system.GetUSBHistory()
	if err != nil {
		fmt.Fprintf(w, "[-] Failed to parse USB history from registry: %v\n", err)
		return
	}

	if len(devices) == 0 {
		fmt.Fprintf(w, "No USB history found on this system.\n")
		return
	}

	for _, dev := range devices {
		fmt.Fprintf(w, "Device Serial: %s\n", dev.SerialNumber)
		if dev.VolumeName != "" {
			fmt.Fprintf(w, "  Volume Name:   %s\n", dev.VolumeName)
		}
		fmt.Fprintf(w, "  Friendly Name: %s\n", dev.FriendlyName)
		fmt.Fprintf(w, "  Vendor:        %s\n", dev.Vendor)
		fmt.Fprintf(w, "  Product:       %s\n", dev.Product)
		fmt.Fprintf(w, "  Revision:      %s\n", dev.Revision)
		if len(dev.HardwareID) > 0 {
			fmt.Fprintf(w, "  Hardware ID:   %s\n", strings.Join(dev.HardwareID, ", "))
		}
		fmt.Fprintln(w, strings.Repeat("-", 50))
	}
}

func printRunMRU(w io.Writer) {
	fmt.Fprintf(w, "\n>>> CATEGORY: RunMRU (Win + R) <<<\n")
	fmt.Fprintln(w, strings.Repeat("=", 60))

	entries, err := mru.ParseRunMRU()
	if err != nil {
		fmt.Fprintf(w, "[-] Failed to parse RunMRU registry key: %v\n", err)
		return
	}

	if len(entries) == 0 {
		fmt.Fprintf(w, "No RunMRU history found.\n")
		return
	}

	for _, entry := range entries {
		fmt.Fprintf(w, "[%d] %s\n", entry.Order, entry.Data)
	}
	fmt.Fprintln(w, strings.Repeat("-", 60))
}

func printRecentDocs(w io.Writer) {
	fmt.Fprintf(w, "\n>>> CATEGORY: RecentDocs <<<\n")
	fmt.Fprintln(w, strings.Repeat("=", 60))

	entries, err := mru.ParseRecentDocs()
	if err != nil {
		fmt.Fprintf(w, "[-] Failed to parse RecentDocs registry key: %v\n", err)
		return
	}

	if len(entries) == 0 {
		fmt.Fprintf(w, "No RecentDocs history found.\n")
		return
	}

	for _, entry := range entries {
		fmt.Fprintf(w, "[%d] %s\n", entry.Order, entry.Data)
	}
	fmt.Fprintln(w, strings.Repeat("-", 60))
}

func getFilteredEntries(guid string, lastRunOnly, focusTimeOnly bool) []userassist.Entry {
	entries, err := userassist.Parse(guid)
	if err != nil {
		log.Printf("[-] Failed to parse registry key for %s: %v\n", guid, err)
		return nil
	}
	var filtered []userassist.Entry
	for _, entry := range entries {
		if lastRunOnly && entry.LastRun.IsZero() {
			continue
		}
		if focusTimeOnly && entry.FocusTimeMs <= 0 {
			continue
		}
		filtered = append(filtered, entry)
	}
	return filtered
}

func main() {
	lastRunOnly := flag.Bool("last-run", false, "Print only entries with a non-zero last execution time")
	focusTimeOnly := flag.Bool("focus-time", false, "Print only entries with a non-zero focus time")
	category := flag.String("category", "all", "Select category to display: 'exe', 'lnk', 'usb', 'runmru', 'recentdocs', or 'all'")
	outputFile := flag.String("o", "", "Write output to the specified .txt file")
	guiOutput := flag.Bool("gui", false, "Generate an HTML report with a minimal and polished design")
	flag.Parse()

	if *outputFile == "" {
		*outputFile = "output.txt"
	} else {
		if index := strings.LastIndex(*outputFile, "."); index != -1 {
			*outputFile = (*outputFile)[:index] + ".txt"
		} else {
			*outputFile += ".txt"
		}
	}

	guiOutputFile := *outputFile

	if *guiOutput {
		index := strings.LastIndex(guiOutputFile, ".")
		guiOutputFile = (guiOutputFile)[:index] + ".html"

		f, err := os.Create(guiOutputFile)
		if err != nil {
			log.Fatalf("Failed to create HTML output file: %v", err)
		}
		defer f.Close()

		data := ReportData{
			Timestamp: time.Now(),
		}

		cat := strings.ToLower(*category)

		if cat == "all" || cat == "exe" {
			data.Executables = getFilteredEntries(userassist.GUIDExecutables, *lastRunOnly, *focusTimeOnly)
		}
		if cat == "all" || cat == "lnk" {
			data.Shortcuts = getFilteredEntries(userassist.GUIDShortcuts, *lastRunOnly, *focusTimeOnly)
		}
		if cat == "all" || cat == "usb" {
			devices, err := system.GetUSBHistory()
			if err == nil {
				data.USBDevices = devices
			} else {
				log.Printf("[-] Failed to parse USB history: %v\n", err)
			}
		}
		if cat == "all" || cat == "runmru" {
			entries, err := mru.ParseRunMRU()
			if err == nil {
				data.RunMRU = entries
			} else {
				log.Printf("[-] Failed to parse RunMRU: %v\n", err)
			}
		}
		if cat == "all" || cat == "recentdocs" {
			entries, err := mru.ParseRecentDocs()
			if err == nil {
				data.RecentDocs = entries
			} else {
				log.Printf("[-] Failed to parse RecentDocs: %v\n", err)
			}
		}

		err = generateHTMLReport(f, data)
		if err != nil {
			log.Fatalf("Failed to generate HTML report: %v", err)
		}
		fmt.Printf("[+] HTML report generated successfully at: %s\n", *outputFile)
	}

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

	if cat == "all" || cat == "usb" {
		printUSBHistory(out)
	}

	if cat == "all" || cat == "runmru" {
		printRunMRU(out)
	}

	if cat == "all" || cat == "recentdocs" {
		printRecentDocs(out)
	}
}
