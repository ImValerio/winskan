package main

import (
	"html/template"
	"io"
	"time"

	"github.com/imvalerio/winskan/pkg/registers/mru"
	"github.com/imvalerio/winskan/pkg/registers/system"
	"github.com/imvalerio/winskan/pkg/registers/userassist"
)

type ReportData struct {
	Timestamp   time.Time
	Executables []userassist.Entry
	Shortcuts   []userassist.Entry
	USBDevices  []system.USBDevice
	RunMRU      []mru.Entry
	RecentDocs  []mru.Entry
}

const htmlTemplateStr = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Winskan Report</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Roboto:wght@300;400;500;700&display=swap" rel="stylesheet">
    <style>
        :root {
            color-scheme: light;
            --bg-color: #f5f5f5;
            --surface-color: #ffffff;
            --primary-color: #1a73e8; /* Google Blue */
            --primary-hover: #174ea6;
            --text-primary: rgba(0, 0, 0, 0.87);
            --text-secondary: rgba(0, 0, 0, 0.60);
            --divider-color: rgba(0, 0, 0, 0.12);
            --hover-color: rgba(0, 0, 0, 0.04);
            --elevation-1: 0 1px 3px rgba(0,0,0,0.12), 0 1px 2px rgba(0,0,0,0.24);
            --elevation-2: 0 3px 6px rgba(0,0,0,0.16), 0 3px 6px rgba(0,0,0,0.23);
        }

        :root[data-theme="dark"] {
            color-scheme: dark;
            --bg-color: #121212;
            --surface-color: #1e1e1e;
            --primary-color: #8ab4f8;
            --primary-hover: #aecbfa;
            --text-primary: rgba(255, 255, 255, 0.92);
            --text-secondary: rgba(255, 255, 255, 0.68);
            --divider-color: rgba(255, 255, 255, 0.16);
            --hover-color: rgba(255, 255, 255, 0.08);
            --elevation-1: 0 1px 3px rgba(0,0,0,0.48), 0 1px 2px rgba(0,0,0,0.72);
            --elevation-2: 0 3px 6px rgba(0,0,0,0.52), 0 3px 6px rgba(0,0,0,0.68);
        }

        body {
            margin: 0;
            padding: 0;
            background-color: var(--bg-color);
            color: var(--text-primary);
            font-family: 'Roboto', sans-serif;
            line-height: 1.5;
            -webkit-font-smoothing: antialiased;
            transition: background-color 0.2s ease, color 0.2s ease;
        }

        [hidden] {
            display: none !important;
        }

        header.app-bar {
            background-color: var(--primary-color);
            color: var(--surface-color);
            padding: 24px;
            box-shadow: var(--elevation-2);
            transition: background-color 0.2s ease, color 0.2s ease;
        }

        .app-bar-content {
            max-width: 1200px;
            margin: 0 auto;
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 24px;
        }

        .report-title {
            flex: 1;
            text-align: center;
        }

        header.app-bar h1 {
            margin: 0;
            font-size: 2.25rem;
            font-weight: 500;
            letter-spacing: 0.01em;
        }

        .timestamp {
            font-size: 0.9rem;
            opacity: 0.85;
            margin-top: 8px;
        }

        .theme-toggle {
            display: inline-flex;
            align-items: center;
            gap: 10px;
            user-select: none;
            white-space: nowrap;
            font-size: 0.9rem;
            font-weight: 500;
        }

        .theme-toggle input {
            position: absolute;
            opacity: 0;
            width: 1px;
            height: 1px;
        }

        .theme-toggle-slider {
            position: relative;
            width: 48px;
            height: 26px;
            border-radius: 999px;
            background-color: rgba(255, 255, 255, 0.35);
            cursor: pointer;
            transition: background-color 0.2s ease, box-shadow 0.2s ease;
        }

        .theme-toggle-slider::before {
            content: "";
            position: absolute;
            top: 3px;
            left: 3px;
            width: 20px;
            height: 20px;
            border-radius: 50%;
            background-color: var(--surface-color);
            box-shadow: 0 1px 3px rgba(0,0,0,0.35);
            transition: transform 0.2s ease, background-color 0.2s ease;
        }

        .theme-toggle input:checked + .theme-toggle-slider {
            background-color: rgba(0, 0, 0, 0.38);
        }

        .theme-toggle input:checked + .theme-toggle-slider::before {
            transform: translateX(22px);
        }

        .theme-toggle input:focus-visible + .theme-toggle-slider {
            box-shadow: 0 0 0 3px rgba(255, 255, 255, 0.5);
        }

        .theme-toggle-text {
            cursor: pointer;
        }

        .container {
            max-width: 1200px;
            margin: 32px auto;
            padding: 0 24px;
        }

        .data-section {
            background-color: var(--surface-color);
            border-radius: 8px;
            padding: 24px;
            margin-bottom: 32px;
            box-shadow: var(--elevation-1);
            overflow-x: auto;
            transition: background-color 0.2s ease, box-shadow 0.2s ease;
        }

        h2 {
            color: var(--text-primary);
            margin-top: 0;
            margin-bottom: 24px;
            font-size: 1.5rem;
            font-weight: 500;
            border-bottom: 1px solid var(--divider-color);
            padding-bottom: 16px;
        }

        /* Controls: Search and Pagination Page Size */
        .controls {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 16px;
            flex-wrap: wrap;
            gap: 16px;
        }

        .search-input, .page-size {
            padding: 12px 16px;
            border-radius: 4px;
            border: 1px solid var(--divider-color);
            background-color: var(--surface-color);
            color: var(--text-primary);
            font-family: inherit;
            font-size: 1rem;
            transition: border-color 0.2s, box-shadow 0.2s;
        }

        .search-input {
            flex: 1;
            max-width: 320px;
        }

        .search-input:focus, .page-size:focus {
            outline: none;
            border-color: var(--primary-color);
            box-shadow: 0 0 0 1px var(--primary-color);
        }

        table {
            width: 100%;
            border-collapse: collapse;
            text-align: left;
        }

        th, td {
            padding: 12px 16px;
            border-bottom: 1px solid var(--divider-color);
        }

        th {
            font-weight: 500;
            color: var(--text-secondary);
            font-size: 0.875rem;
        }

        td {
            font-size: 0.875rem;
            color: var(--text-primary);
        }

        tr:last-child td {
            border-bottom: none;
        }

        tbody tr {
            transition: background-color 0.2s ease;
        }

        tbody tr:hover {
            background-color: var(--hover-color);
        }

        /* Pagination Buttons */
        .pagination {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-top: 16px;
            border-top: 1px solid var(--divider-color);
            padding-top: 16px;
        }

        .pagination button {
            background-color: transparent;
            color: var(--primary-color);
            border: none;
            padding: 8px 16px;
            border-radius: 4px;
            cursor: pointer;
            font-family: inherit;
            font-weight: 500;
            text-transform: uppercase;
            letter-spacing: 0.05em;
            transition: background-color 0.2s ease;
        }

        .pagination button:hover:not(:disabled) {
            background-color: var(--hover-color);
        }

        .pagination button:disabled {
            color: var(--text-secondary);
            cursor: not-allowed;
            background-color: transparent;
        }

        .page-info {
            font-size: 0.875rem;
            color: var(--text-secondary);
        }

        .hidden {
            display: none !important;
        }

        .no-data {
            color: var(--text-secondary);
            font-style: italic;
        }

        /* Responsive */
        @media (max-width: 768px) {
            .app-bar-content {
                flex-direction: column;
                gap: 16px;
            }
            .report-title {
                text-align: center;
            }
            .container {
                padding: 0 16px;
                margin: 24px auto;
            }
            .data-section {
                padding: 16px;
            }
            th, td {
                padding: 10px 12px;
            }
            .controls {
                flex-direction: column;
                align-items: stretch;
            }
            .search-input {
                max-width: none;
            }
        }
    </style>
