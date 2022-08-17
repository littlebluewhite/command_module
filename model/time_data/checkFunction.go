package time_data

import (
	"encoding/json"
	"fmt"
	"new_command/util"
	"strings"
	"time"
)

func (td *TimeData) CheckTimeData() (err error) {
	err = nil
	ch := make(chan error, 3)
	defer close(ch)
	go func(td *TimeData, ch chan error) {
		ch <- td.checkRepeatType()
	}(td, ch)
	go func(td *TimeData, ch chan error) {
		ch <- td.checkTime()
	}(td, ch)
	go func(td *TimeData, ch chan error) {
		ch <- td.checkDate()
	}(td, ch)
	for i := 0; i < 3; i++ {
		select {
		case e := <-ch:
			if e != nil {
				err = e
				fmt.Println(err)
			}
		}
	}
	return
}

func (td *TimeData) checkRepeatType() (err error) {
	err = nil
	switch td.RepeatType {
	case daily:
		if td.ConditionType != "" {

			err = &util.MyError{
				When: time.Now(),
				What: "doesn't need condition type",
			}
			return
		}
		if td.Condition != nil {
			err = &util.MyError{
				When: time.Now(),
				What: "doesn't need condition",
			}
			return
		}
	case weekly:
		if td.ConditionType != weeklyDay {
			err = &util.MyError{
				When: time.Now(),
				What: "condition type error",
			}
			return
		}
		if err = td.checkWeeklyCondition(); err != nil {
			return err
		}
	case monthly:
		if td.ConditionType == weeklyDay || td.ConditionType == "" {
			err = &util.MyError{
				When: time.Now(),
				What: "condition type error",
			}
			return
		}
		if err = td.checkMonthlyCondition(); err != nil {
			return err
		}
	default:
		err = &util.MyError{
			When: time.Now(),
			What: "repeat_type error",
		}
	}
	return
}

func (td *TimeData) checkWeeklyCondition() (err error) {
	err = nil
	if td.Condition == nil {
		err = &util.MyError{
			When: time.Now(),
			What: "can't find condition",
		}
		return
	}
	var conditions []int
	err = json.Unmarshal(td.Condition, &conditions)
	if err != nil {
		err = &util.MyError{
			When: time.Now(),
			What: "condition format are not correct",
		}
		return
	}
	if ok := util.Contains(conditions, allWeekDay); !ok {
		err = &util.MyError{
			When: time.Now(),
			What: "weekly condition number are not correct",
		}
		return
	}
	return
}

func (td *TimeData) checkMonthlyCondition() (err error) {
	err = nil
	if td.Condition == nil {
		err = &util.MyError{
			When: time.Now(),
			What: "can't find condition",
		}
		return
	}
	var conditions []int
	err = json.Unmarshal(td.Condition, &conditions)
	if err != nil {
		err = &util.MyError{
			When: time.Now(),
			What: "condition format are not correct",
		}
		return
	}
	monthType := strings.Split(string(td.ConditionType), "_")
	switch monthType[0] {
	case "weekly":
		if ok := util.Contains(conditions, allWeekDay); !ok {
			err = &util.MyError{
				When: time.Now(),
				What: "monthly condition number are not correct",
			}
			return
		}
	case "monthly":
		if ok := util.Contains(conditions, allMonthDay); !ok {
			err = &util.MyError{
				When: time.Now(),
				What: "monthly condition number are not correct",
			}
			return
		}
	}
	return
}

func (td *TimeData) checkTime() (err error) {
	err = nil
	if td.StartTime > td.EndTime {
		err = &util.MyError{
			When: time.Now(),
			What: "start time and end time error",
		}
	}
	return
}

func (td *TimeData) checkDate() (err error) {
	err = nil
	startDate := time.Time(td.StartDate).Unix()
	endDate := time.Time(td.EndDate).Unix()
	if startDate > endDate {
		err = &util.MyError{
			When: time.Now(),
			What: "start date and end date error",
		}
	}
	if endDate < 100 {
		err = nil
	}
	return
}
