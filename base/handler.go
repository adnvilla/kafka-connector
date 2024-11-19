package base

import "context"

type ConsumerHandler func(ctx context.Context, message interface{}) error
