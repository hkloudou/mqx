package face

type ACL interface {
	SubcribeAble(userName string, topicPattern string)
	PublishAble(userName string, topicPattern string)
}
