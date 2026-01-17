// Package utils provides utility functions for USC platform services.
package utils

import (
	"fmt"
	"time"
)

// TimeUtils provides time utility functions
type TimeUtils struct{}

// NewTimeUtils creates a new time utils instance
func NewTimeUtils() *TimeUtils {
	return &TimeUtils{}
}

// Now returns the current time
func (tu *TimeUtils) Now() time.Time {
	return time.Now()
}

// UTCNow returns the current time in UTC
func (tu *TimeUtils) UTCNow() time.Time {
	return time.Now().UTC()
}

// ParseTime parses a time string
func (tu *TimeUtils) ParseTime(timeStr string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02",
		"15:04:05",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, timeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", timeStr)
}

// FormatTime formats a time using the specified layout
func (tu *TimeUtils) FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

// FormatTimeISO formats a time in ISO 8601 format
func (tu *TimeUtils) FormatTimeISO(t time.Time) string {
	return t.Format(time.RFC3339)
}

// FormatTimeDate formats a time as date only
func (tu *TimeUtils) FormatTimeDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatTimeDateTime formats a time as date and time
func (tu *TimeUtils) FormatTimeDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// AddDays adds days to a time
func (tu *TimeUtils) AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// AddHours adds hours to a time
func (tu *TimeUtils) AddHours(t time.Time, hours int) time.Time {
	return t.Add(time.Duration(hours) * time.Hour)
}

// AddMinutes adds minutes to a time
func (tu *TimeUtils) AddMinutes(t time.Time, minutes int) time.Time {
	return t.Add(time.Duration(minutes) * time.Minute)
}

// AddSeconds adds seconds to a time
func (tu *TimeUtils) AddSeconds(t time.Time, seconds int) time.Time {
	return t.Add(time.Duration(seconds) * time.Second)
}

// AddMonths adds months to a time
func (tu *TimeUtils) AddMonths(t time.Time, months int) time.Time {
	return t.AddDate(0, months, 0)
}

// AddYears adds years to a time
func (tu *TimeUtils) AddYears(t time.Time, years int) time.Time {
	return t.AddDate(years, 0, 0)
}

// StartOfDay returns the start of day for a time
func (tu *TimeUtils) StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the end of day for a time
func (tu *TimeUtils) EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}

// StartOfWeek returns the start of week for a time
func (tu *TimeUtils) StartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday is 0, make it 7
	}
	return tu.StartOfDay(t.AddDate(0, 0, -weekday+1))
}

// EndOfWeek returns the end of week for a time
func (tu *TimeUtils) EndOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday is 0, make it 7
	}
	return tu.EndOfDay(t.AddDate(0, 0, 7-weekday))
}

// StartOfMonth returns the start of month for a time
func (tu *TimeUtils) StartOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth returns the end of month for a time
func (tu *TimeUtils) EndOfMonth(t time.Time) time.Time {
	return tu.StartOfMonth(t.AddDate(0, 1, 0)).Add(-time.Nanosecond)
}

