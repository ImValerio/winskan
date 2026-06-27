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
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;600;700&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg-color: #0f172a;
            --container-bg: #1e293b;
            --text-main: #f8fafc;
            --text-muted: #94a3b8;
            --accent-color: #38bdf8;
            --border-color: #334155;
            --hover-bg: #334155;
        }

        body {
            margin: 0;
            padding: 0;
            background-color: var(--bg-color);
            color: var(--text-main);
            font-family: 'Inter', sans-serif;
            line-height: 1.6;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 2rem;
        }

        header {
            text-align: center;
            margin-bottom: 3rem;
            padding-bottom: 2rem;
            border-bottom: 1px solid var(--border-color);
        }

        h1 {
            color: var(--accent-color);
            font-size: 2.5rem;
            margin-bottom: 0.5rem;
            font-weight: 700;
        }

        .timestamp {
            color: var(--text-muted);
            font-size: 0.9rem;
        }

        section {
            background-color: var(--container-bg);
            border-radius: 12px;
            padding: 2rem;
            margin-bottom: 2rem;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
            overflow-x: auto;
        }

        h2 {
            color: var(--accent-color);
            margin-top: 0;
            margin-bottom: 1.5rem;
            font-size: 1.5rem;
            border-bottom: 2px solid var(--border-color);
            padding-bottom: 0.5rem;
            font-weight: 600;
        }

        /* Controls: Search and Pagination Page Size */
        .controls {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 1rem;
            flex-wrap: wrap;
            gap: 1rem;
        }

        .search-input {
            padding: 0.6rem 1rem;
            border-radius: 6px;
            border: 1px solid var(--border-color);
            background-color: var(--bg-color);
            color: var(--text-main);
            font-family: inherit;
            flex: 1;
            max-width: 300px;
        }

        .search-input:focus {
            outline: none;
            border-color: var(--accent-color);
        }
        
        .page-size {
            padding: 0.6rem;
            border-radius: 6px;
            border: 1px solid var(--border-color);
            background-color: var(--bg-color);
            color: var(--text-main);
            font-family: inherit;
        }

        table {
            width: 100%;
            border-collapse: collapse;
            text-align: left;
        }

        th, td {
            padding: 1rem;
            border-bottom: 1px solid var(--border-color);
        }

        th {
            font-weight: 600;
            color: var(--text-muted);
            text-transform: uppercase;
            font-size: 0.85rem;
            letter-spacing: 0.05em;
        }

        tr:last-child td {
            border-bottom: none;
        }

        tbody tr {
            transition: background-color 0.2s ease;
        }

        tbody tr:hover {
            background-color: var(--hover-bg);
        }

        /* Pagination Buttons */
        .pagination {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-top: 1rem;
            border-top: 1px solid var(--border-color);
            padding-top: 1rem;
        }

        .pagination button {
            background-color: var(--bg-color);
            color: var(--text-main);
            border: 1px solid var(--border-color);
            padding: 0.5rem 1rem;
            border-radius: 6px;
            cursor: pointer;
            font-family: inherit;
            transition: all 0.2s ease;
        }

        .pagination button:hover:not(:disabled) {
            background-color: var(--border-color);
        }

        .pagination button:disabled {
            opacity: 0.5;
            cursor: not-allowed;
        }

        .page-info {
            font-size: 0.9rem;
            color: var(--text-muted);
        }

        .hidden {
            display: none !important;
        }

        .no-data {
            color: var(--text-muted);
            font-style: italic;
        }

        /* Responsive */
        @media (max-width: 768px) {
            .container {
                padding: 1rem;
            }
            section {
                padding: 1rem;
            }
            th, td {
                padding: 0.75rem 0.5rem;
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
    <div class="container">
        <header>
            <h1>Winskan Forensic Report</h1>
            <div class="timestamp">Generated on {{.Timestamp.Format "Jan 02, 2006 15:04:05 MST"}}</div>
        </header>

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
        <section>
            <h2>No Data Found</h2>
            <p class="no-data">The requested categories did not return any data, or the system does not have this information available.</p>
        </section>
        {{end}}
    </div>
    
    <script>
        document.addEventListener('DOMContentLoaded', () => {
            const sections = document.querySelectorAll('.data-section');
            
            sections.forEach(section => {
                const searchInput = section.querySelector('.search-input');
                const pageSizeSelect = section.querySelector('.page-size');
                const prevButton = section.querySelector('.prev-page');
                const nextButton = section.querySelector('.next-page');
                const pageInfo = section.querySelector('.page-info');
                const tbody = section.querySelector('tbody');
                
                if (!tbody) return;

                const allRows = Array.from(tbody.querySelectorAll('tr'));
                let filteredRows = [...allRows];
                let currentPage = 1;
                let pageSize = parseInt(pageSizeSelect.value, 10);

                const updateTable = () => {
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
