package tools

import (
	"fmt"
	"math"
	"time"
)

type DtTimeRange struct {
	Period        string `json:"period"`
	TimeType      string `json:"time_type"`
	Offset        int    `json:"offset"`    // UTC offset of From fields (e.g. for UTC+1 this would be -1)
	OffsetTo      int    `json:"offset_to"` // UTC offset of To fields
	Year          int    `json:"year"`
	YearTo        int    `json:"year_to"`
	Month         int    `json:"month"`
	MonthTo       int    `json:"month_to"`
	Week          int    `json:"week"`
	Day           int    `json:"day"`
	DayTo         int    `json:"day_to"`
	Hour          int    `json:"hour"`
	HourTo        int    `json:"hour_to"`
	Minute        int    `json:"minute"`
	MinuteTo      int    `json:"minute_to"`
	Second        int    `json:"second"`
	SecondTo      int    `json:"second_to"`
	Shift         int    `json:"shift"`
	Quarter       int    `json:"quarter"`
	ShrinkedFrom  int64  `json:"shrinked_from"`
	ShrinkedTo    int64  `json:"shrinked_to"`
	CustomFrom    int64  `json:"custom_from"`
	CustomTo      int64  `json:"custom_to"`
	TimestampFrom int64  `json:"timestamp_from"`
	TimestampTo   int64  `json:"timestamp_to"`
}

// DtTimeRangeRunningStatus defines possible states
type DtTimeRangeRunningStatus int

const (
	IN DtTimeRangeRunningStatus = iota
	OUT
)

// NewDtTimeRange constructor
func NewDtTimeRange() *DtTimeRange {
	now := GetTimeNow()

	return &DtTimeRange{
		Period:       "month",
		TimeType:     "day",
		Offset:       0,
		OffsetTo:     0,
		Year:         now.Year(),
		YearTo:       now.Year(),
		Month:        int(now.Month()),
		MonthTo:      int(now.Month()),
		Week:         GetWeek(now),
		Day:          0,
		DayTo:        0,
		Hour:         0,
		HourTo:       23,
		Minute:       0,
		MinuteTo:     59,
		Second:       0,
		SecondTo:     59,
		Shift:        0,
		Quarter:      1,
		ShrinkedFrom: math.MaxInt64,
		ShrinkedTo:   0,
		CustomFrom:   0,
		CustomTo:     0,
	}
}

// FromStartEnd sets the start and end datetime range
func (self *DtTimeRange) FromStartEnd(start, end time.Time) {
	self.Period = "custom"
	self.TimeType = "custom"

	self.Year = start.Year()
	self.YearTo = end.Year()
	self.Month = int(start.Month())
	self.MonthTo = int(end.Month())
	self.Day = start.Day()
	self.Hour = start.Hour()
	self.Minute = start.Minute()
	self.DayTo = end.Day()
	self.HourTo = end.Hour()
	self.MinuteTo = end.Minute()
	self.Second = start.Second()
	self.SecondTo = end.Second()
}

// FromStartEndMillis sets the start and end datetime range from milliseconds
func (self *DtTimeRange) FromStartEndMillis(start, end int64) {
	self.CustomFrom = start
	self.CustomTo = end

	startDate := time.Unix(start/1000, 0).UTC()
	endDate := time.Unix(end/1000, 0).UTC()
	self.FromStartEnd(startDate, endDate)
}

// IsAll checks if the period is set to "all"
func (self *DtTimeRange) IsAll() bool {
	return self.Period == "all"
}

// Shrink shrinks the time range based on the stored shrinked values
func (self *DtTimeRange) Shrink() {
	if self.ShrinkedTo == 0 {
		return
	}
	self.FromStartEndMillis(self.ShrinkedFrom, self.ShrinkedTo)
}

func (self *DtTimeRange) AddShrinkTime(time_stamp int64) {
	if time_stamp > self.ShrinkedTo {
		self.ShrinkedTo = time_stamp
	}
	if time_stamp < self.ShrinkedFrom {
		self.ShrinkedFrom = time_stamp
	}
}

