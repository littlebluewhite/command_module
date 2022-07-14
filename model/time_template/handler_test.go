package time_template

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"new_command/app"
	"new_command/app/database"
	"os"
	"testing"
)

func setUpHandler() (modelConfig app.ModelConfig) {
	DB, _ := database.NewDB("mySQL", "DB_test.log", "db_test")
	if err := DB.AutoMigrate(&TimeTemplate{}, &WeeklyRepeat{},
		&MonthlyRepeat{}); err != nil {
		log.Println("Error occurred while Migrate test DB :", err)
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

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
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
		w := performRequest(model.Router, "GET", "/time_template/api/1", nil)
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
		w := performRequest(model.Router, "GET", "/time_template/api", nil)
		require.Equal(t, 200, w.Code)
		var response []map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		require.Nil(t, err)
		require.Equal(t, response, []map[string]string{})
	})
	t.Run("create data", func(t *testing.T) {
		t.Run("error", func(t *testing.T) {
			body := []byte(`{
    "name": "test",
    "repeat_type": "weekly",
    "start_date": "2022-07-04T01:01:01+08:00",
    "start_time": "12:12:12",
    "end_time": "22:13:13",
    "weekly_repeat":{
        "weekly_condition":[1, 2, 3, 7]
    }
}`)
			w := performRequest(model.Router, "POST", "/time_template/api", bytes.NewBuffer(body))
			require.Equal(t, 406, w.Code)
			var response map[string]string
			err := json.Unmarshal([]byte(w.Body.String()), &response)
			value, exists := response["What"]
			require.Nil(t, err)
			require.True(t, exists)
			require.Equal(t, value, "weekly conditions number are not correct")
		})
		t.Run("error2", func(t *testing.T) {
			body := []byte(`{
    "name": "test",
    "repeat_type": "weekly",
    "start_date": "2022-07-04T01:01:01+08:00",
    "start_time": "12:12:12",
    "end_time": "02:13:13",
    "weekly_repeat":{
        "weekly_condition":[1, 2, 3]
    }
}`)
			w := performRequest(model.Router, "POST", "/time_template/api", bytes.NewBuffer(body))
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
    "name": "test",
    "repeat_type": "weekly",
    "start_date": "2022-07-04T01:01:01+08:00",
    "end_date": "2022-07-02T01:01:01+08:00",
    "start_time": "12:12:12",
    "end_time": "22:13:13",
    "weekly_repeat":{
        "weekly_condition":[1, 2, 3]
    }
}`)
			w := performRequest(model.Router, "POST", "/time_template/api", bytes.NewBuffer(body))
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
    "name": "test",
    "repeat_type": "weekly",
    "start_date": "2022-07-04T01:01:01+08:00",
    "end_date": "2022-07-06T01:01:01+08:00",
    "start_time": "12:12:12",
    "end_time": "22:13:13",
    "weekly_repeat":{
        "weekly_condition":{}
    }
}`)
			w := performRequest(model.Router, "POST", "/time_template/api", bytes.NewBuffer(body))
			require.Equal(t, 406, w.Code)
			var response map[string]string
			err := json.Unmarshal([]byte(w.Body.String()), &response)
			value, exists := response["What"]
			require.Nil(t, err)
			require.True(t, exists)
			require.Equal(t, value, "weekly conditions are not correct")
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
    "name": "test",
    "repeat_type": "weekly",
    "start_date": "2022-07-04T01:01:01+08:00",
    "end_date": "2022-07-06T01:01:01+08:00",
    "start_time": "12:12:12",
    "end_time": "22:13:13",
    "weekly_repeat":{
        "weekly_condition":[0, 1, 2]
    }
}`)
		w := performRequest(model.Router, "POST", "/time_template/api", bytes.NewBuffer(body))
		require.Equal(t, http.StatusCreated, w.Code)
		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		value, exists := response["message"]
		require.Nil(t, err)
		require.True(t, exists)
		require.Equal(t, value, "created success")
	})
	t.Run("get data id 1", func(t *testing.T) {
		w := performRequest(model.Router, "GET", "/time_template/api/1", nil)
		require.Equal(t, http.StatusOK, w.Code)
		var timeTemplate TimeTemplate
		err := json.Unmarshal([]byte(w.Body.String()), &timeTemplate)
		value := timeTemplate.ID
		require.Nil(t, err)
		require.Equal(t, value, 1)
	})
	t.Run("get data", func(t *testing.T) {
		w := performRequest(model.Router, "GET", "/time_template/api", nil)
		require.Equal(t, http.StatusOK, w.Code)
		var timeTemplate []TimeTemplate
		err := json.Unmarshal([]byte(w.Body.String()), &timeTemplate)
		value := timeTemplate[0].ID
		require.Nil(t, err)
		require.Equal(t, value, 1)
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
    "name": "test2",
    "repeat_type": "weekly",
    "start_date": "2022-07-04T01:01:01+08:00",
    "end_date": "2022-07-06T01:01:01+08:00",
    "start_time": "12:12:12",
    "end_time": "22:13:13",
    "weekly_repeat":{
        "weekly_condition":[0, 1, 2]
    }
}`)
		w := performRequest(model.Router, "POST", "/time_template/api", bytes.NewBuffer(body))
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
    "name": "test3",
    "repeat_type": "monthly",
    "start_date": "2022-07-04T01:01:01+08:00",
    "end_date": "2022-07-06T01:01:01+08:00",
    "start_time": "12:12:12",
    "end_time": "22:13:13",
    "monthly_repeat": {
        "first_weekly_condition": [0, 1, 2],
		"second_weekly_condition": [],
		"third_weekly_condition": [],
		"fourth_weekly_condition": [],
		"monthly_condition": []
    }
}`)
		w := performRequest(model.Router, "PATCH", "/time_template/api/1", bytes.NewBuffer(body))
		var timeTemplate TimeTemplate
		err := json.Unmarshal([]byte(w.Body.String()), &timeTemplate)
		require.Nil(t, err)
		require.Equal(t, nil, timeTemplate)
	})
}
