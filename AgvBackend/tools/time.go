package tools

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//
// const (
// 	TODAY_FORMAT         = "%Y-%m-%d"
// 	TIME_FORMAT          = "%H:%M:%S"
// 	NOW_FORMAT           = "%Y-%m-%d %H:%M:%S"
// 	PERIOD_FORMAT        = "%Y-%m-%d %H:%M"
// 	NOW_FILE_NAME_FORMAT = "%Y-%m-%d-%H-%M-%S"
// 	THIS_WEEK_FORMAT     = "%Y-%m-%V"
// 	THIS_MONTH_FORMAT    = "%Y-%m"
// 	THIS_YEAR_FORMAT     = "%Y"
// )

// V golangu sa to robi takto naprd
const (
	TODAY_FORMAT         = "2006-01-02"
	TIME_FORMAT          = "15:04:05"
	NOW_FORMAT           = "2006-01-02 15:04:05"
	PERIOD_FORMAT        = "2006-01-02 15:04"
	NOW_FILE_NAME_FORMAT = "2006-01-02-15-04-05"
	THIS_WEEK_FORMAT     = "2006-01-02"
	THIS_MONTH_FORMAT    = "2006-01"
	THIS_YEAR_FORMAT     = "2006"
)

// Duration converts milliseconds to a formatted string
func DurationFromMillis(millis int64) string {
	if millis == 0 {
		return "00:00"
	}

	seconds := millis / 1000
	days := seconds / 86400
	seconds %= 86400
	hours := seconds / 3600
	seconds %= 3600
	minutes := seconds / 60
	seconds %= 60

	daysStr := fmt.Sprintf("%02d:", days)
	hoursStr := fmt.Sprintf("%02d:", hours)
	minutesStr := fmt.Sprintf("%02d:", minutes)
	secondsStr := fmt.Sprintf("%02d", seconds)

	return daysStr + hoursStr + minutesStr + secondsStr
}

// GetTimeNow returns the current UTC time
func GetTimeNow() time.Time {
	return time.Now().UTC()
}

// GetTimeFromMillis converts milliseconds to time.Time
func GetTimeFromMillis(millis int64) time.Time {
	return time.Unix(millis/1000, (millis%1000)*int64(time.Millisecond)).UTC()
}

// GetTimeNowMillis returns the current time in milliseconds
func GetTimeNowMillis() Milliseconds {
	return Milliseconds(GetTimeNow().UnixMilli())
}

// GetTimeNowMillis returns the current time in milliseconds
func GetNow() Milliseconds {
	return Milliseconds(GetTimeNow().UnixMilli())
}

// GetTimeNowSeconds returns the current time in seconds
func GetTimeNowSeconds() int64 {
	return GetTimeNow().Unix()
}

// GetMillisFromDatetime returns milliseconds from a time.Time object
func GetMillisFromDatetime(datetime time.Time) int64 {
	return datetime.UnixMilli()
}

// GetNowFormatted returns the current time formatted
func GetNowFormatted() string {
	return GetTimeNow().Format(NOW_FORMAT)
}

// GetNowFileFormatted returns the current time formatted for file naming
func GetNowFileFormatted() string {
	return GetTimeNow().Format(NOW_FILE_NAME_FORMAT)
}

// GetTodayFormatted returns today's date formatted
func GetTodayFormatted() string {
	return GetTimeNow().Format(TODAY_FORMAT)
}

// GetThisWeekFormatted returns the current week formatted
func GetThisWeekFormatted() string {
	return GetTimeNow().Format(THIS_WEEK_FORMAT)
}

// GetThisMonthFormatted returns the current month formatted
func GetThisMonthFormatted() string {
	return GetTimeNow().Format(THIS_MONTH_FORMAT)
}

// GetThisYearFormatted returns the current year formatted
func GetThisYearFormatted() string {
	return GetTimeNow().Format(THIS_YEAR_FORMAT)
}

// GetServerTimeZone returns the current server time zone
func GetServerTimeZone() string {
	_, offset := time.Now().Zone()
	return fmt.Sprintf("%+03d:%02d", offset/3600, (offset%3600)/60)
}

// DeparseLastVisit parses a visit string into a time.Time object
func DeparseLastVisit(visitString string) (time.Time, error) {
	return time.Parse(NOW_FORMAT, visitString)
}

// IsGreaterThanMinutes checks if timeA is greater than timeB by given minutes
func IsGreaterThanMinutes(timeA, timeB time.Time, minutes int) bool {
	if timeA.Day() != timeB.Day() {
		return true
	}
	diff := timeA.Sub(timeB)
	return diff > time.Duration(minutes)*time.Minute
}

// GetMillisFromMS converts "MM:SS" string to milliseconds
func GetMillisFromMS(timeString string) (int64, error) {
	var minutes, seconds int
	_, err := fmt.Sscanf(timeString, "%d:%d", &minutes, &seconds)
	if err != nil {
		return 0, err
	}
	milliseconds := (minutes*60 + seconds) * 1000
	return int64(milliseconds), nil
}

// GetMillisFromSevilo converts timestamp string to milliseconds
func GetMillisFromSewio(timeString string) int64 {
	truncatedTimeString := timeString[:23]
	datetimeObject, err := time.Parse("2006-01-02 15:04:05.9999999", truncatedTimeString)
	if err != nil {
		return 0
	}
	return GetMillisFromDatetime(datetimeObject)
}

