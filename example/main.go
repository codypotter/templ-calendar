package main

import (
	"log"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/codypotter/templ-calendar/calendar"
)

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("example/assets"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		props := calendar.Props{
			Year:     now.Year(),
			Month:    now.Month(),
			Location: now.Location(),
			Events: []calendar.Event{
				// htmx: response replaces #event-detail
				{
					ID:    "1",
					Title: "Team standup",
					Start: time.Date(now.Year(), now.Month(), 5, 9, 0, 0, 0, time.Local),
					End:   time.Date(now.Year(), now.Month(), 5, 9, 30, 0, 0, time.Local),
					Class: "cursor-pointer hover:opacity-75",
					Attributes: templ.Attributes{
						"hx-get":    "/events/1",
						"hx-target": "#event-detail",
					},
				},
				{
					ID:    "2",
					Title: "Sprint review",
					Start: time.Date(now.Year(), now.Month(), 15, 14, 0, 0, 0, time.Local),
					End:   time.Date(now.Year(), now.Month(), 17, 15, 0, 0, 0, time.Local),
					Class: "cursor-pointer hover:opacity-75",
					Attributes: templ.Attributes{
						"hx-get":    "/events/2",
						"hx-target": "#event-detail",
					},
				},
				// alpine (inline): set selectedEvent directly in @click
				{
					ID:    "3",
					Title: "Lunch with design",
					Start: time.Date(now.Year(), now.Month(), 10, 12, 0, 0, 0, time.Local),
					End:   time.Date(now.Year(), now.Month(), 10, 13, 0, 0, 0, time.Local),
					Class: "cursor-pointer hover:opacity-75",
					Attributes: templ.Attributes{
						"@click": `selectedEvent = {id: '3', title: 'Lunch with design'}`,
					},
				},
				// alpine (data-*): handler reads attributes off $el
				{
					ID:    "4",
					Title: "Retro",
					Start: time.Date(now.Year(), now.Month(), 22, 15, 0, 0, 0, time.Local),
					End:   time.Date(now.Year(), now.Month(), 22, 16, 0, 0, 0, time.Local),
					Class: "cursor-pointer hover:opacity-75",
					Attributes: templ.Attributes{
						"@click":           "openEvent($el)",
						"data-event-id":    "4",
						"data-event-title": "Retro",
					},
				},
			},
		}
		templ.Handler(indexPage(props)).ServeHTTP(w, r)
	})

	log.Println("listening on :8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
