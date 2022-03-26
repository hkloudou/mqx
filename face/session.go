package face

import "context"

type Session interface {
	Add(ctx context.Context, clientid string, patterns ...string) error
	Remove(ctx context.Context, clientid string, patterns ...string) error
	Clear(ctx context.Context, cliendid string) error
	Match(ctx context.Context, topic string) ([]string, error)
	ClientPatterns(ctx context.Context, cliendid string) ([]string, error)
}
