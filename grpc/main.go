package main

import (
	"context"
	"errors"
	"fmt"
	hostname "github.com/trainyao/sofastack_test/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

const HOST1_SERVER_PORT = "20001"
const HOST2_SERVER_PORT = "20002"
const HOST3_SERVER_PORT = "20003"
const HOST4_SERVER_PORT = "20004"

const HOST2_SERVER_HOST = "svc_host2"
const HOST3_SERVER_HOST = "svc_host3"
const HOST4_SERVER_HOST = "svc_host4"

type portMap struct {
	ServerPort   string
	UpstreamHost string
	UpstreamPort string
}

var client hostname.HostnameClient

var serverPortMap = map[string]portMap{
	"HOST1": {
		ServerPort:   HOST1_SERVER_PORT,
		UpstreamHost: HOST2_SERVER_HOST,
		UpstreamPort: HOST2_SERVER_PORT,
	},
	"HOST2": {
		ServerPort:   HOST2_SERVER_PORT,
		UpstreamHost: HOST3_SERVER_HOST,
		UpstreamPort: HOST3_SERVER_PORT,
	},
	"HOST3": {
		ServerPort:   HOST3_SERVER_PORT,
		UpstreamHost: HOST4_SERVER_HOST,
		UpstreamPort: HOST4_SERVER_PORT,
	},
	"HOST4": {
		ServerPort:   HOST4_SERVER_PORT,
		UpstreamHost: "",
		UpstreamPort: "",
	},
}

func main() {
	var host string
	var found bool
	if host, found = os.LookupEnv("GRPCTEST_HOSTNAME"); found {
		log.Printf("Host: %s", host)
	}

	config, err := getPortConfig(host)
	if err != nil {
		panic(err)
	}

	l, err := net.Listen("tcp", "0.0.0.0:"+config.ServerPort)
	if err != nil {
		panic(err)
	}

	hostnameServer := &server{
		Host:        host,
		ClientExist: false,
	}

	hostnameServer.initClient()

	s := grpc.NewServer()
	hostname.RegisterHostnameServer(s, hostnameServer)
	reflection.Register(s)

	log.Println("serving")
	if err = s.Serve(l); err != nil {
		panic(err)
	}

	log.Println("exit")
}

func getPortConfig(host string) (result *portMap, err error) {
	var ok bool
	var ret portMap

	if ret, ok = serverPortMap[host]; !ok {
		return nil, errors.New("not found")
	}

	return &ret, nil
}

type server struct {
	Host        string
	Client      hostname.HostnameClient
	ClientExist bool
}

func (s *server) GetHostname(ctx context.Context, in *hostname.HostnameRequest) (out *hostname.HostnameResponse, err error) {
	log.Println("called")

	var resultString string
	if s.ClientExist {
		r, err := s.Client.GetHostname(context.Background(), &hostname.HostnameRequest{Test: ""})
		if err != nil {
			resultString = "query failed, err: " + err.Error()
		} else {
			resultString = "query success, res: \n" + r.Hostname
		}
	} else {
		resultString = "\nend"
	}

	out = &hostname.HostnameResponse{
		Hostname: fmt.Sprintf("host is %s, result from next server: %s", s.Host, resultString),
	}

	return out, nil
}

func (s *server) initClient() {
	config, err := getPortConfig(s.Host)
	if err != nil {
		panic(err)
	}

	if config.UpstreamPort == "" {
		return
	}

	conn, err := grpc.Dial("localhost:"+config.UpstreamPort, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	s.Client = hostname.NewHostnameClient(conn)
	s.ClientExist = true
}
