package habit

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"sync"
	"time"
)

// Global template instance to parse once and reuse
var (
	tmpl     *template.Template
	tmplOnce sync.Once
)

// InitTemplates initializes the templates with helper functions.
// In a real app, call this in your main() or init().
func getTemplates() *template.Template {
	tmplOnce.Do(func() {
		// Define helper functions for the templates
		funcMap := template.FuncMap{
			"formatDate": func(t time.Time) string {
				return t.Format("Jan 02, 2006")
			},
			"isoDate": func(t time.Time) string {
				return t.Format("2006-01-02")
			},
			"title": func(s interface{}) string {
				str := fmt.Sprintf("%v", s)
				// Handle Go 1.18+ where strings.Title is deprecated
				return strings.ToUpper(string(str[0])) + strings.ToLower(str[1:])
			},
			"join":  strings.Join,
			"sub":   func(a, b int) int { return a - b },
			// Calculate 1 year ago for the graph
			"dateMinusOneYear": func() string {
				return time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
			},
			"dateNow": func() string {
				return time.Now().Format("2006-01-02")
			},
			// Helper to format days of week
			"formatWeekdays": func(days []time.Weekday) string {
				if len(days) == 0 {
					return ""
				}
				shortNames := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
				var result []string
				for _, d := range days {
					result = append(result, shortNames[d])
				}
				return strings.Join(result, ", ")
			},
		}

		// Parse all templates
		var err error
		tmpl, err = template.New("root").Funcs(funcMap).Parse(layoutHTML + habitCardHTML + createFormHTML)
		if err != nil {
			panic(fmt.Sprintf("failed to parse templates: %v", err))
		}
	})
	return tmpl
}

// RenderHabit renders a single habit card (safe from XSS).
func RenderHabit(h *Habit) string {
	var buf bytes.Buffer
	err := getTemplates().ExecuteTemplate(&buf, "habit-card", h)
	if err != nil {
		return fmt.Sprintf("Error rendering habit: %v", err)
	}
	return buf.String()
}

// RenderHabitsList renders just the list of habits (useful for HTMX updates).
func RenderHabitsList(habits []Habit) string {
	var buf bytes.Buffer
	// We iterate manually here because the template expects a single item,
	// or we could define a "list" template. Here is a simple loop:
	t := getTemplates()
	for _, h := range habits {
		err := t.ExecuteTemplate(&buf, "habit-card", h)
		if err != nil {
			return fmt.Sprintf("Error rendering list: %v", err)
		}
	}
	return buf.String()
}

// RenderAllHabits renders the full page.
func RenderAllHabits(habits []Habit) template.HTML {
	var buf bytes.Buffer
	// We wrap the habits in a struct if the page needs more data later
	data := struct {
		Habits []Habit
	}{
		Habits: habits,
	}
	
	err := getTemplates().ExecuteTemplate(&buf, "layout", data)
	if err != nil {
		return template.HTML(fmt.Sprintf("Error rendering page: %v", err))
	}
	return template.HTML(buf.String())
}

// --- CONSTANT TEMPLATES ---

const layoutHTML = `
{{define "layout"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Abitudini - Habit Tracker</title>
    <link rel="icon" type="image/svg+xml" href="/static/logo.svg">
    <link rel="stylesheet" href="/static/style.css">
    <script src="https://cdn.jsdelivr.net/npm/htmx.org@1.9.10"></script>
    <script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"></script>
</head>
<body>
    <header>
        <div class="container">
            <h1><img src="/static/logo.svg" alt="A" class="logo">bitudini</h1>
            <button class="btn btn-primary" 
                onclick="document.querySelector('.create-form').style.display = document.querySelector('.create-form').style.display === 'none' ? 'block' : 'none';">
                + New
            </button>
        </div>
    </header>

    <main class="container">
        {{template "create-form"}}

        <div id="habits-list">
            {{range .Habits}}
                {{template "habit-card" .}}
            {{end}}
        </div>
    </main>
    <script src="/static/main.js"></script>
</body>
</html>
{{end}}
`

const createFormHTML = `
{{define "create-form"}}
<div class="create-form" style="display: none;">
    <form hx-post="/api/habits" hx-target="#habits-list" hx-swap="afterbegin" hx-on::after-request="this.reset()">
        <input type="text" name="description" placeholder="What habit?" required>
        <input type="date" name="start_date" value="{{dateNow}}" required>
        <button type="submit" class="btn btn-primary">Create</button>
    </form>
</div>
{{end}}
`

const habitCardHTML = `
{{define "habit-card"}}
<div id="habit-{{.ID}}" class="card" hx-on::htmx:afterRequest="this.classList.add('pulse')">
    <button class="card-delete-btn"
            hx-delete="/api/habits/{{.ID}}" 
            hx-confirm="Delete this habit and all its data?" 
            hx-target="#habit-{{.ID}}" 
            hx-swap="outerHTML swap:0.5s"
            aria-label="Delete habit"
            title="Delete habit">
        ×
    </button>
    <div class="card-header">
        <div>
            <h2>{{.Description}}</h2>
            <p class="card-meta">Started on {{.StartDate | formatDate}}</p>
        </div>
    </div>

    <div id="contribution-{{.ID}}" 
         class="contribution-container"
         hx-get="/api/habits/{{.ID}}/contribution"
         hx-trigger="load"
         data-habit-id="{{.ID}}">
        <div class="contribution-grid" style="opacity: 0.5;">
            Loading...
        </div>
    </div>

    <div class="card-actions">
        {{if not .CompletedToday}}
        <button class="btn" 
                hx-post="/api/habits/{{.ID}}/done-today" 
                hx-target="#habit-{{.ID}}" 
                hx-swap="outerHTML">
            ✓ Done Today
        </button>
        {{end}}

        <div id="streak-{{.ID}}" hx-get="/api/habits/{{.ID}}/streak" hx-trigger="load"></div>
    </div>
</div>
{{end}}
`
