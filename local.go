package cache

import (
	"context"
	"errors"
	"log"
	"reflect"
	"time"

	"github.com/dgraph-io/ristretto"
)

type cacheLocal struct {
	lc  *ristretto.Cache
	pfx *Key
	ttl time.Duration
}

func (c *cacheLocal) GetByFunc(
	context context.Context,
	uniqkey string,
	container interface{},
	getter Getter) error {

	var b []byte
	key := PfxCacheService.NewKey(uniqkey).Wrap(c.pfx).ToKey()
	//key := models.CacheKey(models.PfxCacheService, c.pfx, uniqkey)
	if intf, exist := c.lc.Get(key); exist {
		var ok bool
		b, ok = intf.([]byte)
		if !ok {
			log.Printf("type of interface is not []byte. type: %v", reflect.TypeOf(intf))
			return errors.New("type of interface is not []byte")
		}
	} else {
		// In case of data not found, invoke getter and set into redis
		data, err2 := getter()
		if err2 != nil {
			log.Printf("getter failed: %v", err2)
			return err2
		}
		if bytes, err := encode(data); err != nil {
			log.Printf("encode failed: %v", err)
			return err
		} else {
			b = bytes
		}
		// Do not care if this gets succeeded or not, because this library has no gaurantee about success of Set Op.
		c.lc.SetWithTTL(key, b, int64(len(b)), c.ttl)
	}

	// Write back to given container
	if err := decode(b, container); err != nil {
		log.Printf("decode failed: %v", err)
		return err
	}

	return nil
}

func (c *cacheLocal) Clean(
	context context.Context,
	uniqkey string) error {
	key := PfxCacheService.NewKey(uniqkey).Wrap(c.pfx).ToKey()
	c.lc.Del(key)
	return nil
}
