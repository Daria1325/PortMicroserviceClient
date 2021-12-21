package client

import (
	"context"
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

func (c *Connection) InitConn(port string) {
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := api.NewPortClient(conn)
	c.Conn = conn
	c.client = client
}

func GetPorts(w http.ResponseWriter, r *http.Request) {
	log.Print("Looking for ports")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := Connect.client.GetPorts(ctx, &api.GetPortsRequest{Name: "1"})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp.GetList())

}
func GetPort(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	log.Print("Looking for port")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := Connect.client.GetPort(ctx, &api.GetPortRequest{Id: params["id"]})
	if err != nil {
		w.WriteHeader(404)
	}
	log.Println(resp.GetItem())
}

func UpsertPorts(w http.ResponseWriter, r *http.Request) {
	jsonData, err := ReadJson(JsonPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Updating ports")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := Connect.client.UpsertPorts(ctx, &api.UpsertPortsRequest{Name: string(jsonData)})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp.GetList())
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
	log.Print(config.BindAddrServer, config.BindAddrOuter)
	err := http.ListenAndServe(config.BindAddrOuter, r)

	if err != nil {
		return err
	}
	return nil
}
