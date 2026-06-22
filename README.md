# Winskan

Winskan is a Go-based tool designed to parse Windows **UserAssist** registry keys. It extracts and displays information about applications executed on a Windows system, including their run counts, focus times, and the last time they were executed.

## Features

- **Registry Parsing**: Reads the `UserAssist` registry keys for both Executables and Shortcuts (.lnk).
- **ROT13 Decoding**: Automatically decodes application paths which are stored using ROT13 encoding in the registry.
- **Binary Data Extraction**: Parses the 72-byte binary structures used by modern Windows (7, 8, 10, and 11) to store execution metadata.
- **Human-Readable Output**: Converts Windows FILETIME values to standard timestamps and displays execution statistics in a clear format.
- **Filtering**: Supports filtering entries by execution time, focus time, or category.

## Flags

| Flag           | Description                                                                  |
| :------------- | :--------------------------------------------------------------------------- |
| `--last-run`   | Shows only entries that have been executed (non-zero last run time).         |
| `--focus-time` | Shows only entries that have recorded focus time (non-zero focus time).      |
| `--category`   | Select category: `exe` (Executables), `lnk` (Shortcuts), or `all` (default). |
| `-o`           | Writes the output to the specified file.                                     |

## How it Works

The tool targets the following registry paths under `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Explorer\UserAssist\`:

- `{CEBFF5CD-ACE2-4F4F-9178-9926F41749EA}`: Tracks direct `.exe` executions.
- `{F4E57C4B-2036-45F0-A9AB-443BCFE33D9F}`: Tracks launches via `.lnk` files (shortcuts).

### Extracted Information

For each entry, Winskan extracts:

- **Entry Name**: The decoded path or name of the executable or shortcut.
- **Run Count**: The total number of times the entry has been launched.
- **Focus Count**: The number of times the window gained focus.
- **Focus Time**: The total time (in milliseconds) the application was in the foreground.
- **Last Run**: The timestamp of the most recent execution.

## Requirements

- **Operating System**: Windows (required to access the registry).
- **Go**: A working Go environment. (Optional, you can compile the source)

## Usage

```bash
go run main.go
```
