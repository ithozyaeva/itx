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

// NextOccurrence возвращает дату ближайшего будущего вхождения события.
// Для обычных событий и для рекуррентных с датой в будущем — исходная дата.
// Для рекуррентных с прошедшей исходной датой вычисляет следующее вхождение
// исходя из repeat_period и repeat_interval.
//
// Логика дублирует фронтенд (platform-frontend/src/composables/useEventOccurrence.ts)
// — любая правка должна идти синхронно в обе стороны.
func NextOccurrence(event *models.Event, now time.Time) time.Time {
	if !event.IsRepeating || event.RepeatPeriod == nil {
		return event.Date
	}
	if !event.Date.Before(now) {
		return event.Date
	}

	interval := 1
	if event.RepeatInterval != nil && *event.RepeatInterval > 0 {
		interval = *event.RepeatInterval
	}

	switch models.RepeatPeriod(*event.RepeatPeriod) {
	case models.RepeatDaily:
		diffDays := int(now.Sub(event.Date) / (24 * time.Hour))
		next := (diffDays/interval + 1) * interval
		return event.Date.AddDate(0, 0, next)
	case models.RepeatWeekly:
		diffWeeks := int(now.Sub(event.Date) / (7 * 24 * time.Hour))
		next := (diffWeeks/interval + 1) * interval
		return event.Date.AddDate(0, 0, next*7)
	case models.RepeatMonthly:
		current := event.Date
		for current.Before(now) {
			current = current.AddDate(0, interval, 0)
		}
		return current
	case models.RepeatYearly:
		current := event.Date
		for current.Before(now) {
			current = current.AddDate(interval, 0, 0)
		}
		return current
	default:
		return event.Date
	}
}
