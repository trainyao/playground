package main

import (
	"context"
	trainyao_hostname "github.com/trainyao/sofastack_test/grpc/proto"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("svc-host1:20001", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	c := trainyao_hostname.NewHostnameClient(conn)

	r, err := c.GetHostname(context.Background(), &trainyao_hostname.HostnameRequest{Test: "123"})
	if err != nil {
		panic(err)
	}

	log.Println(r.Hostname)
}
