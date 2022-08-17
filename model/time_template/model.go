package time_template

import (
	"database/sql/driver"
	"new_command/model/time_data"
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
	ID         int                 `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string              `json:"name" gorm:"unique;not null" binding:"required"`
	UpdatedAt  time.Time           `json:"updated_at" gorm:"column:updated_at;default:null"`
	CreatedAt  time.Time           `json:"created_at" gorm:"column:created_at;default:null"`
	TimeDataID int                 `json:"time_data_id" gorm:"column:time_data_id"`
	TimeData   *time_data.TimeData `json:"time_data" gorm:"foreignKey:TimeDataID" binding:"required"`
}

func (TimeTemplate) TableName() string {
	return "time_templates"
}

func (tt *TimeTemplate) UpdateData(ttp TimeTemplatePatch) {
	controller := util.NewController()
	controller.Add(2)
	go util.GoFunction(controller, tt.updateName, ttp.Name)
	go util.GoFunction(controller, tt.TimeData.UpdateData, *ttp.TimeData)
	controller.Wait()
}

func (tt *TimeTemplate) updateName(name string) {
	if name != "" {
		tt.Name = name
	}
}

type TimeTemplatePatch struct {
	Name     string                   `json:"name"`
	TimeData *time_data.TimeDataPatch `json:"time_data"`
}

type SwaggerResponse struct {
	ID         int                       `json:"id" binding:"required" example:"1"`
	Name       string                    `json:"name" binding:"required"`
	UpdatedAt  time.Time                 `json:"updated_at" binding:"required" example:"2022-08-15T06:36:36Z"`
	CreatedAt  time.Time                 `json:"created_at" binding:"required" example:"2022-08-15T06:36:36Z"`
	TimeDataID int                       `json:"time_data_id" binding:"required" example:"1"`
	TimeData   time_data.SwaggerResponse `json:"time_data" binding:"required"`
}

type SwaggerCreate struct {
	Name     string                  `json:"name" binding:"required"`
	TimeData time_data.SwaggerCreate `json:"time_data" binding:"required"`
}

type SwaggerUpdate struct {
	Name     string                  `json:"name"`
	TimeData time_data.SwaggerUpdate `json:"time_data"`
}
