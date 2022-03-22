package face

func DefaultAclFilterOptions() authRequestOptions {
	return authRequestOptions{
		MaxTokens: 0,
		Discard:   AuthDiscardOld,
	}
}

type ACLFilterPolicy uint8

const (
	ACLFilterPolicyUnknow ACLFilterPolicy = iota
	ACLFilterPolicyAllow
	ACLFilterPolicyDeny
)

type AclFilterOption func(*aclFilterOptions) error

type aclFilterOptions struct {
	// Publish  ACLFilterPolicy
	// Subcribe ACLFilterPolicy
	UserName []string //admin@*		allow regexp
	//192.168.0.0/24   2408:4321:180:1701:94c7:bc38:3bfa:***/128
	//0.0.0.0/0（IPv4）::/0（IPv6）
	CIDR   []string
	Topics []string //mqtt format topic to regexp
	MaxQos uint8    //allowd max qos
}

func WithAclFilterUserName(vals ...string) AclFilterOption {
	return func(o *aclFilterOptions) error {
		o.UserName = vals
		return nil
	}
}

func WithAclFilterMaxQos(vals ...string) AclFilterOption {
	return func(o *aclFilterOptions) error {
		o.UserName = vals
		return nil
	}
}
