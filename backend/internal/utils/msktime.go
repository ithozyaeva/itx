package utils

import "time"

var mskLocation = mustLoadMSK()

func mustLoadMSK() *time.Location {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return time.FixedZone("MSK", 3*60*60)
	}
	return loc
}

// MSKLocation возвращает таймзону Москвы.
func MSKLocation() *time.Location {
	return mskLocation
}

// MSKDay усекает время до начала календарного дня в МСК и возвращает
// time.Time на 00:00:00 МСК. Используется для всех day-привязанных таблиц
// (daily_check_ins.day, daily_task_sets.day, raffles.day_key и т.д.).
func MSKDay(t time.Time) time.Time {
	tm := t.In(mskLocation)
	return time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, mskLocation)
}

// MSKToday — текущая МСК-дата на 00:00.
func MSKToday() time.Time {
	return MSKDay(time.Now())
}

// MSKEndOfDay возвращает 23:59:59 МСК для дня d.
func MSKEndOfDay(d time.Time) time.Time {
	day := MSKDay(d)
	return day.Add(24*time.Hour - time.Second)
}

// MSKNextMidnight возвращает ближайшее 00:00 МСК после now.
func MSKNextMidnight(now time.Time) time.Time {
	return MSKDay(now).Add(24 * time.Hour)
}

// DaysBetweenMSK — количество полных МСК-дней от a до b (b - a).
// Отрицательное значение, если b раньше a.
func DaysBetweenMSK(a, b time.Time) int {
	da := MSKDay(a)
	db := MSKDay(b)
	return int(db.Sub(da).Hours() / 24)
}
