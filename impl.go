package cache

import (
	"encoding/json"
	"reflect"

	"github.com/dgraph-io/ristretto"
)

const (
	b  = 1
	kb = 1024 * b
	mb = 1024 * kb
	gb = 1024 * mb

	// Estimated size per object is 1024 bytes
	estSizePerObj = 4 * kb
)

type impl struct {
	lc *ristretto.Cache
}

func New() Service {
	lc, err := ristretto.NewCache(&ristretto.Config{
		// Following comments come from official doc.
		// NumCounters is the number of 4-bit access counters to keep for admission
		// and eviction. We've seen good performance in setting this to 10x the
		// number of items you expect to keep in the cache when full.
		NumCounters: gb / estSizePerObj * 10,
		MaxCost:     gb,
		BufferItems: 64, // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}

	return &impl{
		lc: lc,
	}
}

func (im *impl) Create(pfx *Key, settingFunc settingFunc) Cache {
	var retCache Cache
	setting := decorateSetting(settingFunc)
	switch setting.t {
	case typeRedis:
		// todo
	case typeLocal:
		retCache = &cacheLocal{
			lc:  im.lc,
			pfx: pfx,
			ttl: setting.ttlLocal,
		}
	case typeChain:
		caches := []Cache{}
		for _, f := range setting.chainSetting {
			caches = append(caches, im.Create(pfx, f))
		}
		retCache = &cacheChain{
			caches: caches,
		}
	}
	return retCache
}

func encode(data interface{}) ([]byte, error) {
	switch v := data.(type) {
	case string:
		return []byte(v), nil
	case *string:
		return []byte(*v), nil
	case []byte:
		return v, nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return b, nil
	}
}

func decode(b []byte, container interface{}) error {
	switch container.(type) {
	case *[]byte:
		reflect.ValueOf(container).Elem().Set(reflect.ValueOf(b))
	case *string:
		reflect.ValueOf(container).Elem().Set(reflect.ValueOf(string(b)))
	default:
		if err := json.Unmarshal(b, container); err != nil {
			return err
		}
	}
	return nil
}
