package web

import (
	"html/template"
	"time"
)

var funcMap = template.FuncMap{
	"formatTime": func(t time.Time) string {
		if t.IsZero() {
			return "never"
		}
		return t.Format("2006-01-02 15:04:05 UTC")
	},
	"timeSince": func(t time.Time) string {
		if t.IsZero() {
			return "never"
		}
		d := time.Since(t).Round(time.Second)
		if d < time.Minute {
			return d.String() + " ago"
		}
		if d < time.Hour {
			return (d / time.Minute).String() + "m ago"
		}
		return (d / time.Hour).String() + "h ago"
	},
	"statusClass": func(overdue bool, failing bool) string {
		if failing {
			return "status-failing"
		}
		if overdue {
			return "status-overdue"
		}
		return "status-ok"
	},
}

var dashboardTmpl = template.Must(template.New("dashboard").Funcs(funcMap).Parse(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Cronitor Local</title>
  <style>
    body { font-family: sans-serif; margin: 2rem; background: #f9f9f9; color: #333; }
    h1 { color: #444; }
    table { width: 100%; border-collapse: collapse; background: #fff; box-shadow: 0 1px 3px rgba(0,0,0,.1); }
    th, td { padding: .75rem 1rem; text-align: left; border-bottom: 1px solid #eee; }
    th { background: #f0f0f0; font-weight: 600; }
    .status-ok { color: #2e7d32; font-weight: bold; }
    .status-overdue { color: #e65100; font-weight: bold; }
    .status-failing { color: #b71c1c; font-weight: bold; }
    .empty { text-align: center; padding: 2rem; color: #888; }
  </style>
</head>
<body>
  <h1>&#128337; Cronitor Local</h1>
  <p>Monitoring <strong>{{len .Jobs}}</strong> job(s) &mdash; last checked {{formatTime .CheckedAt}}</p>
  <table>
    <thead><tr><th>Name</th><th>Schedule</th><th>Last Success</th><th>Last Failure</th><th>Status</th></tr></thead>
    <tbody>
    {{if not .Jobs}}
      <tr><td colspan="5" class="empty">No jobs registered yet. POST to /jobs to add one.</td></tr>
    {{else}}
      {{range .Jobs}}
      <tr>
        <td>{{.Name}}</td>
        <td><code>{{.Schedule}}</code></td>
        <td>{{timeSince .LastSuccess}}</td>
        <td>{{timeSince .LastFailure}}</td>
        <td class="{{statusClass .IsOverdue .Failing}}">{{if .Failing}}FAILING{{else if .IsOverdue}}OVERDUE{{else}}OK{{end}}</td>
      </tr>
      {{end}}
    {{end}}
    </tbody>
  </table>
</body>
</html>
`))
