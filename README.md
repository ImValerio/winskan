# Winskan

Winskan is a Go-based tool designed to parse Windows **UserAssist** registry keys, **USB Device History**, and **MRU Lists**. It extracts and displays information about applications executed on a Windows system, including their run counts, focus times, and the last time they were executed, alongside a detailed history of USB mass storage devices connected to the system, and chronological lists of recently executed commands and accessed documents.

## Features

- **Registry Parsing**: Reads the `UserAssist` registry keys for both Executables and Shortcuts (.lnk).
- **USB Device History**: Enumerates the `SYSTEM` hive to retrieve historical connections of USB storage devices.
- **Volume Label Extraction**: Correlates device serial numbers with the Windows Portable Devices registry key to extract custom user-assigned Volume Names (e.g., "ESD-USB").
- **MRU Parsing**: Decodes and chronologically sorts Most Recently Used artifacts like `RunMRU` (commands typed in the Run dialog) and `RecentDocs` (recently opened files).
- **ROT13 Decoding**: Automatically decodes application paths which are stored using ROT13 encoding in the registry.
- **Binary Data Extraction**: Parses the 72-byte binary structures used by modern Windows (7, 8, 10, and 11) to store execution metadata.
- **Human-Readable Output**: Converts Windows FILETIME values to standard timestamps and displays execution statistics in a clear format.
- **Filtering**: Supports filtering entries by execution time, focus time, or category.

## Flags

| Flag           | Description                                                                  |
| :------------- | :--------------------------------------------------------------------------- |
| `--last-run`   | Shows only entries that have been executed (non-zero last run time).         |
| `--focus-time` | Shows only entries that have recorded focus time (non-zero focus time).      |
| `--category`   | Select category: `exe` (Executables), `lnk` (Shortcuts), `usb` (USB History), `runmru` (Run Dialog), `recentdocs` (Recent Documents), or `all` (default). |
| `-o`           | Writes the output to the specified file.                                     |

## How it Works

### UserAssist

The tool targets the following registry paths under `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Explorer\UserAssist\`:

- `{CEBFF5CD-ACE2-4F4F-9178-9926F41749EA}`: Tracks direct `.exe` executions.
- `{F4E57C4B-2036-45F0-A9AB-443BCFE33D9F}`: Tracks launches via `.lnk` files (shortcuts).

### USB History

The tool correlates data from two separate registry locations to build a comprehensive view of USB devices:
- `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Enum\USBSTOR`: Retrieves connected USB instances, vendor/product details, hardware IDs, and the driver-assigned "Friendly Name".
- `HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows Portable Devices\Devices`: Maps the device's serial number to its custom Volume Name (the label displayed in Windows Explorer).

### MRU Lists

The tool extracts and decodes Most Recently Used items from:
- `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Explorer\RunMRU`: Extracts commands typed into the Windows Run dialog (Win + R).
- `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Explorer\RecentDocs`: Extracts file paths of recently accessed documents and shortcuts, parsing the binary payloads and chronological order markers.

## Extracted Information

### UserAssist Entries

For each entry, Winskan extracts:

- **Entry Name**: The decoded path or name of the executable or shortcut.
- **Run Count**: The total number of times the entry has been launched.
- **Focus Count**: The number of times the window gained focus.
- **Focus Time**: The total time (in milliseconds) the application was in the foreground.
- **Last Run**: The timestamp of the most recent execution.

### USB Device Entries

For each device, Winskan extracts:

- **Device Serial**: The unique instance or serial number of the USB device.
- **Volume Name**: The user-assigned label (e.g., "My USB").
- **Friendly Name**: The hardware-level name (e.g., "SanDisk Cruzer USB Device").
- **Vendor / Product / Revision**: Extracted from the device class ID.
- **Hardware ID**: The full hardware identifier list.

### MRU Entries

For each MRU item, Winskan extracts:

- **Order**: The chronological position (1 being the absolute most recently interacted with).
- **Data**: The parsed string representing the command run or the file path accessed.

## Requirements

- **Operating System**: Windows (required to access the registry).
- **Go**: A working Go environment. (Optional, you can compile the source)

## Usage

```bash
go run main.go
```
