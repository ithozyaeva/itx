package utils

import (
	"fmt"
	"ithozyeva/internal/models"
	"strings"
	"time"
)

func GenerateICS(event *models.Event) string {
	formatTime := func(t time.Time) string {
		return t.UTC().Format("20060102T150405Z") // iCalendar формат UTC
	}

	// Получаем таймзону события для информации
	timezone := event.Timezone
	if timezone == "" {
		timezone = "UTC"
	}

	builder := strings.Builder{}
	builder.WriteString("BEGIN:VCALENDAR\n")
	builder.WriteString("VERSION:2.0\n")
	builder.WriteString("PRODID:-//IT Khoziaeva//Event Calendar//EN\n")
	builder.WriteString("CALSCALE:GREGORIAN\n")

	builder.WriteString("BEGIN:VEVENT\n")
	builder.WriteString(fmt.Sprintf("UID:event-%d@ithozyeva.com\n", event.Id))
	builder.WriteString(fmt.Sprintf("DTSTAMP:%s\n", formatTime(time.Now())))
	// Дата события уже в UTC в базе, просто используем её
	builder.WriteString(fmt.Sprintf("DTSTART:%s\n", formatTime(event.Date)))
	builder.WriteString(fmt.Sprintf("SUMMARY:%s\n", escapeICS(event.Title)))

	// Добавляем информацию о таймзоне в описание для справки
	description := event.Description
	if timezone != "UTC" {
		description = fmt.Sprintf("%s\n\n⏰ Время указано для таймзоны: %s", description, timezone)
	}
	builder.WriteString(fmt.Sprintf("DESCRIPTION:%s\n", escapeICS(description)))

	// Место проведения
	place := event.Place
	if event.PlaceType == models.EventHybrid && event.CustomPlaceType != "" {
		place = event.CustomPlaceType + ": " + event.Place
	}
	// Видеоссылка, если есть
	if event.PlaceType == models.EventOnline {
		builder.WriteString(fmt.Sprintf("LOCATION:%s\n", escapeICS(place)))
		builder.WriteString(fmt.Sprintf("URL:%s\n", escapeICS(event.Place)))
	} else {
		builder.WriteString(fmt.Sprintf("LOCATION:%s\n", escapeICS(place)))
	}

	builder.WriteString("END:VEVENT\n")
	builder.WriteString("END:VCALENDAR\n")
	return builder.String()
}
func escapeICS(s string) string {
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		";", "\\;",
		",", "\\,",
		"\n", "\\n",
	)
	return replacer.Replace(s)
}
