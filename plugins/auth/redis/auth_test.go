package redis

import (
	"context"
	"testing"
	"time"

	"github.com/hkloudou/mqx/face"
)

func TestAuth(t *testing.T) {
	obj, err := New()
	if err != nil {
		t.Fatal(err)
	}
	err = obj.Update(context.TODO(), &face.AuthRequest{
		ClientId: "1",
		UserName: "mqtt",
		PassWord: "public",
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)

	err = obj.Update(context.TODO(), &face.AuthRequest{
		ClientId: "2",
		UserName: "mqtt",
		PassWord: "public",
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	err = obj.Update(context.TODO(), &face.AuthRequest{
		ClientId: "3",
		UserName: "mqtt",
		PassWord: "public",
	}, face.WithAuthRequestDiscardPolicy(face.AuthDiscardOld), face.WithAuthRequestMaxTokens(1))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheck(t *testing.T) {
	// encoding.BinaryUnmarshaler
	// t.Log("xx")
	obj, err := New()
	if err != nil {
		t.Fatal(err)
	}
	authed, err := obj.Check(context.TODO(), &face.AuthRequest{
		ClientId: "1",
		UserName: "mqtt",
		PassWord: "public",
	}, face.WithAuthRequestTtl(-1))

	if err != nil {
		t.Fatal(err)
	}
	t.Log(authed)
}

func Test_Expired(t *testing.T) {
	obj, err := New()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(obj.(*redisAuther).expiredBeforeConnection(&face.AuthRequest{
		ClientId: "1",
		UserName: "mqtt",
		PassWord: "public",
	}, 1, face.AuthDiscardOld))
}
