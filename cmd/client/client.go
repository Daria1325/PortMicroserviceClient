package client

import (
	"context"
	"encoding/json"
	api "github.com/daria/PortMicroserviceClient/api/proto"
	cnfg "github.com/daria/PortMicroserviceClient/data/config"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	Connect Connection
)

type Connection struct {
	Conn   *grpc.ClientConn
	client api.PortClient
}

// InitConn :TODO error handle not panic
func (c *Connection) InitConn(host string, port string) error {
	s := host + port
	conn, err := grpc.Dial(s, grpc.WithInsecure())
	if err != nil {
		return err
	}
	client := api.NewPortClient(conn)
	c.Conn = conn
	c.client = client
	return nil
}
func StringToJson(s string) []map[string]interface{} {
	var jsonMap []map[string]interface{}
	err := json.Unmarshal([]byte(s), &jsonMap)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return jsonMap
}

// WriterResp :TODO test the function
func WriterResp(w http.ResponseWriter, respString string) error {
	jsonRes := StringToJson(respString)
	if jsonRes != nil {
		err := json.NewEncoder(w).Encode(jsonRes)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetPorts(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Print("Looking for ports")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := Connect.client.GetPorts(ctx, &api.GetPortsRequest{Name: ""})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	respString := resp.GetList()
	err = WriterResp(w, respString)
	if err != nil {
		http.Error(w, err.Error(), 520)
		return
	}
}
func GetPort(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	log.Print("Looking for port")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := Connect.client.GetPort(ctx, &api.GetPortRequest{Id: params["id"]})
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)

	respString := resp.GetItem()

	err = WriterResp(w, respString)
	if err != nil {
		http.Error(w, err.Error(), 520)
		return
	}
}

func UpsertPorts(w http.ResponseWriter, r *http.Request) {
	var jsonData []byte
	var err error
	if r.ContentLength != 0 {
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Could not get the data file", 400)
			return
		}
		jsonData, err = ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "Could not read the file", 400)
			return
		}
	} else {
		http.Error(w, "Data isn`t provided", 400)
		return
	}

	log.Print("Updating ports")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := Connect.client.UpsertPorts(ctx, &api.UpsertPortsRequest{Name: string(jsonData)})
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	respString := resp.GetList()
	_, err = w.Write([]byte(respString))
	if err != nil {
		http.Error(w, err.Error(), 520)
		return
	}
}

func Start(config *cnfg.Config) error {
	r := mux.NewRouter()

	r.HandleFunc("/ports", GetPorts).Methods("GET")
	r.HandleFunc("/ports/{id}", GetPort).Methods("GET")
	r.HandleFunc("/ports", UpsertPorts).Methods("POST")

	err := http.ListenAndServe(config.BindAddrOuter, r)
	if err != nil {
		return err
	}
	return nil
}
