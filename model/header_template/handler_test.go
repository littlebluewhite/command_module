package header_template

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"new_command/app"
	"new_command/app/database"
	"new_command/util"
	"os"
	"testing"
	"time"
)

func setUpHandler() (modelConfig app.ModelConfig) {
	DB, _ := database.NewDB("mySQL", "DB_test.log", "db_test")
	if err := DB.AutoMigrate(&HeaderTemplate{}); err != nil {
		log.Println("Error occurred while Migrate forTest DB :", err)
	}
	ginFile, err := os.OpenFile("./log/gin_test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("can not open log file: " + err.Error())
	}
	gin.DefaultWriter = io.MultiWriter(ginFile, os.Stdout)
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	c := cache.New(5*time.Minute, 10*time.Minute)
	modelConfig = app.ModelConfig{
		DB:     DB,
		Router: r,
		Cache:  c,
	}
	Inject(modelConfig)
	return
}

func TestHeader(t *testing.T) {
	model2 := setUpHandler()
	defer func() {
		closeErr := database.CloseDB(model2.DB)
		if closeErr != nil {
			log.Println("Error occurred while closing the DB :", closeErr)
		}
	}()
	t.Run("get by id(no data)", func(t *testing.T) {
		w := util.PerformRequest(model2.Router, "GET", "/header_template/api/1", nil)
		require.Equal(t, http.StatusBadRequest, w.Code)
		body := gin.H{"message": "id is not correct"}
		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		value, exists := response["message"]
		require.Nil(t, err)
		require.True(t, exists)
		require.Equal(t, body["message"], value)
	})
	t.Run("get all(no data)", func(t *testing.T) {
		w := util.PerformRequest(model2.Router, "GET", "/header_template/api", nil)
		require.Equal(t, 200, w.Code)
		var response []map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		require.Nil(t, err)
		require.Equal(t, response, []map[string]string(nil))
	})
}

func TestHeaderTemplateCreate(t *testing.T) {
	model2 := setUpHandler()
	defer func() {
		closeErr := database.CloseDB(model2.DB)
		if closeErr != nil {
			log.Println("Error occurred while closing the DB :", closeErr)
		}
	}()
	t.Run("create header template", func(t *testing.T) {
		body := []byte(`{
    "name": "test_header_template",
    "data": {
		"a": "aa1",
		"b": "bb1"
	}
}`)
		w := util.PerformRequest(model2.Router, "POST", "/header_template/api", bytes.NewBuffer(body))
		require.Equal(t, http.StatusCreated, w.Code)
		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		value, exists := response["message"]
		require.Nil(t, err)
		require.True(t, exists)
		require.Equal(t, value, "created success")
	})
	t.Run("get data id 1", func(t *testing.T) {
		w := util.PerformRequest(model2.Router, "GET", "/header_template/api/1", nil)
		require.Equal(t, http.StatusOK, w.Code)
		var headerTemplate HeaderTemplate
		err := json.Unmarshal([]byte(w.Body.String()), &headerTemplate)
		value := headerTemplate.ID
		require.Nil(t, err)
		require.Equal(t, value, 1)
	})
	t.Run("get data", func(t *testing.T) {
		w := util.PerformRequest(model2.Router, "GET", "/header_template/api", nil)
		require.Equal(t, http.StatusOK, w.Code)
		var headerTemplates []HeaderTemplate
		err := json.Unmarshal([]byte(w.Body.String()), &headerTemplates)
		value := headerTemplates[0].ID
		require.Nil(t, err)
		require.Equal(t, value, 1)
	})
	t.Run("update data", func(t *testing.T) {
		body := []byte(`{
    "name": "test_header_template_update",
    "data": {
		"a": "aa2",
		"b": "bb2"
	}
}`)
		w := util.PerformRequest(model2.Router, "PATCH", "/header_template/api/1", bytes.NewBuffer(body))
		var headerTemplate HeaderTemplate
		err := json.Unmarshal([]byte(w.Body.String()), &headerTemplate)
		require.Nil(t, err)
		require.Equal(t, `{"a":"aa2","b":"bb2"}`, headerTemplate.Data.String())
		require.Equal(t, "test_header_template_update", headerTemplate.Name)
	})
	t.Run("get data", func(t *testing.T) {
		w := util.PerformRequest(model2.Router, "GET", "/header_template/api/1", nil)
		require.Equal(t, http.StatusOK, w.Code)
		var headerTemplate HeaderTemplate
		err := json.Unmarshal([]byte(w.Body.String()), &headerTemplate)
		value := headerTemplate.ID
		require.Nil(t, err)
		require.Equal(t, value, 1)
		require.Equal(t, headerTemplate.Name, "test_header_template_update")
	})
	t.Run("delete data", func(t *testing.T) {
		w := util.PerformRequest(model2.Router, "DELETE", "/header_template/api/1", nil)
		require.Equal(t, http.StatusOK, w.Code)
		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		require.Nil(t, err)
		require.Equal(t, "id: 1 has been deleted successfully", response["message"])
	})
}
