package time_template

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/datatypes"
	"new_command/util"
	"time"
)

type RepeatType string

const (
	daily   RepeatType = "daily"
	weekly  RepeatType = "weekly"
	monthly RepeatType = "monthly"
)

func (rt *RepeatType) Scan(value interface{}) error {
	*rt = RepeatType(value.([]byte))
	return nil
}

func (rt RepeatType) Value() (driver.Value, error) {
	return string(rt), nil
}

type TimeTemplate struct {
	ID              int            `json:"id" gorm:"primaryKey;autoIncrement"`
	Name            string         `json:"name" gorm:"unique;not null" binding:"required"`
	RepeatType      RepeatType     `json:"repeat_type" gorm:"type:enum('daily','weekly','monthly');default:daily" binding:"required"`
	StartDate       datatypes.Date `json:"start_date" gorm:"column:start_date;not null" binding:"required"`
	EndDate         datatypes.Date `json:"end_date,omitempty" gorm:"column:end_date;default:3000-01-01T00:00:00+00:00"`
	StartTime       datatypes.Time `json:"start_time" gorm:"column:start_time;not null" binding:"required"`
	EndTime         datatypes.Time `json:"end_time" gorm:"column:end_time;not null" binding:"required"`
	IntervalSeconds int            `json:"interval_seconds,omitempty" gorm:"column:interval_seconds"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"column:updated_at;default:null"`
	CreatedAt       time.Time      `json:"created_at" gorm:"column:created_at;default:null"`
	WeeklyRepeat    *WeeklyRepeat  `json:"weekly_repeat,omitempty" gorm:"foreignkey:TimeTemplateID"`
	MonthlyRepeat   *MonthlyRepeat `json:"monthly_repeat,omitempty" gorm:"foreignkey:TimeTemplateID"`
}

func (TimeTemplate) TableName() string {
	return "time_templates"
}

type WeeklyRepeat struct {
	ID              int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	TimeTemplateID  uint           `json:"time_template_id" gorm:"column:time_template_id;uniqueIndex;not null"`
	WeeklyCondition datatypes.JSON `json:"weekly_condition" gorm:"column:weekly_condition" binding:"required"`
}

var weekDay = []int{0, 1, 2, 3, 4, 5, 6}
var monthDay = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
	21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}

type MonthlyRepeat struct {
	ID                    int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	TimeTemplateID        uint           `json:"time_template_id" gorm:"column:time_template_id;uniqueIndex;not null"`
	FirstWeeklyCondition  datatypes.JSON `json:"first_weekly_condition" gorm:"column:first_weekly_condition" binding:"required"`
	SecondWeeklyCondition datatypes.JSON `json:"second_weekly_condition" gorm:"column:second_weekly_condition" binding:"required"`
	ThirdWeeklyCondition  datatypes.JSON `json:"third_weekly_condition" gorm:"column:third_weekly_condition" binding:"required"`
	FourthWeeklyCondition datatypes.JSON `json:"fourth_weekly_condition" gorm:"column:fourth_weekly_condition" binding:"required"`
	MonthlyCondition      datatypes.JSON `json:"monthly_condition" gorm:"column:monthly_condition" binding:"required"`
}

func (tt *TimeTemplate) CheckRepeatType() (err error) {
	err = nil
	switch tt.RepeatType {
	case daily:
		if tt.MonthlyRepeat != nil {

			err = &util.MyError{
				When: time.Now(),
				What: "doesn't need monthly conditions",
			}
			return
		}
		if tt.WeeklyRepeat != nil {
			err = &util.MyError{
				When: time.Now(),
				What: "doesn't need weekly conditions",
			}
			return
		}
	case weekly:
		if tt.WeeklyRepeat == nil {
			err = &util.MyError{
				When: time.Now(),
				What: "can't find weekly condition",
			}
			return
		}
		if tt.MonthlyRepeat != nil {
			err = &util.MyError{
				When: time.Now(),
				What: "doesn't need monthly conditions",
			}
			return
		}
		if err = tt.checkWeeklyRepeat(); err != nil {
			return err
		}
	case monthly:
		if tt.MonthlyRepeat == nil {
			err = &util.MyError{
				When: time.Now(),
				What: "can't find monthly condition",
			}
			return
		}
		if tt.WeeklyRepeat != nil {
			err = &util.MyError{
				When: time.Now(),
				What: "doesn't need weekly conditions",
			}
			return
		}
	default:
		err = &util.MyError{
			When: time.Now(),
			What: "repeat_type error",
		}
	}
	return
}

func (tt *TimeTemplate) CheckTime() (err error) {
	err = nil
	if tt.StartTime > tt.EndTime {
		err = &util.MyError{
			When: time.Now(),
			What: "start time and end time error",
		}
	}
	return
}

func (tt *TimeTemplate) CheckDate() (err error) {
	err = nil
	startDate := time.Time(tt.StartDate).Unix()
	endDate := time.Time(tt.EndDate).Unix()
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

func (tt *TimeTemplate) checkWeeklyRepeat() (err error) {
	err = nil
	fmt.Println(tt.WeeklyRepeat.WeeklyCondition)
	var conditions []int
	err = json.Unmarshal(tt.WeeklyRepeat.WeeklyCondition, &conditions)
	fmt.Println(err)
	fmt.Println(conditions)
	if err != nil {
		err = &util.MyError{
			When: time.Now(),
			What: "weekly conditions are not correct",
		}
		return
	}
	if ok := util.Contains(conditions, weekDay); !ok {
		err = &util.MyError{
			When: time.Now(),
			What: "weekly conditions number are not correct",
		}
		return
	}
	return
}

type TimeTemplatePatch struct {
	Name            string         `json:"name"`
	RepeatType      RepeatType     `json:"repeat_type"`
	StartDate       datatypes.Date `json:"start_date"`
	EndDate         datatypes.Date `json:"end_date,omitempty"`
	StartTime       datatypes.Time `json:"start_time"`
	EndTime         datatypes.Time `json:"end_time"`
	IntervalSeconds int            `json:"interval_seconds,omitempty"`
	WeeklyRepeat    *WeeklyRepeat  `json:"weekly_repeat,omitempty"`
	MonthlyRepeat   *MonthlyRepeat `json:"monthly_repeat,omitempty"`
}

func (ttp *TimeTemplatePatch) ToModel() (tt TimeTemplate) {
	tt = TimeTemplate{
		Name:            ttp.Name,
		RepeatType:      ttp.RepeatType,
		StartDate:       ttp.StartDate,
		EndDate:         ttp.EndDate,
		StartTime:       ttp.StartTime,
		EndTime:         ttp.EndTime,
		IntervalSeconds: ttp.IntervalSeconds,
		WeeklyRepeat:    ttp.WeeklyRepeat,
		MonthlyRepeat:   ttp.MonthlyRepeat,
	}
	return
}

func (tt *TimeTemplate) ConvertToUpdateMap() (updateData map[string]interface{}) {

}
