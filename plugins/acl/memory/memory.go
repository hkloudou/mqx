package memory

import (
	"fmt"
	"net"
	"strings"

	"github.com/hkloudou/mqx/face"
)

type model struct {
	User     string
	Cidr     string
	Patterns []string `delim:","`
}

type Models struct {
	PubAllow []model
	PubDeny  []model
	SubAllow []model
	SubDeny  []model
}

func init() {
	face.RegisterPugin[face.Acl]("memory", MustNew)
}

func match(meta *face.MetaInfo, arr []model, topic string, matched bool) bool {
	for i := 0; i < len(arr); i++ {
		item := arr[i]
		// match user
		if item.User != "" && item.User != meta.UserName {
			// println("not match rule:user", item.User, "my", userName)
			continue
		}
		// match cidr
		if item.Cidr != "" {
			_, _mask, err := net.ParseCIDR(item.Cidr)
			if err != nil {
				println("not match rule:cidr err", err)
				continue
			}
			if !_mask.Contains(meta.ClientIP) {
				println("not match rule:cidr", item.Cidr, "ip", meta.ClientIP.String(), meta.ClientIP.String())
				continue
			}
		}
		// match pattern
		if len(item.Patterns) == 0 {
			continue
		}
		for i := 0; i < len(item.Patterns); i++ {
			// println("ready", matched, item.Patterns[i], topic)
			pr := item.Patterns[i]
			pr = strings.ReplaceAll(pr, "<username>", meta.UserName)
			pr = strings.ReplaceAll(pr, "<clientid>", meta.ClientIdentifier)
			if matched {
				if err := face.MatchTopic(pr, topic); err == nil {
					// log.Println("matched rule", item.Patterns, topic)
					return true
				} else {
					// println("err", err.Error())
				}
			} else {
				if item.Patterns[i] == topic {
					return true
				}
			}
		}
	}
	return false
}

func MustNew(conf face.Conf) face.Acl {
	obj, err := New(conf)
	if err != nil {
		panic(err)
	}
	return obj
}

func New(conf face.Conf) (face.Acl, error) {
	obj := &memoryAcl{
		model: Models{
			PubAllow: make([]model, 0),
			PubDeny:  make([]model, 0),
			SubAllow: make([]model, 0),
			SubDeny:  make([]model, 0),
		},
	}
	loopRead := func(fm string) []model {
		i := 0
		items := make([]model, 0)
		for {
			var item model
			if err := conf.MapTo(fmt.Sprintf(fm, i), &item); err != nil {
				break
			}
			if len(item.Cidr) == 0 && len(item.Patterns) == 0 && len(item.User) == 0 {
				break
			}
			if len(item.Patterns) == 0 {
				break
			}
			items = append(items, item)
			i++
		}
		return items
	}

	obj.model.PubAllow = loopRead("acl.plugin.pub.allow.%d")
	obj.model.PubDeny = loopRead("acl.plugin.pub.deny.%d")
	obj.model.SubAllow = loopRead("acl.plugin.sub.allow.%d")
	obj.model.SubDeny = loopRead("acl.plugin.sub.deny.%d")
	return obj, nil
}

type memoryAcl struct {
	model Models
}

func (m *memoryAcl) GetSub(meta *face.MetaInfo, pattern string) face.AclCode {
	if !match(meta, m.model.SubAllow, pattern, true) {
		return face.AclCodeDeny
	}
	if match(meta, m.model.SubDeny, pattern, true) {
		return face.AclCodeDeny
	}
	return face.AclCodeAllow
}
func (m *memoryAcl) GetPub(meta *face.MetaInfo, topic string) face.AclCode {
	if !match(meta, m.model.PubAllow, topic, true) {
		return face.AclCodeDeny
	}
	if match(meta, m.model.PubDeny, topic, true) {
		return face.AclCodeDeny
	}
	return face.AclCodeAllow
}
