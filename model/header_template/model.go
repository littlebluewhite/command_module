package header_template

import (
	"gorm.io/datatypes"
	"new_command/util"
)

type HeaderTemplate struct {
	ID   int            `json:"id" gorm:"primaryKey;autoIncrement" example:"1"`
	Name string         `json:"name" gorm:"unique;not null" binding:"required"`
	Data datatypes.JSON `json:"data" gorm:"column:data" binding:"required"`
}

func (*HeaderTemplate) TableName() string {
	return "header_template"
}

func (ht *HeaderTemplate) UpdateData(htp HeaderTemplatePatch) {
	controller := util.NewController()
	controller.Add(2)
	go util.GoFunction(controller, ht.updateName, htp.Name)
	go util.GoFunction(controller, ht.updateData, htp.Data)
}

func (ht *HeaderTemplate) updateName(name string) {
	if name != "" {
		ht.Name = name
	}
}

func (ht *HeaderTemplate) updateData(data datatypes.JSON) {
	if data != nil {
		ht.Data = data
	}
}

type HeaderTemplatePatch struct {
	Name string         `json:"name"`
	Data datatypes.JSON `json:"data"`
}

type SwaggerResponse struct {
	ID   int      `json:"id" example:"1" binding:"required"`
	Name string   `json:"name" binding:"required"`
	Data struct{} `json:"data" binding:"required"`
}

type SwaggerCreate struct {
	Name string   `json:"name" binding:"required"`
	Data struct{} `json:"data" binding:"required"`
}

type SwaggerUpdate struct {
	Name string   `json:"name"`
	Data struct{} `json:"data"`
}
