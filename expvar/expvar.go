// 对 golang 内置的 expvar 的扩展
package expvar

import (
	"github.com/3th1nk/easygo/util/mapUtil"
	"github.com/3th1nk/easygo/util/mathUtil"
	"github.com/3th1nk/easygo/util/strUtil"
	"sort"
	"sync"
)

var Default = New()

func New() VarSet {
	return &varSetImpl{varMap: make(map[string]interface{}, 64)}
}

type VarSet interface {
	Publish(name string, f VarFunc)
	PublishMap(f VarMapFunc)
	All() mapUtil.StringObjectMap
	Each(f func(key string, val interface{}), sortKeys ...int)
}

type varSetImpl struct {
	lock   sync.RWMutex
	varMap map[string]interface{}
}

type VarFunc func() interface{}

type VarMapFunc func() (vars map[string]interface{})

func IntVarFunc(f func() int) VarFunc {
	return func() interface{} { return f() }
}

func (this *varSetImpl) Publish(name string, f VarFunc) {
	this.lock.Lock()
	this.varMap[name] = f
	this.lock.Unlock()
}

func (this *varSetImpl) PublishMap(f VarMapFunc) {
	this.lock.Lock()
	this.varMap[strUtil.Rand(8)] = f
	this.lock.Unlock()
}

func (this *varSetImpl) All() (data mapUtil.StringObjectMap) {
	data = make(mapUtil.StringObjectMap, mathUtil.MaxInt(len(this.varMap), 64))

	this.lock.RLock()
	defer this.lock.RUnlock()

	for key, val := range this.varMap {
		if f, ok := val.(VarFunc); ok {
			data[key] = f()
		} else if f, ok := val.(VarMapFunc); ok {
			m := f()
			for k, v := range m {
				data[k] = v
			}
		}
	}

	return
}

func (this *varSetImpl) Each(f func(key string, val interface{}), sortKeys ...int) {
	data := this.All()
	if len(sortKeys) == 0 || sortKeys[0] == 0 {
		for key, val := range data {
			f(key, val)
		}
	} else {
		keys := data.Keys()
		sort.Strings(keys)
		if sortKeys[0] < 0 {
			sort.Sort(sort.Reverse(sort.StringSlice(keys)))
		}
		for _, key := range keys {
			f(key, data[key])
		}
	}
}
