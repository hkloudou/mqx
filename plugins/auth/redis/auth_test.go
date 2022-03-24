package redis

import (
	"context"
	"testing"

	"github.com/hkloudou/mqx/face"
)

func TestAuth(t *testing.T) {
	// t.Log("")
	obj, err := New(nil)
	if err != nil {
		t.Fatal(err)
	}
	err = obj.Update(context.TODO(), &face.AuthRequest{
		ClientId: "asdas-xasd-123-axsd",
		UserName: "user1",
		PassWord: "pwd",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheck(t *testing.T) {
	// encoding.BinaryUnmarshaler
	// t.Log("xx")
	obj, err := New(nil)
	if err != nil {
		t.Fatal(err)
	}
	code := obj.Check(context.TODO(), &face.AuthRequest{
		ClientId: "1",
		UserName: "mqtt",
		PassWord: "public",
	}, face.WithAuthRequestTtl(-1))

	if err != nil {
		t.Fatal(err)
	}
	t.Log(code)
}

func Test_Expired(t *testing.T) {
	obj, err := New(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(obj.(*redisAuther).expiredBeforeConnection(&face.AuthRequest{
		ClientId: "1",
		UserName: "mqtt",
		PassWord: "public",
	}, 1, face.AuthDiscardOld))
}
