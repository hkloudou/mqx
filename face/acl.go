package face

import "context"

type Acl interface {
	Subcribe(ctx context.Context, meta *MetaInfo, pattern string) (bool, error)
	Publish(ctx context.Context, meta *MetaInfo, topic string) (bool, error)
}
