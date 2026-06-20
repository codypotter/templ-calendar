package calendar

import "time"

type Event struct {
	ID    string
	Title string
	Date  time.Time
}

type Props struct {
	Year   int
	Month  time.Month
	Events []Event
}
