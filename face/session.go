package face

import "context"

type Session interface {
	Add(ctx context.Context, clientid string, patterns ...string) error
	Remove(ctx context.Context, clientid string, patterns ...string) error
	List(ctx context.Context, pattern string) ([]string, error)
	// Match(ctx context.Context, pattern string) ([]string, error)
	Patterns(ctx context.Context) ([]string, error)
	Clear(ctx context.Context, cliendid string)
}
