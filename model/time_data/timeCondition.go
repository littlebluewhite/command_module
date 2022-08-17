package time_data

import (
	"encoding/json"
	"new_command/util"
	"time"
)

func (td *TimeData) AccessTimeDate(t time.Time) (result bool) {
	result = true
	ch := make(chan bool, 3)
	go func(td *TimeData, t time.Time, ch chan bool) {
		ch <- td.AccessTime(t)
	}(td, t, ch)
	go func(td *TimeData, t time.Time, ch chan bool) {
		ch <- td.AccessDate(t)
	}(td, t, ch)
	go func(td *TimeData, t time.Time, ch chan bool) {
		ch <- td.AccessConditions(t)
	}(td, t, ch)
	for i := 0; i < 3; i++ {
		select {
		case b := <-ch:
			if b == false {
				result = false
			}
		}
	}
	return
}

func (td *TimeData) AccessDate(t time.Time) (result bool) {
	if time.Time(td.StartDate).Unix() <= t.Unix() && time.Time(td.EndDate).Add(24*time.Hour).Unix() > t.Unix() {
		result = true
	}
	return
}

func (td *TimeData) AccessTime(t time.Time) (result bool) {
	startInt := int(td.StartTime)
	nowInt := util.GetTimeInt(t)
	endInt := int(td.EndTime)
	if !(startInt <= nowInt && endInt+999999999 >= nowInt) {
		return
	}
	if td.IntervalSeconds == 0 {
		result = true
		return
	}
	if duration := nowInt - startInt; (duration/int(time.Second))%td.IntervalSeconds == 0 {
		result = true
	}
	return
}

func (td *TimeData) AccessConditions(t time.Time) (result bool) {
	switch td.RepeatType {
	case daily:
		result = true
	case weekly:
		result = td.AccessWeekly(t)
	case monthly:
		result = td.AccessMonthly(t)
	}
	return
}

func (td *TimeData) AccessWeekly(t time.Time) (result bool) {
	if td.ConditionType != weeklyDay {
		return
	}
	var conditions []int
	if err := json.Unmarshal(td.Condition, &conditions); err != nil {
		return
	}
	result = util.Contains[int]([]int{int(t.Weekday())}, conditions)
	return
}

func (td *TimeData) AccessMonthly(t time.Time) (result bool) {
	var conditions []int
	if err := json.Unmarshal(td.Condition, &conditions); err != nil {
		return
	}
	weekCount := util.CountWeek(t)
	switch td.ConditionType {
	case monthDay:
		result = util.Contains[int]([]int{t.Day()}, conditions)
	case weeklyFirst:
		if weekCount == 0 {
			result = util.Contains[int]([]int{int(t.Weekday())}, conditions)
		}
	case weeklySecond:
		if weekCount == 1 {
			result = util.Contains[int]([]int{int(t.Weekday())}, conditions)
		}
	case weeklyThird:
		if weekCount == 2 {
			result = util.Contains[int]([]int{int(t.Weekday())}, conditions)
		}
	case weeklyFourth:
		if weekCount == 3 {
			result = util.Contains[int]([]int{int(t.Weekday())}, conditions)
		}
	}
	return
}
