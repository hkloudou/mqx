package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/hkloudou/nrpc"
	"github.com/hkloudou/xtransport/packets/mqtt"
	"github.com/nats-io/nats.go"
)

// 定义命令行参数对应的变量，这三个变量都是指针类型
var cliKey = flag.String("k", "", "Input Your Key")
var cliValue = flag.String("v", "", "Input Your Value Filepath")
var mqttPublishKey = flag.String("publishkey", "mqtt.bridge.publish.jetstream", "publish_jetstream_key")

func main() {
	flag.Parse()
	c, err := nrpc.Connect("127.0.0.1")
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadFile(*cliValue)
	if err != nil {
		panic(err)
	}
	// m := nats.NewMsg(fmt.Sprintf(*cliKey))
	// m.Data = b
	pk := mqtt.NewControlPacket(mqtt.Publish).(*mqtt.PublishPacket)
	pk.Payload = b
	pk.TopicName = *cliKey
	pk.Retain = true
	var buf bytes.Buffer
	_, err = pk.WriteTo(&buf)
	if err != nil {
		panic(err)
	}
	m := nats.NewMsg(fmt.Sprintf(*mqttPublishKey))
	// m.Header.Set("x-mqtt-retain", "true")
	m.Data = buf.Bytes()
	if err := c.PublishMsg(m); err != nil {
		panic(err)
	}
	log.Println("finish", *cliKey, " - ", *cliValue)
}
