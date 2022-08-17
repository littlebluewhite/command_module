package time_template

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"new_command/app"
	"new_command/app/database"
	"new_command/model/time_data"
	"new_command/util"
	"os"
	"testing"
)

func setUpHandler() (modelConfig app.ModelConfig) {
	DB, _ := database.NewDB("mySQL", "DB_test.log", "db_test")
	if err := DB.AutoMigrate(&TimeTemplate{}, &time_data.TimeData{}); err != nil {
		log.Println("Error occurred while Migrate forTest DB :", err)
	}
	ginFile, err := os.OpenFile("./log/gin_test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("can not open log file: " + err.Error())
	}
	gin.DefaultWriter = io.MultiWriter(ginFile, os.Stdout)
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	modelConfig = app.ModelConfig{
		DB:     DB,
		Router: r,
	}
	Inject(modelConfig)
	return
}

func TestGetTimeTemplate(t *testing.T) {
	model := setUpHandler()
	defer func() {
		closeErr := database.CloseDB(model.DB)
		if closeErr != nil {
			log.Println("Error occurred while closing the DB :", closeErr)
		}
	}()
	t.Run("get by id(no data)", func(t *testing.T) {
		w := util.PerformRequest(model.Router, "GET", "/time_template/api/1", nil)
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
		w := util.PerformRequest(model.Router, "GET", "/time_template/api", nil)
		require.Equal(t, 200, w.Code)
		var response []map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		require.Nil(t, err)
		require.Equal(t, response, []map[string]string{})
	})
	t.Run("create data", func(t *testing.T) {
		t.Run("error", func(t *testing.T) {
			body := []byte(`{
    "name": "forTest",
    "time_data":{
        "repeat_type": "weekly",
        "start_date": "2022-07-06T01:01:01+08:00",
        "end_date": "2022-07-08T01:01:01+08:00",
        "start_time": "00:00:00",
        "end_time": "22:13:13",
        "interval_seconds": 300,
        "condition_type": "weekly_day",
        "condition": [0, 1, 2, 3, 4, 5, 7]
    }
}`)
			w := util.PerformRequest(model.Router, "POST", "/time_template/api", bytes.NewBuffer(body))
			require.Equal(t, 406, w.Code)
			var response map[string]string
			err := json.Unmarshal([]byte(w.Body.String()), &response)
			value, exists := response["What"]
			require.Nil(t, err)
			require.True(t, exists)
			require.Equal(t, value, "weekly condition number are not correct")
		})
		t.Run("error2", func(t *testing.T) {
			body := []byte(`{
    "name": "forTest",
    "time_data":{
        "repeat_type": "weekly",
        "start_date": "2022-07-06T01:01:01+08:00",
        "end_date": "2022-07-08T01:01:01+08:00",
        "start_time": "22:00:00",
        "end_time": "12:13:13",
        "interval_seconds": 300,
        "condition_type": "weekly_day",
        "condition": [0, 1, 2, 3, 4, 5]
    }
}`)
			w := util.PerformRequest(model.Router, "POST", "/time_template/api", bytes.NewBuffer(body))
			require.Equal(t, 406, w.Code)
			var response map[string]string
			err := json.Unmarshal([]byte(w.Body.String()), &response)
			value, exists := response["What"]
			require.Nil(t, err)
			require.True(t, exists)
			require.Equal(t, value, "start time and end time error")
		})
		t.Run("error3", func(t *testing.T) {
			body := []byte(`{
    "name": "forTest",
    "time_data":{
        "repeat_type": "weekly",
        "start_date": "2022-07-16T01:01:01+08:00",
        "end_date": "2022-07-08T01:01:01+08:00",
        "start_time": "02:00:00",
        "end_time": "12:13:13",
        "interval_seconds": 300,
        "condition_type": "weekly_day",
        "condition": [0, 1, 2, 3, 4, 5]
    }
}`)
			w := util.PerformRequest(model.Router, "POST", "/time_template/api", bytes.NewBuffer(body))
			require.Equal(t, 406, w.Code)
			var response map[string]string
			err := json.Unmarshal([]byte(w.Body.String()), &response)
			value, exists := response["What"]
			require.Nil(t, err)
			require.True(t, exists)
			require.Equal(t, value, "start date and end date error")
		})
		t.Run("error3", func(t *testing.T) {
			body := []byte(`{
    "name": "forTest",
    "time_data":{
        "repeat_type": "weekly",
        "start_date": "2022-07-06T01:01:01+08:00",
        "end_date": "2022-07-08T01:01:01+08:00",
        "start_time": "02:00:00",
        "end_time": "12:13:13",
        "interval_seconds": 300,
        "condition_type": "weekly_day",
        "condition": {}
    }
}`)
			w := util.PerformRequest(model.Router, "POST", "/time_template/api", bytes.NewBuffer(body))
			require.Equal(t, 406, w.Code)
			var response map[string]string
			err := json.Unmarshal([]byte(w.Body.String()), &response)
			value, exists := response["What"]
			require.Nil(t, err)
			require.True(t, exists)
			require.Equal(t, value, "condition format are not correct")
		})
	})
}

