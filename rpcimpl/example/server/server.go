package main

import (
	"tiny_rpc/log"
	"tiny_rpc/server"
)

const (
	network = "tcp"
	address = "localhost:8972"
)

func main() {
	log.Init(log.Logrus)
	log.Info("Start..")
	defer log.Info("Stop..")

	var ser = server.NewServer(network, address)
	defer ser.Close()
	ser.Serve()
}
