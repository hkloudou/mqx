package face

import "time"

type AuthDiscardPolicy int

const (
	// DiscardOld will remove older user expired
	// the default.
	AuthDiscardOld AuthDiscardPolicy = iota
	//DiscardNew will fail to accept new auth
	AuthDiscardNew
)

func DefaultAuthRequestOptions() authRequestOptions {
	return authRequestOptions{
		MaxTokens: 0,
		Discard:   AuthDiscardOld,
	}
}

type AuthRequestOption func(*authRequestOptions) error

type authRequestOptions struct {
	Ttl       time.Duration
	UseTtl    bool
	MaxTokens uint64 //no limit
	Discard   AuthDiscardPolicy
}

func WithAuthRequestTtl(ttl time.Duration) AuthRequestOption {
	return func(o *authRequestOptions) error {
		o.Ttl = ttl
		o.UseTtl = true
		return nil
	}
}

func WithAuthRequestMaxTokens(val uint64) AuthRequestOption {
	return func(o *authRequestOptions) error {
		o.MaxTokens = val
		return nil
	}
}

func WithAuthRequestDiscardPolicy(val AuthDiscardPolicy) AuthRequestOption {
	return func(o *authRequestOptions) error {
		o.Discard = val
		return nil
	}
}
