package tests

import (
	"github.com/daria/PortMicroserviceClient/cmd/client"
	cnfg "github.com/daria/PortMicroserviceClient/data/config"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	config, _ := cnfg.NewConfigPath("configs/dataConfig.toml")
	client.JsonPath = config.JsonPath
	client.Connect.InitConn(config.BindAddrServer)

	exitcode := m.Run()
	os.Exit(exitcode)
}

func TestInitConn(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("panic")
		}
	}()
	client.Connect.InitConn(":9080")
}

func TestGetPorts(t *testing.T) {
	t.Parallel()

	r, _ := http.NewRequest("GET", "/ports", nil)
	w := httptest.NewRecorder()

	client.GetPorts(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}
func TestGetPort(t *testing.T) {
	t.Parallel()

	r, _ := http.NewRequest("GET", "/port/id", nil)
	w := httptest.NewRecorder()
	vars := map[string]string{
		"id": "2",
	}
	r = mux.SetURLVars(r, vars)
	client.GetPort(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	vars = map[string]string{
		"id": "0",
	}
	r = mux.SetURLVars(r, vars)
	client.GetPort(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
func TestUpsertPorts(t *testing.T) {
	t.Parallel()

	r, _ := http.NewRequest("POST", "/ports", nil)
	w := httptest.NewRecorder()

	client.UpsertPorts(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestNewConfigPath(t *testing.T) {
	input := "some.toml"
	_, err := cnfg.NewConfigPath(input)
	assert.Error(t, err)
}
func TestReadJson(t *testing.T) {
	input1 := "../data/ports.json"
	_, err := client.ReadJson(input1)
	assert.NoError(t, err)

	input2 := "some.json"
	_, err = client.ReadJson(input2)
	assert.Error(t, err)
}
