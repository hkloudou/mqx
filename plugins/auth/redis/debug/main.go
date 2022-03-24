package main

import (
	"context"
	"time"

	"github.com/hkloudou/mqx/face"
	"github.com/hkloudou/mqx/plugins/auth/redis"
)

func main() {
	obj, err := redis.New(nil)
	if err != nil {
		panic(err)
	}
	go obj.MotionExpired(func(userName, clientId string) error {
		println("userName", userName)
		println("clientId", clientId)
		return nil
	})

	time.Sleep(2 * time.Second)
	err = obj.Update(context.TODO(), &face.AuthRequest{
		ClientId: "1",
		UserName: "mqtt",
		PassWord: "public",
	}, face.WithAuthRequestTtl(10*time.Second))
	if err != nil {
		panic(err)
	}
	// log.Println("wt")

	// time.Sleep(1 * time.Second)
	// err = obj.Update(context.TODO(), &face.AuthRequest{
	// 	ClientId: "1",
	// 	UserName: "mqtt",
	// 	PassWord: "public",
	// }, face.WithAuthRequestTtl(-1))
	// if err != nil {
	// 	panic(err)
	// }
	// time.Sleep(1 * time.Second)
	<-make(chan bool)
}