// GetMillisFromTime converts timestamp string to milliseconds
func GetMillisFromTime(timeString string) (int64, error) {
	datetimeObject, err := time.Parse("2006-01-02 15:04:05", timeString)
	if err != nil {
		return 0, err
	}
	return GetMillisFromDatetime(datetimeObject), nil
}

// GetDatetimeFormattedFromMillis returns formatted datetime from milliseconds
func GetDatetimeFormattedFromMillis(millis int64) string {
	return GetTimeFromMillis(millis).Format(NOW_FORMAT)
}

// GetDatetimeFromMillis returns time.Time from milliseconds
func GetDatetimeFromMillis(millis int64) time.Time {
	return GetTimeFromMillis(millis)
}

// GetFormattedDateTimeFromMillis returns formatted datetime from milliseconds
func GetFormattedDateTimeFromMillis(millis int64) string {
	return GetDatetimeFromMillis(millis).Format(NOW_FORMAT)
}

// MilisToSeconds converts milliseconds to seconds
func MilisToSeconds(millis int64) int64 {
	return millis / 1000
}

// SecondsToMilis converts seconds to milliseconds
func SecondsToMilis(seconds int64) int64 {
	return seconds * 1000
}

// MilisToMinutes converts milliseconds to minutes
func MilisToMinutes(millis int64) int64 {
	return MilisToSeconds(millis) / 60
}

// MilisToHours converts milliseconds to hours
func MilisToHours(millis int64) int64 {
	return MilisToMinutes(millis) / 60
}

// TimeFromHours converts "HH:MM:SS" string to time.Time
func TimeFromHours(timeStr string) (time.Time, error) {
	return time.Parse(TIME_FORMAT, timeStr)
}

// MilisFromTime converts time.Time to milliseconds
func MilisFromTime(dttime time.Time) int64 {
	delta := time.Duration(dttime.Hour())*time.Hour +
		time.Duration(dttime.Minute())*time.Minute +
		time.Duration(dttime.Second())*time.Second
	return int64(delta.Seconds()) * 1000
}

// GetDaysBetween returns a list of formatted dates between two timestamps
func GetDaysBetween(min, max int64) []string {
	result := []string{}

	for currentDateMS := min; currentDateMS <= max; currentDateMS += 24 * 60 * 60 * 1000 {
		currentDate := GetDatetimeFromMillis(currentDateMS)
		result = append(result, currentDate.Format(TODAY_FORMAT))
	}

	// Check if max date needs to be added
	lastDate := GetDatetimeFromMillis(max).Format(TODAY_FORMAT)
	if !contains(result, lastDate) {
		result = append(result, lastDate)
	}
	return result
}

// GetDaysCountBetween returns the number of days between two timestamps
func GetDaysCountBetween(min, max int64) int64 {
	differenceMS := max - min
	return differenceMS / (24 * 60 * 60 * 1000)
}

// GetHoursCountBetween returns the number of hours between two timestamps
func GetHoursCountBetween(min, max int64) int64 {
	differenceMS := max - min
	return differenceMS / (60 * 60 * 1000)
}

// GetHoursOfDay returns the list of hours in a day
func GetHoursOfDay() []string {
	hours := make([]string, 24)
	for i := 0; i < 24; i++ {
		hours[i] = fmt.Sprintf("%02d", i)
	}
	return hours
}

func GetHourOfTimeStamp(timestamp int64, offsetFromUTC int) int {
	seconds := timestamp / 1000

	t := time.Unix(seconds, 0)

	location := time.FixedZone("Offset", offsetFromUTC*3600)
	localTime := t.In(location)

	return localTime.Hour()
}

// GetDayYMDFromMillis returns formatted date from milliseconds
func GetDayYMDFromMillis(millis int64) string {
	return GetTimeFromMillis(millis).Format(TODAY_FORMAT)
}

// GetDayAndHoursFromTo0HFromMillis returns day and hours percentage from milliseconds
func GetDayAndHoursFromTo0HFromMillis(minMillis, maxMillis int64) [][]interface{} {
	result := [][]interface{}{}

	minDate := GetDatetimeFromMillis(minMillis)
	maxDate := GetDatetimeFromMillis(maxMillis)

	for currentDate := minDate; currentDate.Before(maxDate) || currentDate.Equal(maxDate); currentDate = currentDate.Add(time.Hour) {
		hourStart := currentDate
		hourEnd := hourStart.Add(time.Hour)
		if hourEnd.After(maxDate) {
			hourEnd = maxDate
		}
		hourDuration := hourEnd.Sub(hourStart).Minutes()

		if hourDuration == 0 {
			hourDuration = 1
		}

		percentageSpent := hourEnd.Sub(max(hourStart, minDate)).Minutes() / hourDuration

		result = append(result, []interface{}{hourStart.Format(TODAY_FORMAT), hourStart.Format("15"), percentageSpent})
	}
	return result
}

// GetMillisOfWeek returns milliseconds in a week
func GetMillisOfWeek() int64 {
	return 1000 * 60 * 60 * 24 * 7
}

// GetSecondsOfWeek returns seconds in a week
func GetSecondsOfWeek() int64 {
	return 60 * 60 * 24 * 7
}

// GetSecondsOfShift returns seconds in an 8-hour shift
func GetSecondsOfShift() int64 {
	return 60 * 60 * 8
}

// Helper function to check if a date string is contained in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Helper function to get the maximum time.Time
func max(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func ParseElapsedTimeToMillis(timeStr string) (int64, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid time format")
	}

	minutes, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}

	seconds, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	totalMillis := int64(minutes*60*1000 + seconds*1000)
	return totalMillis, nil
}
