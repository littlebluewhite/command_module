package schedule

import (
	"new_command/model/command"
	"new_command/model/time_data"
	"new_command/util"
	"time"
)

type Schedule struct {
	ID          int                 `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string              `json:"name" gorm:"unique;not null" binding:"required"`
	Description string              `json:"description" gorm:"column:description"`
	TimeDataID  int                 `json:"time_data_id" gorm:"column:time_data_id"`
	CommandID   int                 `json:"command_id" gorm:"column:command_id" binding:"required"`
	Enabled     bool                `json:"enabled" gorm:"column:enabled;default:false"`
	UpdatedAt   time.Time           `json:"updated_at" gorm:"column:updated_at;default:null"`
	CreatedAt   time.Time           `json:"created_at" gorm:"column:created_at;default:null"`
	TimeData    *time_data.TimeData `json:"time_data" gorm:"foreignKey:TimeDataID" binding:"required"`
	Command     *command.Command    `json:"command" gorm:"foreignKey:CommandID"`
}

func (*Schedule) TableName() string {
	return "schedules"
}

func (s *Schedule) UpdateData(sp SchedulePatch) {
	controller := util.NewController()
	controller.Add(4)
	go util.GoFunction(controller, s.updateName, sp.Name)
	go util.GoFunction(controller, s.updateDescription, sp.Description)
	go util.GoFunction(controller, s.updateCommandID, sp.CommandID)
	go util.GoFunction(controller, s.updateEnabled, sp.Enabled)
	if sp.TimeData != nil {
		controller.Add(1)
		go util.GoFunction(controller, s.TimeData.UpdateData, *sp.TimeData)
	}
	controller.Wait()
}

func (s *Schedule) updateName(name string) {
	if name != "" {
		s.Name = name
	}
}

func (s *Schedule) updateDescription(description string) {
	if description != "" {
		s.Description = description
	}
}

func (s *Schedule) updateCommandID(commandID int) {
	if commandID != 0 {
		s.CommandID = commandID
	}
}

func (s *Schedule) updateEnabled(enabled *bool) {
	if enabled != nil {
		s.Enabled = *enabled
	}
}

type SchedulePatch struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	CommandID   int                      `json:"command_id"`
	Enabled     *bool                    `json:"enabled"`
	TimeData    *time_data.TimeDataPatch `json:"time_data"`
}

type SwaggerResponse struct {
	ID          int                       `json:"id" binding:"required" example:"1"`
	Name        string                    `json:"name" binding:"required" example:"test"`
	Description string                    `json:"description" binding:"required" example:"describe something"`
	TimeDataID  int                       `json:"time_data_id" binding:"required" example:"1"`
	CommandID   int                       `json:"command_id" binding:"required" example:"1"`
	Enabled     bool                      `json:"enabled" binding:"required" example:"true"`
	UpdatedAt   time.Time                 `json:"updated_at" binding:"required" example:"2022-08-15T06:36:36Z"`
	CreatedAt   time.Time                 `json:"created_at" binding:"required" example:"2022-08-15T06:36:36Z"`
	TimeData    time_data.SwaggerResponse `json:"time_data" binding:"required"`
	Command     command.SwaggerResponse   `json:"command" binding:"required"`
}

type SwaggerCreate struct {
	Name        string                  `json:"name" binding:"required" example:"test"`
	Description string                  `json:"description" binding:"required" example:"describe something"`
	CommandID   int                     `json:"command_id" binding:"required" example:"1"`
	Enabled     bool                    `json:"enabled" binding:"required" example:"true"`
	TimeData    time_data.SwaggerCreate `json:"time_data" binding:"required"`
}

type SwaggerUpdate struct {
	Name        string                  `json:"name" example:"test"`
	Description string                  `json:"description" example:"describe something"`
	CommandID   int                     `json:"command_id" example:"1"`
	Enabled     bool                    `json:"enabled" example:"true"`
	TimeData    time_data.SwaggerUpdate `json:"time_data"`
}
