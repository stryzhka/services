package word

import (
	"context"
	"time"
)

type Cache interface {
	GetOne(context.Context, string) (interface{}, error)
	GetMany(context.Context, string) ([]interface{}, error)
	Set(context.Context, string, time.Duration, interface{}) (bool, error)
	Delete(context.Context, string) (bool, error)
}
