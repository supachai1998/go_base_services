package domain

import (
	"time"
)

const (
	DateLayout     = "2006-01-02"
	DateTimeLayout = "2006-01-02 15:04:05"
	TimeLayout     = "15:04"
)

var TimeNow = time.Now

var TimeDay = 24 * time.Hour

func TimeNowPtr() *time.Time {
	t := TimeNow()
	return &t
}

func NowUnix() int64 {
	return TimeNow().Unix()
}

func NewDate(value string) (Date, error) {
	t, err := time.ParseInLocation(DateLayout, value, time.Local)
	if err != nil {
		return Date{}, err
	}

	return Date{Time: t}, nil
}

// DateTime represents time in yyyy-MM-dd format.
type Date struct {
	time.Time
}

func (d Date) MarshalJSON() ([]byte, error) {
	return timeJSON(DateLayout, d.Time)
}

func (d *Date) UnmarshalJSON(data []byte) error {
	var err error
	d.Time, err = parseTimeJSON(DateLayout, data)

	return err
}

func (d *Date) UnmarshalParam(param string) error {
	var err error
	*d, err = NewDate(param)
	return err
}

func (d *Date) String() string {
	return d.Format(DateLayout)
}

// DateTime represents time in yyyy-MM-dd HH:mm:ss format.
type DateTime struct {
	time.Time
}

func NewDateTime(value string) (DateTime, error) {
	t, err := time.ParseInLocation(DateTimeLayout, value, time.Local)
	if err != nil {
		return DateTime{}, err
	}

	return NewDateTimeT(t), nil
}

func NewDateTimeT(t time.Time) DateTime {
	return DateTime{Time: t}
}

func (d DateTime) MarshalJSON() ([]byte, error) {
	return timeJSON(DateTimeLayout, d.Time)
}

func (d *DateTime) UnmarshalJSON(data []byte) error {
	var err error
	d.Time, err = parseTimeJSON(DateTimeLayout, data)

	return err
}

func (d *DateTime) UnmarshalParam(param string) error {
	var err error
	*d, err = NewDateTime(param)
	return err
}

func (d *DateTime) String() string {
	return d.Format(DateTimeLayout)
}

func timeJSON(layout string, t time.Time) ([]byte, error) {
	b := make([]byte, 0, len(layout)+2)

	b = append(b, '"')
	b = t.In(time.Local).AppendFormat(b, layout)
	b = append(b, '"')

	return b, nil
}

func parseTimeJSON(layout string, data []byte) (time.Time, error) {
	if string(data) == "null" {
		return time.Time{}, nil
	}

	return time.ParseInLocation(`"`+layout+`"`, string(data), time.Local)
}

// EndOfDay returns time.Time with time set to 23:59:59.999
func EndOfDay(t time.Time) time.Time {
	return t.AddDate(0, 0, 1).Add(-time.Millisecond)
}

// StartOfDay returns time.Time with time set to 00:00:00.000
func StartOfThisDay(t time.Time, l ...time.Location) time.Time {
	if len(l) > 0 {
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, &l[0])
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// EndOfThisDay returns time.Time with time set to 23:59:59.999999999
func EndOfThisDay(t time.Time, l ...time.Location) time.Time {
	ns := int((time.Second - time.Nanosecond).Nanoseconds())
	if len(l) > 0 {
		return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, ns, &l[0])
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, ns, time.UTC)
}

// ForEachDate iterates from start date to end date.
func ForEachDate(start, end time.Time, fn func(t time.Time)) {
	for d := start; d.Before(end.AddDate(0, 0, 1)); d = d.AddDate(0, 0, 1) {
		fn(d)
	}
}

// DaysDiff returns number of days from t1 to t2.
func DaysDiff(t1, t2 time.Time) int {
	return int(t2.Sub(t1).Hours() / 24)
}

func IsInTimePeriod(asIsTime, expectStartTime, expectEndTime time.Time) bool {
	asIsHour, asIsMinute := asIsTime.Hour(), asIsTime.Minute()
	if asIsHour < expectStartTime.Hour() || asIsHour > expectEndTime.Hour() {
		return false
	}

	if asIsHour > expectStartTime.Hour() && asIsHour < expectEndTime.Hour() {
		return true
	}

	if asIsHour == expectStartTime.Hour() && asIsMinute < expectStartTime.Minute() {
		return false
	}

	if asIsHour == expectEndTime.Hour() && asIsMinute > expectEndTime.Minute() {
		return false
	}

	return true
}
