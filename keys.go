package cache

import "strings"

type CachePrefix string

const (
	PfxCacheService CachePrefix = "cache:service"
)

func (p CachePrefix) NewKey(ks ...string) *Key {
	k := &Key{
		pkgPrefix: p,
		keys:      []string{},
	}
	k.keys = append(k.keys, ks...)
	return k
}

type Key struct {
	pkgPrefix CachePrefix
	keys      []string
	subKey    *Key
}

func (k *Key) AppendKeys(ks ...string) *Key {
	k.keys = append(k.keys, ks...)
	return k
}

func (k *Key) Wrap(subk *Key) *Key {
	k.subKey = subk
	return k
}

func (k *Key) WrappedBy(p CachePrefix) *Key {
	return p.NewKey().Wrap(k)
}

func (k *Key) ToKey() string {
	finalkey := []string{string(k.pkgPrefix)}
	finalkey = append(finalkey, k.keys...)
	if k.subKey != nil {
		finalkey = append(finalkey, k.subKey.ToKey())
	}
	return JoinKey(finalkey...)
}

/*
	func CacheKey(pfx CachePrefix, keys ...string) string {
		return string(pfx) + ":" + JoinKey(keys...)
	}
*/
func JoinKey(keys ...string) string {
	return strings.Join(keys, ":")
}
