package face

type ACL interface {
	SubcribeAble(userName string, topicPattern string, qos int, retain bool)
	PublishAble(userName string, topicPattern string, qos int, retain bool)
}

/*
flow:
1、check allow
2、diacard deny
*/

func x() {

}
