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

//START THE SERVER

func TestMain(m *testing.M) {
	config, _ := cnfg.NewConfigPath("configs/dataConfig.toml")
	client.Connect.InitConn("localhost", config.BindAddrServer)

	exitcode := m.Run()
	os.Exit(exitcode)
}

func TestInitConn(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("panic")
		}
	}()
	client.Connect.InitConn("localhost", ":9080")
}

func TestGetPorts(t *testing.T) {
	t.Parallel()
	//передать файл
	r, _ := http.NewRequest("GET", "/ports", nil)
	w := httptest.NewRecorder()

	client.GetPorts(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}
func TestGetPortValid(t *testing.T) {
	t.Parallel()

	r, _ := http.NewRequest("GET", "/port/id", nil)
	w := httptest.NewRecorder()
	vars := map[string]string{
		"id": "2",
	}
	r = mux.SetURLVars(r, vars)
	client.GetPort(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

}
func TestGetPortInvalid(t *testing.T) {
	t.Parallel()

	r, _ := http.NewRequest("GET", "/port/id", nil)
	w := httptest.NewRecorder()
	vars := map[string]string{
		"id": "1000",
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
