package cache

import (
	"context"
)

type Service interface {
	Create(pfx *Key, settingFunc settingFunc) Cache
}

type Getter func() (interface{}, error)

type Cache interface {
	GetByFunc(ctx context.Context, uniqkey string, container interface{}, getter Getter) error
	Clean(context context.Context, uniqkey string) error
}
