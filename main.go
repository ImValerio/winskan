package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/sys/windows/registry"
)

const (
	guidExecutables = "{CEBFF5CD-ACE2-4F4F-9178-9926F41749EA}"
	guidShortcuts   = "{F4E57C4B-2036-45F0-A9AB-443BCFE33D9F}"
)

// rot13 decodes the ROT13 encoded registry value names
func rot13(s string) string {
	var result strings.Builder
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			result.WriteRune('a' + (r-'a'+13)%26)
		} else if r >= 'A' && r <= 'Z' {
			result.WriteRune('A' + (r-'A'+13)%26)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// filetimeToTime converts a Windows FILETIME uint64 to a standard Go time.Time
func filetimeToTime(ft uint64) time.Time {
	if ft == 0 {
		return time.Time{}
	}
	const unixEpochOffset = 116444736000000000
	nsec := int64(ft-unixEpochOffset) * 100
	return time.Unix(0, nsec)
}

func parseUserAssist(w io.Writer, guid string, lastRunOnly, focusTimeOnly bool) {
	categoryName := "Executables"
	if guid == guidShortcuts {
		categoryName = "Shortcuts"
	}

	fmt.Fprintf(w, "\n>>> CATEGORY: %s (%s) <<<\n", categoryName, guid)
	fmt.Fprintln(w, strings.Repeat("=", 60))

	basePath := fmt.Sprintf(`Software\Microsoft\Windows\CurrentVersion\Explorer\UserAssist\%s\Count`, guid)

	key, err := registry.OpenKey(registry.CURRENT_USER, basePath, registry.QUERY_VALUE)
	if err != nil {
		fmt.Fprintf(w, "[-] Failed to open registry key for %s. It might not exist on this system.\n", categoryName)
		return
	}
	defer key.Close()

	valNames, err := key.ReadValueNames(-1)
	if err != nil {
		log.Printf("Failed to read value names for %s: %v", categoryName, err)
		return
	}

	for _, valName := range valNames {
		decodedName := rot13(valName)
		if strings.HasPrefix(decodedName, "UEME_") {
			continue
		}

		data, _, err := key.GetBinaryValue(valName)
		if err != nil {
			continue
		}

		if len(data) == 72 {
			runCount := binary.LittleEndian.Uint32(data[4:8])
			focusCount := binary.LittleEndian.Uint32(data[8:12])
			focusTimeMs := binary.LittleEndian.Uint32(data[12:16])
			lastExecutionRaw := binary.LittleEndian.Uint64(data[60:68])
			lastExecution := filetimeToTime(lastExecutionRaw)

			if lastRunOnly && lastExecution.IsZero() {
				continue
			}
			if focusTimeOnly && focusTimeMs <= 0 {
				continue
			}

			fmt.Fprintf(w, "Entry: %s\n", decodedName)
			fmt.Fprintf(w, "  Run Count:   %d\n", runCount)
			fmt.Fprintf(w, "  Focus Count: %d\n", focusCount)
			fmt.Fprintf(w, "  Focus Time:  %d ms (%.2f seconds)\n", focusTimeMs, float64(focusTimeMs)/1000.0)

			if !lastExecution.IsZero() {
				fmt.Fprintf(w, "  Last Run:    %s\n", lastExecution.Format(time.RFC1123))
			} else {
				fmt.Fprintf(w, "  Last Run:    Never\n")
			}
			fmt.Fprintln(w, strings.Repeat("-", 50))
		}
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
		parseUserAssist(out, guidExecutables, *lastRunOnly, *focusTimeOnly)
	}

	if cat == "all" || cat == "lnk" {
		parseUserAssist(out, guidShortcuts, *lastRunOnly, *focusTimeOnly)
	}
}
