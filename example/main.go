package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/codypotter/templ-calendar/calendar"
)

func main() {
	staticOut := flag.String("static", "", "render static HTML to this file and exit")
	flag.Parse()

	if *staticOut != "" {
		if err := renderStatic(*staticOut); err != nil {
			log.Fatal(err)
		}
		return
	}

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("example/assets"))))
	http.HandleFunc("/events/", handleEvent)
	http.HandleFunc("/", handleIndex)
	log.Println("listening on :8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}

func renderStatic(out string) error {
	now := time.Now()
	year, month := now.Year(), now.Month()
	calProps := calendar.Props{
		Year:        year,
		Month:       month,
		Location:    time.Local,
		Events:      exampleEvents(year, month),
		HideHeading: true,
	}
	navProps := calendar.NavigatorProps{Year: year, Month: month}
	jumperProps := calendar.JumperProps{Year: year, Month: month}

	if err := os.MkdirAll(filepath.Dir(out), 0755); err != nil {
		return err
	}
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	return indexPage(calProps, navProps, jumperProps).Render(context.Background(), f)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	year, month := parseMonthYear(r)

	prev := time.Date(year, month-1, 1, 0, 0, 0, 0, time.Local)
	next := time.Date(year, month+1, 1, 0, 0, 0, 0, time.Local)

	navProps := calendar.NavigatorProps{
		Year:  year,
		Month: month,
		PrevAttributes: templ.Attributes{
			"hx-get":    fmt.Sprintf("/?year=%d&month=%d", prev.Year(), int(prev.Month())),
			"hx-target": "#calendar-container",
			"hx-select": "#calendar-container",
			"hx-swap":   "outerHTML",
		},
		NextAttributes: templ.Attributes{
			"hx-get":    fmt.Sprintf("/?year=%d&month=%d", next.Year(), int(next.Month())),
			"hx-target": "#calendar-container",
			"hx-select": "#calendar-container",
			"hx-swap":   "outerHTML",
		},
	}

	calProps := calendar.Props{
		Year:        year,
		Month:       month,
		Location:    time.Local,
		Events:      exampleEvents(year, month),
		HideHeading: true,
	}

	jumperProps := calendar.JumperProps{
		Year:  year,
		Month: month,
		FormAttributes: templ.Attributes{
			"hx-get":     "/",
			"hx-target":  "#calendar-container",
			"hx-select":  "#calendar-container",
			"hx-swap":    "outerHTML",
			"hx-trigger": "change",
		},
	}

	templ.Handler(indexPage(calProps, navProps, jumperProps)).ServeHTTP(w, r)
}

func handleEvent(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/events/"):]
	fmt.Fprintf(w, `<p>Event <strong>%s</strong> detail would go here.</p>`, id)
}

func parseMonthYear(r *http.Request) (int, time.Month) {
	now := time.Now()
	year, month := now.Year(), now.Month()
	if y := r.URL.Query().Get("year"); y != "" {
		if parsed, err := strconv.Atoi(y); err == nil {
			year = parsed
		}
	}
	if m := r.URL.Query().Get("month"); m != "" {
		if parsed, err := strconv.Atoi(m); err == nil && parsed >= 1 && parsed <= 12 {
			month = time.Month(parsed)
		}
	}
	return year, month
}

func exampleEvents(year int, month time.Month) []calendar.Event {
	return []calendar.Event{
		// htmx: response replaces #event-detail
		{
			ID:    "1",
			Title: "Team standup",
			Start: time.Date(year, month, 5, 9, 0, 0, 0, time.Local),
			End:   time.Date(year, month, 5, 9, 30, 0, 0, time.Local),
			Class: "cursor-pointer hover:opacity-75",
			Attributes: templ.Attributes{
				"hx-get":    "/events/1",
				"hx-target": "#event-detail",
			},
		},
		{
			ID:    "2",
			Title: "Sprint review",
			Start: time.Date(year, month, 15, 14, 0, 0, 0, time.Local),
			End:   time.Date(year, month, 17, 15, 0, 0, 0, time.Local),
			Class: "cursor-pointer hover:opacity-75 bg-violet-100 text-violet-800 dark:bg-violet-900 dark:text-violet-200",
			Attributes: templ.Attributes{
				"hx-get":    "/events/2",
				"hx-target": "#event-detail",
			},
		},
		// alpine (inline): set selectedEvent directly in @click
		{
			ID:    "3",
			Title: "Lunch with design",
			Start: time.Date(year, month, 10, 12, 0, 0, 0, time.Local),
			End:   time.Date(year, month, 10, 13, 0, 0, 0, time.Local),
			Class: "cursor-pointer hover:opacity-75",
			Attributes: templ.Attributes{
				"@click": `selectedEvent = {id: '3', title: 'Lunch with design'}`,
			},
		},
		// alpine (data-*): handler reads attributes off $el
		{
			ID:    "4",
			Title: "Retro",
			Start: time.Date(year, month, 22, 15, 0, 0, 0, time.Local),
			End:   time.Date(year, month, 22, 16, 0, 0, 0, time.Local),
			Class: "cursor-pointer hover:opacity-75",
			Attributes: templ.Attributes{
				"@click":           "openEvent($el)",
				"data-event-id":    "4",
				"data-event-title": "Retro",
			},
		},
	}
}