</head>
<body>
    <header class="app-bar">
        <div class="app-bar-content">
            <div class="report-title">
                <h1>Winskan Forensic Report</h1>
                <div class="timestamp">Generated on {{.Timestamp.Format "Jan 02, 2006 15:04:05 MST"}}</div>
            </div>
            <label class="theme-toggle" for="themeToggle">
                <input type="checkbox" id="themeToggle" aria-label="Toggle dark mode">
                <span class="theme-toggle-slider" aria-hidden="true"></span>
                <span class="theme-toggle-text">Dark mode</span>
            </label>
        </div>
    </header>
    
    <div class="container">
        {{if .Executables}}
        <section class="data-section">
            <h2>UserAssist: Executables</h2>
            <div class="controls">
                <input type="text" class="search-input" placeholder="Search executables...">
                <select class="page-size">
                    <option value="10">10 rows</option>
                    <option value="25">25 rows</option>
                    <option value="50">50 rows</option>
                    <option value="100">100 rows</option>
                    <option value="-1">All rows</option>
                </select>
            </div>
            <table>
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Run Count</th>
                        <th>Focus Count</th>
                        <th>Focus Time (s)</th>
                        <th>Last Run</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .Executables}}
                    <tr>
                        <td>{{.Name}}</td>
                        <td>{{.RunCount}}</td>
                        <td>{{.FocusCount}}</td>
                        <td>{{printf "%.2f" (divide .FocusTimeMs 1000.0)}}</td>
                        <td>{{if .LastRun.IsZero}}Never{{else}}{{.LastRun.Format "Jan 02, 2006 15:04:05"}}{{end}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            <div class="pagination">
                <button class="prev-page">Previous</button>
                <span class="page-info"></span>
                <button class="next-page">Next</button>
            </div>
        </section>
        {{end}}

        {{if .Shortcuts}}
        <section class="data-section">
            <h2>UserAssist: Shortcuts</h2>
            <div class="controls">
                <input type="text" class="search-input" placeholder="Search shortcuts...">
                <select class="page-size">
                    <option value="10">10 rows</option>
                    <option value="25">25 rows</option>
                    <option value="50">50 rows</option>
                    <option value="100">100 rows</option>
                    <option value="-1">All rows</option>
                </select>
            </div>
            <table>
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Run Count</th>
                        <th>Focus Count</th>
                        <th>Focus Time (s)</th>
                        <th>Last Run</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .Shortcuts}}
                    <tr>
                        <td>{{.Name}}</td>
                        <td>{{.RunCount}}</td>
                        <td>{{.FocusCount}}</td>
                        <td>{{printf "%.2f" (divide .FocusTimeMs 1000.0)}}</td>
                        <td>{{if .LastRun.IsZero}}Never{{else}}{{.LastRun.Format "Jan 02, 2006 15:04:05"}}{{end}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            <div class="pagination">
                <button class="prev-page">Previous</button>
                <span class="page-info"></span>
                <button class="next-page">Next</button>
            </div>
        </section>
        {{end}}

        {{if .USBDevices}}
        <section class="data-section">
            <h2>USB History (SYSTEM Hive)</h2>
            <div class="controls">
                <input type="text" class="search-input" placeholder="Search USB devices...">
                <select class="page-size">
                    <option value="10">10 rows</option>
                    <option value="25">25 rows</option>
                    <option value="50">50 rows</option>
                    <option value="100">100 rows</option>
                    <option value="-1">All rows</option>
                </select>
            </div>
            <table>
                <thead>
                    <tr>
                        <th>Friendly Name</th>
                        <th>Serial Number</th>
                        <th>Volume Name</th>
                        <th>Vendor</th>
                        <th>Product</th>
                        <th>Revision</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .USBDevices}}
                    <tr>
                        <td>{{if .FriendlyName}}{{.FriendlyName}}{{else}}-{{end}}</td>
                        <td>{{.SerialNumber}}</td>
                        <td>{{if .VolumeName}}{{.VolumeName}}{{else}}-{{end}}</td>
                        <td>{{if .Vendor}}{{.Vendor}}{{else}}-{{end}}</td>
                        <td>{{if .Product}}{{.Product}}{{else}}-{{end}}</td>
                        <td>{{if .Revision}}{{.Revision}}{{else}}-{{end}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            <div class="pagination">
                <button class="prev-page">Previous</button>
                <span class="page-info"></span>
                <button class="next-page">Next</button>
            </div>
        </section>
        {{end}}

        {{if .RunMRU}}
        <section class="data-section">
            <h2>RunMRU (Win + R)</h2>
            <div class="controls">
                <input type="text" class="search-input" placeholder="Search commands...">
                <select class="page-size">
                    <option value="10">10 rows</option>
                    <option value="25">25 rows</option>
                    <option value="50">50 rows</option>
                    <option value="100">100 rows</option>
                    <option value="-1">All rows</option>
                </select>
            </div>
            <table>
                <thead>
                    <tr>
                        <th style="width: 50px;">Order</th>
                        <th>Command / File</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .RunMRU}}
                    <tr>
                        <td>{{.Order}}</td>
                        <td>{{.Data}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            <div class="pagination">
                <button class="prev-page">Previous</button>
                <span class="page-info"></span>
                <button class="next-page">Next</button>
            </div>
        </section>
        {{end}}

        {{if .RecentDocs}}
        <section class="data-section">
            <h2>RecentDocs</h2>
            <div class="controls">
                <input type="text" class="search-input" placeholder="Search recent docs...">
                <select class="page-size">
                    <option value="10">10 rows</option>
                    <option value="25">25 rows</option>
                    <option value="50">50 rows</option>
                    <option value="100">100 rows</option>
                    <option value="-1">All rows</option>
                </select>
            </div>
            <table>
                <thead>
                    <tr>
                        <th style="width: 50px;">Order</th>
                        <th>File Path / Name</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .RecentDocs}}
                    <tr>
                        <td>{{.Order}}</td>
                        <td>{{.Data}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            <div class="pagination">
                <button class="prev-page">Previous</button>
                <span class="page-info"></span>
                <button class="next-page">Next</button>
            </div>
        </section>
        {{end}}

        {{if and (not .Executables) (not .Shortcuts) (not .USBDevices) (not .RunMRU) (not .RecentDocs)}}
        <section class="data-section">
            <h2>No Data Found</h2>
            <p class="no-data">The requested categories did not return any data, or the system does not have this information available.</p>
        </section>
        {{end}}
    </div>
    
    <script>
        document.addEventListener('DOMContentLoaded', () => {
            const themeToggle = document.getElementById('themeToggle');
            const themeStorageKey = 'winskan-report-theme';

            const getStoredTheme = () => {
                try {
                    return localStorage.getItem(themeStorageKey);
                } catch (error) {
                    return null;
                }
            };

            const setStoredTheme = (theme) => {
                try {
                    localStorage.setItem(themeStorageKey, theme);
                } catch (error) {
                    // Some browsers disable localStorage for local files or privacy settings.
                }
            };

            const prefersDarkMode = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
            const initialTheme = getStoredTheme() || (prefersDarkMode ? 'dark' : 'light');
            const applyTheme = (theme) => {
                const isDark = theme === 'dark';
                document.documentElement.dataset.theme = isDark ? 'dark' : 'light';
                if (themeToggle) {
                    themeToggle.checked = isDark;
                    themeToggle.setAttribute('aria-label', isDark ? 'Switch to light mode' : 'Switch to dark mode');
                }
            };

            applyTheme(initialTheme);

            if (themeToggle) {
                themeToggle.addEventListener('change', () => {
                    const selectedTheme = themeToggle.checked ? 'dark' : 'light';
                    applyTheme(selectedTheme);
                    setStoredTheme(selectedTheme);
                });
            }

            const sections = document.querySelectorAll('.data-section');
            
            sections.forEach(section => {
                const searchInput = section.querySelector('.search-input');
                const pageSizeSelect = section.querySelector('.page-size');
                const prevButton = section.querySelector('.prev-page');
                const nextButton = section.querySelector('.next-page');
                const pageInfo = section.querySelector('.page-info');
                const paginationFooter = section.querySelector('.pagination');
                const tbody = section.querySelector('tbody');
                
                if (!tbody) return;

                const allRows = Array.from(tbody.querySelectorAll('tr'));
                let filteredRows = [...allRows];
                let currentPage = 1;
                
                const minimumPageSize = 10;
                pageSizeSelect.value = String(minimumPageSize);
                let pageSize = minimumPageSize;

                const updateTable = () => {
                    pageSizeSelect.hidden = filteredRows.length < minimumPageSize;
                    paginationFooter.hidden = filteredRows.length <= minimumPageSize;
                    paginationFooter.style.display = filteredRows.length <= minimumPageSize ? 'none' : '';

                    const totalPages = pageSize === -1 ? 1 : Math.ceil(filteredRows.length / pageSize) || 1;
                    if (currentPage > totalPages) currentPage = totalPages;
                    if (currentPage < 1) currentPage = 1;

                    const startIndex = (currentPage - 1) * pageSize;
                    const endIndex = pageSize === -1 ? filteredRows.length : startIndex + pageSize;

                    allRows.forEach(row => row.classList.add('hidden'));
                    
                    for (let i = startIndex; i < endIndex && i < filteredRows.length; i++) {
                        filteredRows[i].classList.remove('hidden');
                    }

                    pageInfo.textContent = "Page " + currentPage + " of " + totalPages + " (" + filteredRows.length + " entries)";
                    prevButton.disabled = currentPage === 1;
                    nextButton.disabled = currentPage === totalPages;
                };

                const filterData = () => {
                    const query = searchInput.value.toLowerCase();
                    filteredRows = allRows.filter(row => {
                        return Array.from(row.querySelectorAll('td')).some(cell => 
                            cell.textContent.toLowerCase().includes(query)
                        );
                    });
                    currentPage = 1;
                    updateTable();
                };

                searchInput.addEventListener('input', filterData);
                
                pageSizeSelect.addEventListener('change', (e) => {
                    pageSize = parseInt(e.target.value, 10);
                    currentPage = 1;
                    updateTable();
                });

                prevButton.addEventListener('click', () => {
                    if (currentPage > 1) {
                        currentPage--;
                        updateTable();
                    }
                });

                nextButton.addEventListener('click', () => {
                    const totalPages = pageSize === -1 ? 1 : Math.ceil(filteredRows.length / pageSize);
                    if (currentPage < totalPages) {
                        currentPage++;
                        updateTable();
                    }
                });

                // Initialize table state
                updateTable();
            });
        });
    </script>
</body>
</html>
`

func generateHTMLReport(w io.Writer, data ReportData) error {
	funcMap := template.FuncMap{
		"divide": func(a uint32, b float64) float64 {
			return float64(a) / b
		},
	}
	tmpl, err := template.New("report").Funcs(funcMap).Parse(htmlTemplateStr)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}
