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
				{
					ID:    "1",
					Title: "Team standup",
					Start: time.Date(now.Year(), now.Month(), 5, 9, 0, 0, 0, time.Local),
					End:   time.Date(now.Year(), now.Month(), 5, 9, 30, 0, 0, time.Local),
				},
				{
					ID:    "2",
					Title: "Sprint review",
					Start: time.Date(now.Year(), now.Month(), 15, 14, 0, 0, 0, time.Local),
					End:   time.Date(now.Year(), now.Month(), 17, 15, 0, 0, 0, time.Local),
				},
			},
		}
		templ.Handler(indexPage(props)).ServeHTTP(w, r)
	})

	log.Println("listening on :8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
