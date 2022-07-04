package time_template

import (
	"database/sql/driver"
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
	ID              int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name            string         `json:"name" gorm:"unique;not null" binding:"required"`
	RepeatType      RepeatType     `json:"repeat_type" gorm:"type:enum('daily','weekly','monthly');default:daily" binding:"required"`
	StartDate       datatypes.Date `json:"start_date" gorm:"column:start_date;not null" binding:"required"`
	EndDate         datatypes.Date `json:"end_date" gorm:"column:end_date;default:null"`
	StartTime       datatypes.Time `json:"start_time" gorm:"column:start_time;not null" binding:"required"`
	EndTime         datatypes.Time `json:"end_time" gorm:"column:end_time;not null" binding:"required"`
	IntervalSeconds int64          `json:"interval_seconds" gorm:"column:interval_seconds;default:null"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"column:updated_at;default:null"`
	CreatedAt       time.Time      `json:"created_at" gorm:"column:created_at;default:null"`
	WeeklyRepeat    WeeklyRepeat   `json:"weekly_repeat" gorm:"foreignkey:TimeTemplateID"`
	MonthlyRepeat   MonthlyRepeat  `json:"monthly_repeat" gorm:"foreignkey:TimeTemplateID"`
}

func (TimeTemplate) TableName() string {
	return "time_templates"
}

type WeeklyRepeat struct {
	ID              int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	TimeTemplateID  uint           `json:"time_template_id" gorm:"column:time_template_id;uniqueIndex;not null"`
	WeeklyCondition datatypes.JSON `json:"weekly_condition" gorm:"column:weekly_condition"`
}

type MonthlyRepeat struct {
	ID                    int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	TimeTemplateID        uint           `json:"time_template_id" gorm:"column:time_template_id;uniqueIndex;not null"`
	FirstWeeklyCondition  datatypes.JSON `json:"first_weekly_condition" gorm:"column:first_weekly_condition"`
	SecondWeeklyCondition datatypes.JSON `json:"second_weekly_condition" gorm:"column:second_weekly_condition"`
	ThirdWeeklyCondition  datatypes.JSON `json:"third_weekly_condition" gorm:"column:third_weekly_condition"`
	FourthWeeklyCondition datatypes.JSON `json:"fifth_weekly_condition" gorm:"column:fifth_weekly_condition"`
	MonthlyCondition      datatypes.JSON `json:"monthly_condition" gorm:"column:monthly_condition"`
}

func (tt *TimeTemplate) CheckRepeatType() (err error) {
	err = nil
	switch tt.RepeatType {
	case daily:
	case weekly:
		if tt.WeeklyRepeat.WeeklyCondition == nil {
			err = &util.MyError{
				When: time.Now(),
				What: "can't find weekly condition",
			}
		}
	case monthly:
		if conditions := tt.MonthlyRepeat; conditions.MonthlyCondition == nil || conditions.FirstWeeklyCondition == nil || conditions.SecondWeeklyCondition == nil || conditions.ThirdWeeklyCondition == nil || conditions.FourthWeeklyCondition == nil {
			err = &util.MyError{
				When: time.Now(),
				What: "can't find monthly condition",
			}
		}
	}
	return
}