func TestGetTimeTemplateCreate(t *testing.T) {
	model := setUpHandler()
	defer func() {
		closeErr := database.CloseDB(model.DB)
		if closeErr != nil {
			log.Println("Error occurred while closing the DB :", closeErr)
		}
	}()
	t.Run("create", func(t *testing.T) {
		body := []byte(`{
    "name": "forTest",
    "time_data":{
        "repeat_type": "weekly",
        "start_date": "2022-07-06T01:01:01+08:00",
        "end_date": "2022-07-08T01:01:01+08:00",
        "start_time": "02:00:00",
        "end_time": "12:13:13",
        "interval_seconds": 300,
        "condition_type": "weekly_day",
        "condition": [1, 2, 3]
    }
}`)
		w := util.PerformRequest(model.Router, "POST", "/time_template/api", bytes.NewBuffer(body))
		require.Equal(t, http.StatusCreated, w.Code)
		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		value, exists := response["message"]
		require.Nil(t, err)
		require.True(t, exists)
		require.Equal(t, value, "created success")
	})
	t.Run("get data id 1", func(t *testing.T) {
		w := util.PerformRequest(model.Router, "GET", "/time_template/api/1", nil)
		require.Equal(t, http.StatusOK, w.Code)
		var timeTemplate TimeTemplate
		err := json.Unmarshal([]byte(w.Body.String()), &timeTemplate)
		value := timeTemplate.ID
		require.Nil(t, err)
		require.Equal(t, value, 1)
	})
	t.Run("get data", func(t *testing.T) {
		w := util.PerformRequest(model.Router, "GET", "/time_template/api", nil)
		require.Equal(t, http.StatusOK, w.Code)
		var timeTemplate []TimeTemplate
		err := json.Unmarshal([]byte(w.Body.String()), &timeTemplate)
		value := timeTemplate[0].ID
		require.Nil(t, err)
		require.Equal(t, value, 1)
	})
	t.Run("delete data", func(t *testing.T) {
		w := util.PerformRequest(model.Router, "DELETE", "/time_template/api/1", nil)
		require.Equal(t, http.StatusOK, w.Code)
		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		require.Nil(t, err)
		require.Equal(t, "id: 1 has been deleted successfully", response["message"])
	})
}

func TestUpdateTimeTemplate(t *testing.T) {
	model := setUpHandler()
	defer func() {
		closeErr := database.CloseDB(model.DB)
		if closeErr != nil {
			log.Println("Error occurred while closing the DB :", closeErr)
		}
	}()
	t.Run("create data", func(t *testing.T) {
		body := []byte(`{
    "name": "forTest",
	"created_at": "2022-07-06T01:01:01+08:00",
    "updated_at": "2022-07-06T01:01:01+08:00",
    "time_data":{
        "repeat_type": "weekly",
        "start_date": "2022-07-06T01:01:01+08:00",
        "end_date": "2022-07-08T01:01:01+08:00",
        "start_time": "02:00:00",
        "end_time": "12:13:13",
        "interval_seconds": 300,
        "condition_type": "weekly_day",
        "condition": [1, 2, 3]
    }
}`)
		w := util.PerformRequest(model.Router, "POST", "/time_template/api", bytes.NewBuffer(body))
		require.Equal(t, http.StatusCreated, w.Code)
		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		value, exists := response["message"]
		require.Nil(t, err)
		require.True(t, exists)
		require.Equal(t, value, "created success")
	})
	t.Run("update data", func(t *testing.T) {
		body := []byte(`{
    "name": "asdf",
    "time_data":{
        "repeat_type": "weekly",
        "start_date": "2022-07-02T01:01:01+08:00",
        "end_date": "2022-07-10T01:01:01+08:00",
        "start_time": "03:13:00",
        "end_time": "14:14:14",
        "interval_seconds": 150,
        "condition_type": "weekly_day",
		"condition": [1, 2]
    }
}`)
		w := util.PerformRequest(model.Router, "PATCH", "/time_template/api/1", bytes.NewBuffer(body))
		var timeTemplate TimeTemplate
		err := json.Unmarshal([]byte(w.Body.String()), &timeTemplate)
		require.Nil(t, err)
		require.Equal(t, "asdf", timeTemplate.Name)
	})
	t.Run("delete data", func(t *testing.T) {
		w := util.PerformRequest(model.Router, "DELETE", "/time_template/api/1", nil)
		require.Equal(t, http.StatusOK, w.Code)
		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		require.Nil(t, err)
		require.Equal(t, "id: 1 has been deleted successfully", response["message"])
	})
}
