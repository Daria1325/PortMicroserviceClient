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
	"os"
	"time"
)

var (
	JsonPath string
	Connect  Connection
)

type Connection struct {
	Conn   *grpc.ClientConn
	client api.PortClient
}

func (c *Connection) InitConn(host string, port string) {
	s := host + port
	conn, err := grpc.Dial(s, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := api.NewPortClient(conn)
	c.Conn = conn
	c.client = client
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

func GetPorts(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Print("Looking for ports")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := Connect.client.GetPorts(ctx, &api.GetPortsRequest{Name: "1"})
	if err != nil {
		log.Fatal(err)
	}
	respString := resp.GetList()
	jsonRes := StringToJson(respString)
	if jsonRes != nil {
		err := json.NewEncoder(w).Encode(jsonRes)
		if err != nil {
			log.Fatal(err)
			return
		}
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
		return
	}
	w.WriteHeader(200)

	respString := resp.GetItem()
	jsonRes := StringToJson(respString)
	if jsonRes != nil {
		err := json.NewEncoder(w).Encode(jsonRes)
		if err != nil {
			return
		}
	}
}

func UpsertPorts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var jsonData []byte
	var err error
	if r.ContentLength != 0 {
		//jsonData, err = json.Marshal(r.Body)
		jsonData, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		jsonData, err = ReadJson(JsonPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Print("Updating ports")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := Connect.client.UpsertPorts(ctx, &api.UpsertPortsRequest{Name: string(jsonData)})
	if err != nil {
		log.Fatal(err)
	}

	respString := resp.GetList()
	jsonRes := StringToJson(respString)
	if jsonRes != nil {
		err := json.NewEncoder(w).Encode(jsonRes)
		if err != nil {
			return
		}
	}
}

func ReadJson(path string) ([]byte, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	defer jsonFile.Close()
	return byteValue, nil
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
