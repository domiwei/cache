package cache

import (
	"context"
	"log"
)

type cacheChain struct {
	caches []Cache
}

func (c *cacheChain) GetByFunc(
	context context.Context,
	uniqkey string,
	container interface{},
	getter Getter) error {

	nowgetter := getter
	for i := range c.caches {
		innercache := c.caches[len(c.caches)-i-1]
		// NOTE: Must need to re-assign nowgetter to another local variable in each loop iteration,
		// or it may lead in refrence to each other.
		innergetter := nowgetter
		nowgetter = func() (interface{}, error) {
			if err := innercache.GetByFunc(context, uniqkey, container, innergetter); err != nil {
				log.Printf("innercache.GetByFunc failed: %v", err)
				return nil, err
			}
			return container, nil
		}
	}

	if _, err := nowgetter(); err != nil {
		log.Printf("nowgetter failed: %v", err)
		return err
	}
	return nil
}

func (c *cacheChain) Clean(
	context context.Context,
	uniqkey string) error {

	for i, c := range c.caches {
		if err := c.Clean(context, uniqkey); err != nil {
			log.Printf("caches[%d].Clean failed: %v", i, err)
			// Do not return error
		}
	}
	return nil
}
