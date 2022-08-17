package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"new_command/app"
	"new_command/app/database"
	"new_command/model/https_command"
	"new_command/util"
	"os"
	"testing"
)

func setUpHandler() (modelConfig app.ModelConfig) {
	DB, _ := database.NewDB("mySQL", "DB_test.log", "db_test")
	if err := DB.AutoMigrate(&Command{}, &https_command.HttpsCommand{}); err != nil {
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

func TestCommand(t *testing.T) {
	model := setUpHandler()
	defer func() {
		closeErr := database.CloseDB(model.DB)
		if closeErr != nil {
			log.Println("Error occurred while closing the DB :", closeErr)
		}
	}()
	t.Run("get by id(no data)", func(t *testing.T) {
		w := util.PerformRequest(model.Router, "GET", "/command/api/1", nil)
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
		w := util.PerformRequest(model.Router, "GET", "/command/api", nil)
		require.Equal(t, 200, w.Code)
		var response []map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		require.Nil(t, err)
		require.Equal(t, response, []map[string]string{})
	})
	t.Run("create error data", func(t *testing.T) {
		t.Run("error3", func(t *testing.T) {
			body := []byte(`{
    "name": "forTest",
    "description": "asdfjkl",
    "protocol":"http",
    "https_command": {
        "url":"asdlfkjaslkdfj",
        "Method":"POST",
        "header":{
            "aaa":"aaa",
            "bbb":2
        },
        "body":{
            "as":"ff",
            "bb":{
                "ff":2
            }
        }
    }
}`)
			w := util.PerformRequest(model.Router, "POST", "/command/api", bytes.NewBuffer(body))
			require.Equal(t, 406, w.Code)
			var response map[string]string
			fmt.Println(w.Body.String())
			err := json.Unmarshal([]byte(w.Body.String()), &response)
			value, exists := response["What"]
			require.Nil(t, err)
			require.True(t, exists)
			require.Equal(t, value, "http command need body type")
		})
		t.Run("error4", func(t *testing.T) {
			body := []byte(`{
    "name": "forTest",
    "description": "asdfjkl",
    "protocol":"http",
    "https_command": {
        "url":"asdlfkjaslkdfj",
        "Method":"PATCH",
        "header":{
            "aaa":"aaa",
            "bbb":2
        },
        "body_type":"json"
    }
}`)
			w := util.PerformRequest(model.Router, "POST", "/command/api", bytes.NewBuffer(body))
			require.Equal(t, 406, w.Code)
			var response map[string]string
			err := json.Unmarshal([]byte(w.Body.String()), &response)
			value, exists := response["What"]
			require.Nil(t, err)
			require.True(t, exists)
			require.Equal(t, value, "http command need body")
		})
		t.Run("error5", func(t *testing.T) {
			body := []byte(`{
    "name": "forTest",
    "description": "asdfjkl",
    "protocol":"http"
}`)
			w := util.PerformRequest(model.Router, "POST", "/command/api", bytes.NewBuffer(body))
			require.Equal(t, 406, w.Code)
			var response map[string]string
			err := json.Unmarshal([]byte(w.Body.String()), &response)
			value, exists := response["What"]
			require.Nil(t, err)
			require.True(t, exists)
			require.Equal(t, value, "lose https command!")
		})
	})
}

func TestCommandCreate(t *testing.T) {
	model := setUpHandler()
	defer func() {
		closeErr := database.CloseDB(model.DB)
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
        "Method":"POST",
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
		w := util.PerformRequest(model.Router, "POST", "/command/api", bytes.NewBuffer(body))
		require.Equal(t, http.StatusCreated, w.Code)
		var response map[string]string
		//err := json.Unmarshal([]byte(w.Body.String()), &response)
		err := json.Unmarshal(w.Body.Bytes(), &response)
		value, exists := response["message"]
		require.Nil(t, err)
		require.True(t, exists)
		require.Equal(t, value, "created success")
	})
	t.Run("get data id 1", func(t *testing.T) {
		w := util.PerformRequest(model.Router, "GET", "/command/api/1", nil)
		require.Equal(t, http.StatusOK, w.Code)
		var command Command
		err := json.Unmarshal([]byte(w.Body.String()), &command)
		value := command.ID
		require.Nil(t, err)
		require.NotNil(t, command.HttpsCommand)
		require.Equal(t, value, 1)
	})
	t.Run("get data", func(t *testing.T) {
		w := util.PerformRequest(model.Router, "GET", "/command/api", nil)
		require.Equal(t, http.StatusOK, w.Code)
		var command []Command
		err := json.Unmarshal([]byte(w.Body.String()), &command)
		value := command[0].ID
		require.Nil(t, err)
		require.NotNil(t, command[0].HttpsCommand)
		require.Equal(t, value, 1)
	})
	t.Run("update data", func(t *testing.T) {
		body := []byte(`{
    "name": "test_patch",
    "description": "bbbbbbbbbb",
    "protocol":"http",
    "https_command": {
        "Method":"PATCH",
        "header":{
            "ee":"er",
            "qw":87
        },
        "body_type":"json",
        "body":{
            "ee":"v",
            "we":{
                "be":5
            }
        }
    }
}`)
		w := util.PerformRequest(model.Router, "PATCH", "/command/api/1", bytes.NewBuffer(body))
		var command Command
		err := json.Unmarshal([]byte(w.Body.String()), &command)
		log.Println(command.HttpsCommand)
		require.Nil(t, err)
		require.Equal(t, `{"ee":"v","we":{"be":5}}`, string(command.HttpsCommand.Body))
		require.Equal(t, "test_patch", command.Name)
	})
	t.Run("get data", func(t *testing.T) {
		w := util.PerformRequest(model.Router, "GET", "/command/api/1", nil)
		require.Equal(t, http.StatusOK, w.Code)
		var command Command
		err := json.Unmarshal([]byte(w.Body.String()), &command)
		value := command.ID
		require.Nil(t, err)
		require.NotNil(t, command.HttpsCommand)
		require.Equal(t, value, 1)
		require.Equal(t, command.Description, "bbbbbbbbbb")
	})
	t.Run("delete data", func(t *testing.T) {
		w := util.PerformRequest(model.Router, "DELETE", "/command/api/1", nil)
		require.Equal(t, http.StatusOK, w.Code)
		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		require.Nil(t, err)
		require.Equal(t, "id: 1 has been deleted successfully", response["message"])
	})
}
