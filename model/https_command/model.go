package https_command

import (
	"database/sql/driver"
	"gorm.io/datatypes"
	"new_command/pkg/logFile"
	"new_command/util"
)

var httpsCommandLog = logFile.NewLogFile("model", "httpsCommand.log")

type HttpMethod string

const (
	Get    HttpMethod = "GET"
	Post   HttpMethod = "POST"
	Put    HttpMethod = "PUT"
	Patch  HttpMethod = "PATCH"
	Delete HttpMethod = "DELETE"
)

func (hm *HttpMethod) Scan(value interface{}) error {
	*hm = HttpMethod(value.([]byte))
	return nil
}

func (hm HttpMethod) Value() (driver.Value, error) {
	return string(hm), nil
}

type AuthorizationType string

const (
	Basic AuthorizationType = "basic"
	Token AuthorizationType = "token"
)

func (at *AuthorizationType) Scan(value interface{}) error {
	*at = AuthorizationType(value.([]byte))
	return nil
}

func (at AuthorizationType) Value() (driver.Value, error) {
	return string(at), nil
}

type BodyType string

const (
	Text               BodyType = "text"
	Html               BodyType = "html"
	Xml                BodyType = "xml"
	FormData           BodyType = "form_data"
	XWwwFormUrlencoded BodyType = "x_www_form_urlencoded"
	Json               BodyType = "json"
)

func (bt *BodyType) Scan(value interface{}) error {
	*bt = BodyType(value.([]byte))
	return nil
}

func (bt BodyType) Value() (driver.Value, error) {
	return string(bt), nil
}

type HttpsCommand struct {
	ID                int               `json:"id" gorm:"primaryKey;autoIncrement"`
	Url               string            `json:"url" gorm:"column:url" binding:"required"`
	Method            HttpMethod        `json:"method" gorm:"column:method;type:enum('GET','POST','PUT','PATCH','DELETE')" binding:"required,oneof=GET POST PUT PATCH DELETE"`
	AuthorizationType AuthorizationType `json:"authorization_type,omitempty" gorm:"column:authorization_type;type:enum('basic', 'token', '');default:''"`
	Header            datatypes.JSON    `json:"header,omitempty" gorm:"column:header"`
	BodyType          BodyType          `json:"body_type" gorm:"column:body_type;type:enum('text','html','xml','form_data','x_www_form_urlencoded','json','');default:''"`
	Body              datatypes.JSON    `json:"body" gorm:"column:body;default:null"`
	CommandID         int               `json:"command_id" gorm:"column:command_id"`
}

func (*HttpsCommand) TableName() string {
	return "https_commands"
}

func (hc *HttpsCommand) UpdateData(hcp HttpsCommandPatch) {
	controller := util.NewController()
	controller.Add(6)
	go util.GoFunction(controller, hc.updateUrl, hcp.Url)
	go util.GoFunction(controller, hc.updateMethod, hcp.Method)
	go util.GoFunction(controller, hc.updateAuthorizationType, hcp.AuthorizationType)
	go util.GoFunction(controller, hc.updateHeader, hcp.Header)
	go util.GoFunction(controller, hc.updateBodyType, hcp.BodyType)
	go util.GoFunction(controller, hc.updateBody, hcp.Body)
	controller.Wait()
}

func (hc *HttpsCommand) updateUrl(url string) {
	if url != "" {
		hc.Url = url
	}
}

func (hc *HttpsCommand) updateMethod(method HttpMethod) {
	if method != "" {
		hc.Method = method
	}
}

func (hc *HttpsCommand) updateAuthorizationType(authorizationType AuthorizationType) {
	if authorizationType != "" {
		hc.AuthorizationType = authorizationType
	}
}

func (hc *HttpsCommand) updateHeader(header datatypes.JSON) {
	if header != nil {
		hc.Header = header
	}
}

func (hc *HttpsCommand) updateBodyType(bodyType BodyType) {
	if bodyType != "" {
		hc.BodyType = bodyType
	}
}

func (hc *HttpsCommand) updateBody(body datatypes.JSON) {
	if body != nil {
		hc.Body = body
	}
}

type HttpsCommandPatch struct {
	Url               string            `json:"url"`
	Method            HttpMethod        `json:"method"`
	AuthorizationType AuthorizationType `json:"authorization_type"`
	Header            datatypes.JSON    `json:"header,omitempty"`
	BodyType          BodyType          `json:"body_type"`
	Body              datatypes.JSON    `json:"body"`
}

type SwaggerResponse struct {
	ID                int               `json:"id" binding:"required" example:"1"`
	Url               string            `json:"url" binding:"required" example:"http://localhost:9800/api"`
	Method            HttpMethod        `json:"method" binding:"required" example:"1"`
	AuthorizationType AuthorizationType `json:"authorization_type" binding:"required" example:"token"`
	Header            struct{}          `json:"header,omitempty" binding:"required"`
	BodyType          BodyType          `json:"body_type" binding:"required" example:"json"`
	Body              struct{}          `json:"body" binding:"required"`
	CommandID         int               `json:"command_id" binding:"required" example:"1"`
}

type SwaggerCreate struct {
	Url               string            `json:"url" binding:"required" example:"http://localhost:9800/api"`
	Method            HttpMethod        `json:"method" binding:"required" example:"POST"`
	AuthorizationType AuthorizationType `json:"authorization_type" example:"token"`
	Header            struct{}          `json:"header"`
	BodyType          BodyType          `json:"body_type" example:"json"`
	Body              struct{}          `json:"body"`
}

type SwaggerUpdate struct {
	Url               string            `json:"url" example:"http://localhost:9800/api"`
	Method            HttpMethod        `json:"method" example:"POST"`
	AuthorizationType AuthorizationType `json:"authorization_type" example:"token"`
	Header            struct{}          `json:"header"`
	BodyType          BodyType          `json:"body_type" example:"json"`
	Body              struct{}          `json:"body"`
}
