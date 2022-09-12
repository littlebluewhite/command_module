package time_data

import (
	"database/sql/driver"
	"gorm.io/datatypes"
	"new_command/util"
	"time"
)

var allWeekDay = []int{0, 1, 2, 3, 4, 5, 6}
var allMonthDay = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
	21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}

type RepeatType string

const (
	daily   RepeatType = "daily"
	weekly  RepeatType = "weekly"
	monthly RepeatType = "monthly"
)

type ConditionType string

const (
	monthDay     ConditionType = "monthly_day"
	weeklyDay    ConditionType = "weekly_day"
	weeklyFirst  ConditionType = "weekly_first"
	weeklySecond ConditionType = "weekly_second"
	weeklyThird  ConditionType = "weekly_third"
	weeklyFourth ConditionType = "weekly_fourth"
)

func (rt *RepeatType) Scan(value interface{}) error {
	*rt = RepeatType(value.([]byte))
	return nil
}

func (rt RepeatType) Value() (driver.Value, error) {
	return string(rt), nil
}

func (ct *ConditionType) Scan(value interface{}) error {
	*ct = ConditionType(value.([]byte))
	return nil
}

func (ct ConditionType) Value() (driver.Value, error) {
	return string(ct), nil
}

type TimeData struct {
	ID              int            `json:"id" gorm:"primaryKey;autoIncrement"`
	RepeatType      RepeatType     `json:"repeat_type" gorm:"type:enum('daily','weekly','monthly', '');default:''"`
	StartDate       datatypes.Date `json:"start_date" gorm:"column:start_date;not null" binding:"required"`
	EndDate         datatypes.Date `json:"end_date,omitempty" gorm:"column:end_date;defaremiult:9999-01-01T00:00:00+00:00"`
	StartTime       datatypes.Time `json:"start_time" gorm:"column:start_time;not null"`
	EndTime         datatypes.Time `json:"end_time" gorm:"column:end_time;not null" binding:"required"`
	IntervalSeconds int            `json:"interval_seconds,omitempty" gorm:"column:interval_seconds"`
	ConditionType   ConditionType  `json:"condition_type" gorm:"column:condition_type;type:enum('monthly_day','weekly_day','weekly_first','weekly_second','weekly_third','weekly_fourth','');default:''"`
	Condition       datatypes.JSON `json:"condition" gorm:"column:condition;default:null"`
}

func (*TimeData) TableName() string {
	return "time_data"
}

func (td *TimeData) UpdateData(tdp TimeDataPatch) {
	controller := util.NewController()
	controller.Add(8)
	go util.GoFunction(controller, td.updateRepeatType, tdp.RepeatType)
	go util.GoFunction(controller, td.updateStartDate, tdp.StartDate)
	go util.GoFunction(controller, td.updateEndDate, tdp.EndDate)
	go util.GoFunction(controller, td.updateStartTime, tdp.StartTime)
	go util.GoFunction(controller, td.updateEndTime, tdp.EndTime)
	go util.GoFunction(controller, td.updateIntervalSeconds, tdp.IntervalSeconds)
	go util.GoFunction(controller, td.updateConditionType, tdp.ConditionType)
	go util.GoFunction(controller, td.updateCondition, tdp.Condition)
	controller.Wait()
}

func (td *TimeData) updateRepeatType(repeatType RepeatType) {
	if repeatType != "" {
		td.RepeatType = repeatType
	}
}

func (td *TimeData) updateStartDate(startDate datatypes.Date) {
	if time.Time(startDate).Year() != 1 {
		td.StartDate = startDate
	}
}

func (td *TimeData) updateEndDate(endDate datatypes.Date) {
	if time.Time(endDate).Year() != 1 {
		td.EndDate = endDate
	}
}

func (td *TimeData) updateStartTime(startTime datatypes.Time) {
	if int(startTime) != 0 {
		td.StartTime = startTime
	}
}

func (td *TimeData) updateEndTime(endTime datatypes.Time) {
	if int(endTime) != 0 {
		td.EndTime = endTime
	}
}

func (td *TimeData) updateIntervalSeconds(intervalSeconds int) {
	if intervalSeconds != 0 {
		td.IntervalSeconds = intervalSeconds
	}
}

func (td *TimeData) updateConditionType(conditionType ConditionType) {
	if conditionType != "" {
		td.ConditionType = conditionType
	}
}

func (td *TimeData) updateCondition(condition datatypes.JSON) {
	if condition != nil {
		td.Condition = condition
	}
	return
}

type TimeDataPatch struct {
	RepeatType      RepeatType     `json:"repeat_type"`
	StartDate       datatypes.Date `json:"start_date"`
	EndDate         datatypes.Date `json:"end_date"`
	StartTime       datatypes.Time `json:"start_time"`
	EndTime         datatypes.Time `json:"end_time"`
	IntervalSeconds int            `json:"interval_seconds"`
	ConditionType   ConditionType  `json:"condition_type"`
	Condition       datatypes.JSON `json:"condition"`
}

type SwaggerResponse struct {
	ID              int            `json:"id" binding:"required" example:"1"`
	RepeatType      RepeatType     `json:"repeat_type" binding:"required" example:"weekly"`
	StartDate       datatypes.Date `json:"start_date" binding:"required" example:"2022-08-15T06:36:36Z"`
	EndDate         datatypes.Date `json:"end_date,omitempty" binding:"required" example:"2022-08-15T06:36:36Z"`
	StartTime       datatypes.Time `json:"start_time" binding:"required" example:"00:00:22"`
	EndTime         datatypes.Time `json:"end_time" binding:"required" example:"21:33:22"`
	IntervalSeconds int            `json:"interval_seconds,omitempty" binding:"required" example:"5"`
	ConditionType   ConditionType  `json:"condition_type" binding:"required" example:"weekly_day"`
	Condition       []int          `json:"condition" binding:"required"`
}

type SwaggerCreate struct {
	RepeatType      RepeatType     `json:"repeat_type" binding:"required" example:"weekly"`
	StartDate       datatypes.Date `json:"start_date" binding:"required" example:"2022-08-15T06:36:36Z"`
	EndDate         datatypes.Date `json:"end_date,omitempty" binding:"required" example:"2022-08-15T06:36:36Z"`
	StartTime       datatypes.Time `json:"start_time" example:"00:00:22"`
	EndTime         datatypes.Time `json:"end_time" example:"21:33:22"`
	IntervalSeconds int            `json:"interval_seconds" example:"5"`
	ConditionType   ConditionType  `json:"condition_type" example:"weekly_day"`
	Condition       []int          `json:"condition"`
}

type SwaggerUpdate struct {
	RepeatType      RepeatType     `json:"repeat_type" example:"weekly"`
	StartDate       datatypes.Date `json:"start_date" example:"2022-08-15T06:36:36Z"`
	EndDate         datatypes.Date `json:"end_date,omitempty" example:"2022-08-15T06:36:36Z"`
	StartTime       datatypes.Time `json:"start_time" example:"00:00:22"`
	EndTime         datatypes.Time `json:"end_time" example:"21:33:22"`
	IntervalSeconds int            `json:"interval_seconds,omitempty" example:"5"`
	ConditionType   ConditionType  `json:"condition_type" example:"weekly_day"`
	Condition       []int          `json:"condition"`
}