func (calc *DtTimeRange) CalculateStartEnd() ([]time.Time, error) {
	var startEnd []time.Time

	// Check if period is "all"
	if calc.Period == "all" {
		startDate := time.Date(calc.Year, time.Month(calc.Month), calc.Day, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(calc.Year, time.Month(calc.Month), calc.Day, 23, 59, 59, 0, time.UTC)
		return []time.Time{startDate, endDate}, nil
	}

	// Check timeType and period
	if calc.TimeType == "day" {
		switch calc.Period {
		case "day":
			startDate := time.Date(calc.Year, time.Month(calc.Month), calc.Day, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(calc.Year, time.Month(calc.Month), calc.Day, 23, 59, 59, 0, time.UTC)
			startEnd = append(startEnd, startDate, endDate)

		case "week":
			startDate := time.Date(calc.Year, time.Month(calc.Month), calc.Day, 0, 0, 0, 0, time.UTC)
			startDate = startDate.AddDate(0, 0, -(int(startDate.Weekday()) - 1)) // Adjust to the first day of the week
			endDate := startDate.AddDate(0, 0, 6).Add(time.Hour*23 + time.Minute*59 + time.Second*59)
			startEnd = append(startEnd, startDate, endDate)

		case "month":
			daysInMonth := time.Date(calc.Year, time.Month(calc.Month)+1, 0, 23, 59, 59, 0, time.UTC).Day()
			startDate := time.Date(calc.Year, time.Month(calc.Month), 1, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(calc.Year, time.Month(calc.Month), daysInMonth, 23, 59, 59, 0, time.UTC)
			startEnd = append(startEnd, startDate, endDate)

		case "quarter":
			startMonth := 3*(calc.Quarter-1) + 1
			startDate := time.Date(calc.Year, time.Month(startMonth), 1, 0, 0, 0, 0, time.UTC)
			endDate := startDate.AddDate(0, 3, 0).Add(-time.Second) // End at the last second of the quarter

			startEnd = append(startEnd, startDate, endDate)

		case "year":
			startDate := time.Date(calc.Year, 1, 1, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(calc.Year, 12, 31, 23, 59, 59, 0, time.UTC)
			startEnd = append(startEnd, startDate, endDate)

		case "custom":
			startDate := time.Date(calc.Year, time.Month(calc.Month), calc.Day, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(calc.YearTo, time.Month(calc.MonthTo), calc.DayTo, 23, 59, 59, 0, time.UTC)
			startEnd = append(startEnd, startDate, endDate)
		}

	} else { // Handle timeType != "day"
		switch calc.Period {
		case "day":
			startDate := time.Date(calc.Year, time.Month(calc.Month), calc.Day, calc.Hour, calc.Minute, 0, 0, time.UTC)
			endDate := time.Date(calc.YearTo, time.Month(calc.MonthTo), calc.DayTo, calc.HourTo, calc.MinuteTo, 59, 0, time.UTC)
			startEnd = append(startEnd, startDate, endDate)

		case "week":
			startDate := time.Date(calc.Year, time.Month(calc.Month), calc.Day, calc.Hour, calc.Minute, 0, 0, time.UTC)
			endDate := time.Date(calc.Year, time.Month(calc.Month), calc.Day, calc.HourTo, calc.MinuteTo, 59, 0, time.UTC)
			endDate = endDate.AddDate(0, 0, 6) // Add 6 days to the end date
			startEnd = append(startEnd, startDate, endDate)

		case "month":
			daysInMonth := time.Date(calc.Year, time.Month(calc.Month)+1, 0, 23, 59, 59, 0, time.UTC).Day()
			startDate := time.Date(calc.Year, time.Month(calc.Month), 1, calc.Hour, calc.Minute, 0, 0, time.UTC)
			endDate := time.Date(calc.Year, time.Month(calc.Month), daysInMonth, calc.HourTo, calc.MinuteTo, 59, 0, time.UTC)
			startEnd = append(startEnd, startDate, endDate)
		case "quarter":
			// it is responsibility of calling code, that all fields are properly set in the structure
			startDate := time.Date(calc.Year, time.Month(calc.Month), calc.Day, calc.Hour, calc.Minute, 0, 0, time.UTC)
			endDate := time.Date(calc.YearTo, time.Month(calc.MonthTo), calc.DayTo, calc.HourTo, calc.MinuteTo, 59, 0, time.UTC)
			startEnd = append(startEnd, startDate, endDate)

		case "year":
			startDate := time.Date(calc.Year, 1, 1, calc.Hour, calc.Minute, 0, 0, time.UTC)
			endDate := time.Date(calc.Year, 12, 31, calc.HourTo, calc.MinuteTo, 59, 0, time.UTC)
			startEnd = append(startEnd, startDate, endDate)

		case "custom":
			startDate := time.Date(calc.Year, time.Month(calc.Month), calc.Day, calc.Hour, calc.Minute, 0, 0, time.UTC)
			endDate := time.Date(calc.YearTo, time.Month(calc.MonthTo), calc.DayTo, calc.HourTo, calc.MinuteTo, 59, 0, time.UTC)
			startEnd = append(startEnd, startDate, endDate)
		}
	}

	// Adjust the time range based on the offset
	if len(startEnd) > 0 {
		if startEnd[1].Before(startEnd[0]) {
			startEnd[1] = startEnd[1].Add(time.Hour * 24) // Add one day if end is before start
		}

		// Apply offset (in hours) to both start and end times
		startEnd[0] = startEnd[0].Add(time.Duration(calc.Offset) * time.Hour)
		startEnd[1] = startEnd[1].Add(time.Duration(calc.OffsetTo) * time.Hour)

		return startEnd, nil
	}

	return nil, fmt.Errorf("invalid period or time range")
}

func (self *DtTimeRange) GetStartEndMilis() []int64 {
	if self.Period != "all" {
		tmp, _ := self.CalculateStartEnd()
		return []int64{tmp[0].UnixMilli(), tmp[1].UnixMilli()}
	} else {
		return []int64{0, MaxInt64()}

	}
}

func (self *DtTimeRange) StartTime_Pretty() string {
	startEnd, _ := self.CalculateStartEnd()
	start := startEnd[0].Add(time.Hour * time.Duration(-self.Offset))
	return start.Format(NOW_FORMAT)
}

func (self *DtTimeRange) EndTime_Pretty() string {
	startEnd, _ := self.CalculateStartEnd()
	end := startEnd[1].Add(time.Hour * time.Duration(-self.OffsetTo))
	return end.Format(NOW_FORMAT)
}

func (self *DtTimeRange) GetDurationMilis() int64 {
	tmp := self.GetStartEndMilis()
	return tmp[1] - tmp[0]
}

func (self *DtTimeRange) SetOffsetsUsingLocalDst() {
	/*
		Warning! Server's and client's local dst timezones are not necessarily the same!
	*/
	d := time.Date(self.Year, time.Month(self.Month), self.Day, self.Hour, self.Minute, self.Second, 0, time.Local)
	dTo := time.Date(self.YearTo, time.Month(self.MonthTo), self.DayTo, self.HourTo, self.MinuteTo, self.SecondTo, 0, time.Local)
	_, offset := d.Zone()
	_, offsetTo := dTo.Zone()
	const hoursToSec = 3600
	self.Offset = -offset / hoursToSec
	self.OffsetTo = -offsetTo / hoursToSec
	// on rare occasions around shifting of summer and winter time, this also adjusts hour, so we set it again
	self.Hour = d.Hour()
	self.HourTo = dTo.Hour()
}

// FromJSON populates the fields from a JSON object
func (self *DtTimeRange) FromJSON(jObj map[string]interface{}) {

	if v, ok := jObj["period"].(string); ok {
		self.Period = v
	}
	if v, ok := jObj["time_type"].(string); ok {
		self.TimeType = v
	}
	if v, ok := jObj["offset"].(float64); ok {
		self.Offset = int(v)
	}
	if v, ok := jObj["offset_to"].(float64); ok {
		self.OffsetTo = int(v)
	}
	if v, ok := jObj["year"].(float64); ok {
		self.Year = int(v)
	}
	if v, ok := jObj["year_to"].(float64); ok {
		self.YearTo = int(v)
	}
	if v, ok := jObj["month"].(float64); ok {
		self.Month = int(v)
	}
	if v, ok := jObj["month_to"].(float64); ok {
		self.MonthTo = int(v)
	}
	if v, ok := jObj["week"].(float64); ok {
		self.Week = int(v)
	}
	if v, ok := jObj["day"].(float64); ok {
		self.Day = int(v)
	}
	if v, ok := jObj["hour"].(float64); ok {
		self.Hour = int(v)
	}
	if v, ok := jObj["minute"].(float64); ok {
		self.Minute = int(v)
	}
	if v, ok := jObj["day_to"].(float64); ok {
		self.DayTo = int(v)
	}
	if v, ok := jObj["hour_to"].(float64); ok {
		self.HourTo = int(v)
	}
	if v, ok := jObj["minute_to"].(float64); ok {
		self.MinuteTo = int(v)
	}
	if v, ok := jObj["quarter"].(float64); ok {
		self.Quarter = int(v)
	}
	if v, ok := jObj["second"].(float64); ok {
		self.Second = int(v)
	}
	if v, ok := jObj["second_to"].(float64); ok {
		self.SecondTo = int(v)
	}
}

// ToJSON converts the structure to JSON format
func (self *DtTimeRange) ToJSON(setLocalOffsets bool) map[string]interface{} {
	if setLocalOffsets {
		self.SetOffsetsUsingLocalDst()
	}
	return map[string]interface{}{
		"period":    self.Period,
		"time_type": self.TimeType,
		"offset":    self.Offset,
		"offset_to": self.OffsetTo,
		"year":      self.Year,
		"year_to":   self.YearTo,
		"month":     self.Month,
		"month_to":  self.MonthTo,
		"week":      self.Week,
		"day":       self.Day,
		"day_to":    self.DayTo,
		"quarter":   self.Quarter,
		"hour":      self.Hour,
		"minute":    self.Minute,
		"hour_to":   self.HourTo,
		"minute_to": self.MinuteTo,
		"second":    self.Second,
		"second_to": self.SecondTo,
	}
}

func (self *DtTimeRange) Clone() DtTimeRange {
	tmp := NewDtTimeRange()
	tmp.FromJSON(self.ToJSON(true))
	return *tmp
}

func (self *DtTimeRange) CreateMongoTimeQueries() []map[string]interface{} {
	var ranges []map[string]interface{}

	if self.Period == "all" {
		return []map[string]interface{}{
			{
				"$gte": 0,
				"$lte": math.MaxInt64,
			},
		}
	} else if self.TimeType == "day" {
		startEnd, _ := self.CalculateStartEnd()

		startDate := startEnd[0]
		endDate := startEnd[1]

		ranges = append(ranges, map[string]interface{}{
			"$gte": startDate.UnixNano() / int64(time.Millisecond),
			"$lte": endDate.UnixNano() / int64(time.Millisecond),
		})
	} else {
		startEnd, _ := self.CalculateStartEnd()

		// Assuming CalculateRanges generates multiple ranges based on the period
		moreRanges := self.CalculateRanges(startEnd[0], startEnd[1])
		ranges = append(ranges, moreRanges...)
	}
	return ranges
}
func (self *DtTimeRange) CalculateRanges(startDate, endDate time.Time) []map[string]interface{} {
	var ranges []map[string]interface{}

	startTime := startDate
	endTime := endDate

	for startTime.Before(endDate) {
		ranges = append(ranges, map[string]interface{}{
			"$gte": startTime.UnixNano() / int64(time.Millisecond),
			"$lte": endTime.UnixNano() / int64(time.Millisecond),
		})

		// Increment by one day
		startTime = startTime.Add(24 * time.Hour)
		endTime = endTime.Add(24 * time.Hour)
	}

	return ranges
}

// CreateMongoTimeQuery creates a MongoDB time query based on the period and time type
func (self *DtTimeRange) CreateMongoTimeQuery() map[string]interface{} {
	//TODO sem mato treba dorobit aby sa generovalo mongo query v $IN

	query := make(map[string]interface{})

	startEnd := self.GetStartEndMilis()

	query["$gte"] = startEnd[0]
	query["$lte"] = startEnd[1]
	//lgr.Info(fmt.Sprintf("Query: from %v to %v", time.UnixMilli(startEnd[0]), time.UnixMilli(startEnd[1])))
	return query
}

// getDaysInMonth returns the number of days in a month
func getDaysInMonth(year, month int) int {
	return time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()
}

// GetWeek returns the ISO week number
func GetWeek(t time.Time) int {
	_, week := t.ISOWeek()
	return week
}
