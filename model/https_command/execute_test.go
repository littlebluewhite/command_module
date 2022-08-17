package https_command

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExecute(t *testing.T) {
	t.Run("test 1", func(t *testing.T) {
		hc := HttpsCommand{
			CommandID: 1,
			Url:       `fdsa`,
			Method:    Get,
			Header:    []byte(`{"accept":"application/json"}`),
		}
		result := hc.Execute()
		require.Equal(t, ``, result)
	})
	t.Run("test 2", func(t *testing.T) {
		hc := HttpsCommand{
			CommandID: 2,
			Url:       `http://localhost:9800/dispatch_level/api/1`,
			Method:    Get,
			Header:    []byte(`{"accept":"application/json"}`),
		}
		result := hc.Execute()
		require.Equal(t, `{"name":"low","description":"level low","color_code":["#ffffff","#000000"],"id":1}`, result)
	})
	t.Run("test 3", func(t *testing.T) {
		hc := HttpsCommand{
			CommandID: 3,
			Url:       `http://localhost:9800/dispatch_level/api/1`,
			Method:    Patch,
			Header:    []byte(`{"accept":"application/json"}`),
			BodyType:  Json,
			Body:      []byte(`{"description":"level test low"}`),
		}
		result := hc.Execute()
		require.Equal(t, `{"name":"low","description":"level test low","color_code":["#ffffff","#000000"],"id":1}`, result)
	})
	t.Run("test 4", func(t *testing.T) {
		hc := HttpsCommand{
			CommandID: 4,
			Url:       `http://localhost:9800/dispatch_level/api/5`,
			Method:    Patch,
			Header:    []byte(`{"accept":"application/json"}`),
			BodyType:  Json,
			Body:      []byte(`{"description":"level test low"}`),
		}
		result := hc.Execute()
		require.Equal(t, `{"message":"id:5 is not exist"}`, result)
	})
	t.Run("test 5", func(t *testing.T) {
		hc := HttpsCommand{
			CommandID: 5,
			Url:       `http://localhost:9800/dispatch_level/api/1`,
			Method:    Patch,
			Header:    []byte(`{"accept":"application/json"}`),
			BodyType:  Json,
			Body: []byte(`{
  "description": "level low"
}`),
		}
		result := hc.Execute()
		require.Equal(t, `{"name":"low","description":"level low","color_code":["#ffffff","#000000"],"id":1}`, result)
	})
}
