package main

import (
	"Jottings/tiny_rpc/log"
	"Jottings/tiny_rpc/msg"
	"Jottings/tiny_rpc/server"
)

const (
	network = "tcp"
	address = "localhost:8972"
)

func main() {
	log.SetDefLog(log.Logrus{})
	log.Info("Start..")
	defer log.Info("Stop..")
	var ser = server.NewServer(network, address)
	defer ser.Close()
	msg.SetSerializer(msg.SerializerJson)
	ser.Serve()
}
