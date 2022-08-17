package schedule

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
	"new_command/model/command"
	"new_command/model/https_command"
	"new_command/model/time_data"
	"new_command/util"
	"os"
	"testing"
	"time"
)

func setUpHandler() (modelConfig app.ModelConfig) {
	DB, _ := database.NewDB("mySQL", "DB_test.log", "db_test")
	if err := DB.AutoMigrate(&Schedule{}, &https_command.HttpsCommand{},
		&command.Command{}, &time_data.TimeData{}); err != nil {
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
	command.Inject(modelConfig)
	return
}

func TestScheduler(t *testing.T) {
	model2 := setUpHandler()
	defer func() {
		closeErr := database.CloseDB(model2.DB)
		if closeErr != nil {
			log.Println("Error occurred while closing the DB :", closeErr)
		}
	}()
	t.Run("get by id(no data)", func(t *testing.T) {
		w := util.PerformRequest(model2.Router, "GET", "/schedule/api/1", nil)
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
		w := util.PerformRequest(model2.Router, "GET", "/schedule/api", nil)
		require.Equal(t, 200, w.Code)
		var response []map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		require.Nil(t, err)
		require.Equal(t, response, []map[string]string(nil))
	})
}

func TestScheduleCreate(t *testing.T) {
	model2 := setUpHandler()
	defer func() {
		closeErr := database.CloseDB(model2.DB)
		if closeErr != nil {
			log.Println("Error occurred while closing the DB :", closeErr)
		}
	}()
	t.Run("create http command", func(t *testing.T) {
		body := []byte(`{
    "name": "forTest",
    "description": "asdfjkl",
    "protocol":"http",
    "updated_at": "2022-07-25T14:27:15.199+08:00",
    "created_at": "2022-07-25T06:27:08+08:00",
    "https_command": {
        "url":"asdlfkjaslkdfj",
        "Method":"post",
        "header":{
            "aaa":"aaa",
            "bbb":2
        },
        "body_type":"json",
        "body":{
            "as":"ff",
            "bb":{
                "ff":2
            }
        }
    }
}`)
		w := util.PerformRequest(model2.Router, "POST", "/command/api", bytes.NewBuffer(body))
		require.Equal(t, http.StatusCreated, w.Code)
		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		value, exists := response["message"]
		require.Nil(t, err)
		require.True(t, exists)
		require.Equal(t, value, "created success")
	})
	t.Run("create schedule", func(t *testing.T) {
		body := []byte(`{
    "name": "test_schedule",
    "description": "asdfqwersasf",
    "enabled": true,
    "command_id":1,
    "updated_at": "2022-07-25T14:27:15.199+08:00",
    "created_at": "2022-07-25T06:27:08+08:00",
    "time_data":{
        "repeat_type": "weekly",
        "start_date": "2022-07-03T01:01:01+08:00",
        "end_date": "2022-07-13T01:01:01+08:00",
        "start_time": "00:00:01",
        "end_time": "02:13:13",
        "interval_seconds": 300,
        "condition_type": "weekly_day",
        "condition": [3, 4]
    }
}`)
		w := util.PerformRequest(model2.Router, "POST", "/schedule/api", bytes.NewBuffer(body))
		require.Equal(t, http.StatusCreated, w.Code)
		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		value, exists := response["message"]
		require.Nil(t, err)
		require.True(t, exists)
		require.Equal(t, value, "created success")
	})
	t.Run("get data id 1", func(t *testing.T) {
		w := util.PerformRequest(model2.Router, "GET", "/schedule/api/1", nil)
		require.Equal(t, http.StatusOK, w.Code)
		var schedule Schedule
		err := json.Unmarshal([]byte(w.Body.String()), &schedule)
		value := schedule.ID
		require.Nil(t, err)
		require.NotNil(t, schedule.TimeData)
		require.NotNil(t, schedule.Command)
		require.Equal(t, value, 1)
	})
	t.Run("get data", func(t *testing.T) {
		w := util.PerformRequest(model2.Router, "GET", "/schedule/api", nil)
		require.Equal(t, http.StatusOK, w.Code)
		var schedule []Schedule
		err := json.Unmarshal([]byte(w.Body.String()), &schedule)
		value := schedule[0].ID
		require.Nil(t, err)
		require.NotNil(t, schedule[0].TimeData)
		require.NotNil(t, schedule[0].Command)
		require.Equal(t, value, 1)
	})
	t.Run("update data", func(t *testing.T) {
		body := []byte(`{
    "name": "test_patch_schedule",
    "description": "eeeeeeeeeess",
    "enabled": false,
    "time_data":{
        "repeat_type": "monthly",
        "start_date": "2022-07-03T01:01:01+08:00",
        "end_date": "2022-07-13T01:01:01+08:00",
        "start_time": "00:00:01",
        "end_time": "02:13:13",
        "interval_seconds": 500,
        "condition_type": "monthly_day",
        "condition": [2, 4, 15, 16, 20, 22, 30]
    }
}`)
		w := util.PerformRequest(model2.Router, "PATCH", "/schedule/api/1", bytes.NewBuffer(body))
		var schedule Schedule
		err := json.Unmarshal([]byte(w.Body.String()), &schedule)
		log.Println(schedule.TimeData)
		require.Nil(t, err)
		require.Equal(t, `[2,4,15,16,20,22,30]`, schedule.TimeData.Condition.String())
		require.Equal(t, "test_patch_schedule", schedule.Name)
	})
	t.Run("get data", func(t *testing.T) {
		w := util.PerformRequest(model2.Router, "GET", "/schedule/api/1", nil)
		require.Equal(t, http.StatusOK, w.Code)
		var schedule Schedule
		err := json.Unmarshal([]byte(w.Body.String()), &schedule)
		value := schedule.ID
		require.Nil(t, err)
		require.NotNil(t, schedule.Command)
		require.NotNil(t, schedule.TimeData)
		require.Equal(t, value, 1)
		require.Equal(t, schedule.Description, "eeeeeeeeeess")
	})
	t.Run("delete data", func(t *testing.T) {
		w := util.PerformRequest(model2.Router, "DELETE", "/schedule/api/1", nil)
		require.Equal(t, http.StatusOK, w.Code)
		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		require.Nil(t, err)
		require.Equal(t, "id: 1 has been deleted successfully", response["message"])
	})
}
