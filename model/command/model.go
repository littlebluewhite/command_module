package command

import (
	"database/sql/driver"
	"new_command/model/https_command"
	"new_command/util"
	"time"
)

type Protocol string

const (
	Http      Protocol = "http"
	Websocket Protocol = "websocket"
	Socket    Protocol = "socket"
	Snmp      Protocol = "snmp"
)

func (p *Protocol) Scan(value interface{}) error {
	*p = Protocol(value.([]byte))
	return nil
}

func (p Protocol) Value() (driver.Value, error) {
	return string(p), nil
}

type Command struct {
	ID           int                         `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string                      `json:"name" gorm:"unique;not null" binding:"required"`
	Description  string                      `json:"description,omitempty" gorm:"column:description"`
	Protocol     Protocol                    `json:"protocol" gorm:"type:enum('http','socket','websocket','snmp');default:null" binding:"required"`
	UpdatedAt    time.Time                   `json:"updated_at" gorm:"column:updated_at;default:null"`
	CreatedAt    time.Time                   `json:"created_at" gorm:"column:created_at;default:null"`
	HttpsCommand *https_command.HttpsCommand `json:"https_command,omitempty"`
}

func (*Command) TableName() string {
	return "commands"
}

func (c *Command) UpdateData(cp CommandPatch) {
	controller := util.NewController()
	controller.Add(3)
	go util.GoFunction(controller, c.updateName, cp.Name)
	go util.GoFunction(controller, c.updateDescription, cp.Description)
	go util.GoFunction(controller, c.HttpsCommand.UpdateData, *cp.HttpsCommand)
	controller.Wait()
}

func (c *Command) updateName(name string) {
	if name != "" {
		c.Name = name
	}
}

func (c *Command) updateDescription(description string) {
	if description != "" {
		c.Description = description
	}
}

type CommandPatch struct {
	Name         string                           `json:"name"`
	Description  string                           `json:"description"`
	HttpsCommand *https_command.HttpsCommandPatch `json:"https_command"`
}

type SwaggerResponse struct {
	ID           int                           `json:"id" binding:"required" example:"1"`
	Name         string                        `json:"name" binding:"required" example:"test"`
	Description  string                        `json:"description" binding:"required" example:"describe something"`
	Protocol     Protocol                      `json:"protocol" binding:"required" example:"http"`
	UpdatedAt    time.Time                     `json:"updated_at" binding:"required" example:"2022-08-15T06:36:36Z"`
	CreatedAt    time.Time                     `json:"created_at" binding:"required" example:"2022-08-15T06:36:36Z"`
	HttpsCommand https_command.SwaggerResponse `json:"https_command" binding:"required"`
}

type SwaggerCreate struct {
	Name         string                      `json:"name" binding:"required" example:"test"`
	Description  string                      `json:"description" binding:"required" example:"describe something"`
	Protocol     Protocol                    `json:"protocol" binding:"required" example:"http"`
	HttpsCommand https_command.SwaggerCreate `json:"https_command" binding:"required"`
}

type SwaggerUpdate struct {
	Name         string                      `json:"name" example:"test"`
	Description  string                      `json:"description" example:"describe something"`
	Protocol     Protocol                    `json:"protocol" example:"http"`
	HttpsCommand https_command.SwaggerUpdate `json:"https_command"`
}
