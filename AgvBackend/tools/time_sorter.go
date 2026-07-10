package tools

import (
	"time"
)

type WeeksSorter struct {
	Weeks map[int][]interface{}
}

func NewWeeksSorter() *WeeksSorter {
	return &WeeksSorter{
		Weeks: make(map[int][]interface{}),
	}
}

func (ws *WeeksSorter) Add(at int64, element interface{}) {
	t := time.UnixMilli(at)
	_, week := t.ISOWeek()

	if _, exists := ws.Weeks[week]; exists {
		ws.Weeks[week] = append(ws.Weeks[week], element)
	} else {
		ws.Weeks[week] = []interface{}{element}
	}
}

func (ws *WeeksSorter) GetWeeksSorted() [][]interface{} {
	// Implement sorting logic if needed
	return [][]interface{}{}
}

type WeekDaySorter struct {
	Days map[int][]interface{}
}

func NewWeekDaySorter() *WeekDaySorter {
	return &WeekDaySorter{
		Days: make(map[int][]interface{}),
	}
}

func (wds *WeekDaySorter) Add(at int64, element interface{}) {
	t := time.UnixMilli(at)
	day := int(t.Weekday())

	if _, exists := wds.Days[day]; exists {
		wds.Days[day] = append(wds.Days[day], element)
	} else {
		wds.Days[day] = []interface{}{element}
	}
}

func (wds *WeekDaySorter) GetDaysSorted() [][]interface{} {
	// Implement sorting logic if needed
	return [][]interface{}{}
}

type DaySorter struct {
	Days map[int][]interface{}
}

func NewDaySorter() *DaySorter {
	return &DaySorter{
		Days: make(map[int][]interface{}),
	}
}

func (ds *DaySorter) Add(at int64, element interface{}) {
	t := time.UnixMilli(at)
	day := int(t.Weekday())

	if _, exists := ds.Days[day]; exists {
		ds.Days[day] = append(ds.Days[day], element)
	} else {
		ds.Days[day] = []interface{}{element}
	}
}

func (ds *DaySorter) GetDaysSorted() [][]interface{} {
	// Implement sorting logic if needed
	return [][]interface{}{}
}
