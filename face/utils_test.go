package face_test

import (
	"testing"

	"github.com/hkloudou/mqx/face"
)

func et(t bool) {
	if !t {
		panic("not equal true")
	}
}

func ef(f bool) {
	if f {
		panic("not equal falase")
	}
}

func TestTopicMatch(t *testing.T) {
	et(face.MatchTopic("t/#", "t/1"))
	ef(face.MatchTopic("t/#", "t2/1"))
	et(face.MatchTopic("t/+/x", "t/1/x"))
	ef(face.MatchTopic("t/+/x", "t/1/y"))
}