// StartOfYear returns the start of year for a time
func (tu *TimeUtils) StartOfYear(t time.Time) time.Time {
	year, _, _ := t.Date()
	return time.Date(year, 1, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear returns the end of year for a time
func (tu *TimeUtils) EndOfYear(t time.Time) time.Time {
	return tu.StartOfYear(t.AddDate(1, 0, 0)).Add(-time.Nanosecond)
}

// IsSameDay checks if two times are on the same day
func (tu *TimeUtils) IsSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// IsSameWeek checks if two times are in the same week
func (tu *TimeUtils) IsSameWeek(t1, t2 time.Time) bool {
	s1 := tu.StartOfWeek(t1)
	s2 := tu.StartOfWeek(t2)
	return s1.Equal(s2)
}

// IsSameMonth checks if two times are in the same month
func (tu *TimeUtils) IsSameMonth(t1, t2 time.Time) bool {
	y1, m1, _ := t1.Date()
	y2, m2, _ := t2.Date()
	return y1 == y2 && m1 == m2
}

// IsSameYear checks if two times are in the same year
func (tu *TimeUtils) IsSameYear(t1, t2 time.Time) bool {
	y1, _, _ := t1.Date()
	y2, _, _ := t2.Date()
	return y1 == y2
}

// IsToday checks if a time is today
func (tu *TimeUtils) IsToday(t time.Time) bool {
	return tu.IsSameDay(t, time.Now())
}

// IsYesterday checks if a time is yesterday
func (tu *TimeUtils) IsYesterday(t time.Time) bool {
	return tu.IsSameDay(t, time.Now().AddDate(0, 0, -1))
}

// IsTomorrow checks if a time is tomorrow
func (tu *TimeUtils) IsTomorrow(t time.Time) bool {
	return tu.IsSameDay(t, time.Now().AddDate(0, 0, 1))
}

// IsWeekend checks if a time is on a weekend
func (tu *TimeUtils) IsWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// IsWeekday checks if a time is on a weekday
func (tu *TimeUtils) IsWeekday(t time.Time) bool {
	return !tu.IsWeekend(t)
}

// DaysBetween returns the number of days between two times
func (tu *TimeUtils) DaysBetween(t1, t2 time.Time) int {
	duration := t2.Sub(t1)
	return int(duration.Hours() / 24)
}

// HoursBetween returns the number of hours between two times
func (tu *TimeUtils) HoursBetween(t1, t2 time.Time) int {
	duration := t2.Sub(t1)
	return int(duration.Hours())
}

// MinutesBetween returns the number of minutes between two times
func (tu *TimeUtils) MinutesBetween(t1, t2 time.Time) int {
	duration := t2.Sub(t1)
	return int(duration.Minutes())
}

// SecondsBetween returns the number of seconds between two times
func (tu *TimeUtils) SecondsBetween(t1, t2 time.Time) int {
	duration := t2.Sub(t1)
	return int(duration.Seconds())
}

// Age calculates the age from a birth date
func (tu *TimeUtils) Age(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()
	if now.YearDay() < birthDate.YearDay() {
		age--
	}
	return age
}

// IsLeapYear checks if a year is a leap year
func (tu *TimeUtils) IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// DaysInMonth returns the number of days in a month
func (tu *TimeUtils) DaysInMonth(year, month int) int {
	return time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()
}

// DaysInYear returns the number of days in a year
func (tu *TimeUtils) DaysInYear(year int) int {
	if tu.IsLeapYear(year) {
		return 366
	}
	return 365
}

// TimezoneInfo provides timezone information
type TimezoneInfo struct {
	Name     string         `json:"name"`
	Offset   int            `json:"offset"`
	Location *time.Location `json:"-"`
}

// GetTimezoneInfo returns timezone information
func (tu *TimeUtils) GetTimezoneInfo(tz string) (*TimezoneInfo, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, err
	}

	now := time.Now().In(loc)
	_, offset := now.Zone()

	return &TimezoneInfo{
		Name:     tz,
		Offset:   offset,
		Location: loc,
	}, nil
}

// ConvertTimezone converts a time to a different timezone
func (tu *TimeUtils) ConvertTimezone(t time.Time, tz string) (time.Time, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(loc), nil
}

// TimeRange represents a time range
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// NewTimeRange creates a new time range
func (tu *TimeUtils) NewTimeRange(start, end time.Time) *TimeRange {
	return &TimeRange{
		Start: start,
		End:   end,
	}
}

// Contains checks if a time is within the range
func (tr *TimeRange) Contains(t time.Time) bool {
	return t.After(tr.Start) && t.Before(tr.End)
}

// Overlaps checks if two time ranges overlap
func (tr *TimeRange) Overlaps(other *TimeRange) bool {
	return tr.Start.Before(other.End) && tr.End.After(other.Start)
}

// Duration returns the duration of the time range
func (tr *TimeRange) Duration() time.Duration {
	return tr.End.Sub(tr.Start)
}

// IsValid checks if the time range is valid
func (tr *TimeRange) IsValid() bool {
	return tr.Start.Before(tr.End)
}
