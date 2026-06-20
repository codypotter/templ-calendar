package templcalendar

import "embed"

//go:embed calendar/calendar.templ calendar/navigator.templ calendar/jumper.templ
var Files embed.FS
