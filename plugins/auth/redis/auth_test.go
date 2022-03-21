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
	}, time.Minute*20)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheck(t *testing.T) {
	obj, err := New()
	if err != nil {
		t.Fatal(err)
	}
	authed, err := obj.Check(context.TODO(), &face.AuthRequest{
		ClientId: "1",
		UserName: "mqtt",
		PassWord: "public",
	}, time.Minute*20)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(authed)
}
