package timeutils

import (
	"strconv"
	"time"
)

// GetActualGMTDate - Returns the current GMT date time into the RFC1123 format
func GetActualGMTDate() string {

	loc, _ := time.LoadLocation("GMT")
	currentTime := time.Now().In(loc)

	return currentTime.Format(time.RFC1123)
}

// GetTimeID - Generate one unic ID based on time
func GetTimeID() string {
	t := time.Now()

	y := t.Year()
	mon := t.Month()
	d := t.Day()
	h := t.Hour()
	m := t.Minute()
	s := t.Second()
	n := t.Nanosecond()

	timeID := strconv.Itoa(y) + strconv.Itoa(int(mon)) + strconv.Itoa(d) + strconv.Itoa(h) + strconv.Itoa(m) + strconv.Itoa(s) + strconv.Itoa(n)

	return timeID
}
