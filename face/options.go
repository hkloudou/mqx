package face

type AuthDiscardPolicy int

const (
	// DiscardOld will remove older user expired
	// the default.
	AuthDiscardOld AuthDiscardPolicy = iota
	//DiscardNew will fail to accept new auth
	AuthDiscardNew
)

type AuthOptionConfiger struct {
	GOpt *authOptions
}

func (m *AuthOptionConfiger) Config(options ...authOption) error {
	opts := defaultAuthOptions()
	for _, opt := range options {
		if opt != nil {
			if err := opt(&opts); err != nil {
				return err
			}
		}
	}
	m.GOpt = &opts
	return nil
}

type authOption func(*authOptions) error

type authOptions struct {
	// Prefix    string
	MaxTokens uint64
	Discard   AuthDiscardPolicy
}

func defaultAuthOptions() authOptions {
	return authOptions{
		// Prefix:    "mqtt.auth",
		MaxTokens: 0, //no limit
		Discard:   AuthDiscardOld,
	}
}

func WithAuthMaxTokens(val uint64) authOption {
	return func(o *authOptions) error {
		o.MaxTokens = val
		return nil
	}
}
