package timeRoundedCounter

import (
	"sort"
	"sync"
	"time"
)

type CounterInterface interface {
	Name() string
	Round() time.Duration
	GetAll(rtrimZero ...bool) []CounterItem
}

type CounterItem interface {
	IsZero() bool
	Reset()
}

type PropertyCounterInterface interface {
	CounterInterface

	ItemProperties() []string
}

type PropertyCounterItem interface {
	CounterItem

	Values() []interface{}
}

var (
	namedCounter     = make([]*counterWrap, 0, 4)
	namedCounterLock = sync.RWMutex{}
)

type counterWrap struct {
	name    string
	counter CounterInterface
}

func AllNamedCounter() []CounterInterface {
	src := namedCounter
	dst := make([]CounterInterface, len(src))
	for i, v := range src {
		dst[i] = v.counter
	}
	return dst
}

func registerNamedCounter(c CounterInterface) {
	name := c.Name()
	if name == "" {
		return
	}

	namedCounterLock.Lock()
	defer namedCounterLock.Unlock()

	for _, v := range namedCounter {
		if v.name == name {
			return
		}
	}

	namedCounter = append(namedCounter, &counterWrap{name: name, counter: c})
	sort.Slice(namedCounter, func(i, j int) bool { return namedCounter[i].name < namedCounter[j].name })
}

func toPropertyCounterItem(src []CounterItem) []PropertyCounterItem {
	dst := make([]PropertyCounterItem, len(src))
	for i, v := range src {
		dst[i] = v.(PropertyCounterItem)
	}
	return dst
}
