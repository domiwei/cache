package cache

import "time"

type cacheType int

const (
	_ cacheType = iota
	typeRedis
	typeLocal
	typeChain
)

type setting struct {
	t            cacheType
	ttlRedis     time.Duration
	ttlLocal     time.Duration
	chainSetting []settingFunc
}

type settingFunc func(s *setting)

func Local(ttl time.Duration) settingFunc {
	return func(s *setting) {
		s.t = typeLocal
		s.ttlLocal = ttl
	}
}

func Redis(ttl time.Duration) settingFunc {
	return func(s *setting) {
		s.t = typeRedis
		s.ttlRedis = ttl
	}
}

func Chain(fs ...settingFunc) settingFunc {
	return func(s *setting) {
		s.t = typeChain
		s.chainSetting = fs
	}
}

func decorateSetting(f settingFunc) *setting {
	s := setting{}
	f(&s)
	return &s
}
