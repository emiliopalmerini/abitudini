# Abitudini - Simple Habit Tracker

A minimalist habit tracking application built with Go, HTMX, and Alpine.js using vertical slice architecture.

## Features

- **Track daily habits** with GitHub-style contribution graphs
- **Inline editing** - click habit title to edit description and color
- **Streak tracking** - automatic calculation of current streaks
- **Read-only contribution graph** - visual representation of habit completion
- **Single-user** - personal use optimized

## Tech Stack

- **Backend**: Go (stdlib + `database/sql`)
- **Database**: SQLite3
- **Frontend**: HTMX + Alpine.js
- **Architecture**: Vertical slice (feature-driven)

## Project Structure

```
abitudini/
├── internal/
│   ├── db/
│   │   ├── db.go          # Database initialization
│   │   └── migrations.go   # Schema migrations
│   └── habit/
│       ├── handler.go      # HTTP handlers
│       ├── models.go       # Domain models
│       ├── store.go        # Data access layer
│       └── views.go        # Template rendering
├── static/
│   ├── style.css           # Styles
│   └── main.js             # Client-side logic
├── main.go                 # Server setup
├── go.mod
└── go.sum
```

## Getting Started

### Prerequisites

- Go 1.23+
- SQLite3

### Installation

```bash
git clone <repo>
cd abitudini
go mod download
go build -o abitudini .
```

### Running

```bash
./abitudini
```

Server starts on `http://localhost:8080`

## API Endpoints

### Habits
- `POST /api/habits` - Create habit
- `GET /api/habits` - Get all habits
- `GET /api/habits/{id}` - Get habit by ID
- `PUT /api/habits/{id}` - Update habit
- `POST /api/habits/{id}/done-today` - Mark as done today
- `GET /api/habits/{id}/streak` - Get streak count
- `GET /api/habits/{id}/contribution?from=YYYY-MM-DD&to=YYYY-MM-DD` - Get contribution data

## Data Model

### Habit
- `id`: Integer (PK)
- `description`: String
- `start_date`: Date
- `color`: Hex color
- `created_at`: Timestamp

### Record
- `habit_id`: FK to habits
- `record_date`: Date (unique per habit)
- `completed_at`: Timestamp

## Features Detail

### Streak Logic
- Broken after one missed day
- Calculates backward from today
- Consecutive days with completion

### Contribution Graph
- Displays completed vs. incomplete days
- Grouped by week (7 days)
- Read-only (shows past activity)
- Hover for date tooltip
- Completed days styled with accent color

## Database Migrations

Migrations run automatically on startup. No manual setup needed.

Schema includes:
- `habits` table
- `records` table (completion history)
- Indexes on frequently queried columns

## Development Notes

### Vertical Slice Architecture

Each feature is organized vertically:
```
habit/
├── models.go    # Domain types
├── store.go     # Data layer
├── handler.go   # HTTP handlers
└── views.go     # Rendering
```

No traditional layer separation—each "slice" (habit, streak, record) contains all needed code.

### HTMX Integration

- Form submissions return updated HTML
- HTMX handles all dynamic updates
- `hx-swap="outerHTML"` for card replacement
- Server renders complete card HTML on response

### Alpine.js

- Minimal interactivity: toggle edit form visibility
- Used for event dispatch (`$dispatch`)
- No complex state management needed
- Plain HTML enhanced with directives

## Future Enhancements

- Multiple user support
- Custom reminders/notifications
- Data export (CSV/JSON)
- Dark mode
- Mobile app
