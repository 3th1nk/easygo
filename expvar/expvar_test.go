package expvar

import (
	"fmt"
	"github.com/3th1nk/easygo/util/mapUtil"
	"github.com/shirou/gopsutil/process"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestVarSet(t *testing.T) {
	vars := New()
	vars.Publish("goroutine", IntVarFunc(runtime.NumGoroutine))
	vars.PublishMap(getProcessInfo)

	end := time.Now().Add(10 * time.Second)
	for range time.Tick(time.Second) {
		fmt.Println("------------------------------------------")
		vars.Each(func(key string, val interface{}) {
			fmt.Println(fmt.Sprintf("    %s: %v", key, val))
		}, 1)
		if time.Now().After(end) {
			break
		}
	}
}

func getProcessInfo() map[string]interface{} {
	dict := mapUtil.StringObjectMap{}

	p, _ := process.NewProcess(int32(os.Getpid()))
	if p == nil {
		return nil
	}

	{
		val, _ := p.Cwd()
		dict["process.cwd"] = val
	}
	{
		val, _ := p.NumFDs()
		dict["process.num_fd"] = val
	}
	{
		val, _ := p.NumThreads()
		dict["process.num_threads"] = val
	}

	if pm, _ := p.MemoryInfo(); pm != nil {
		dict.SetMulti(map[string]interface{}{
			"process.mem_rss":  pm.RSS,
			"process.mem_vms":  pm.VMS,
			"process.mem_swap": pm.Swap,
		})
	}

	if pt, _ := p.Times(); pt != nil {
		dict.SetMulti(map[string]interface{}{
			"process.time_cpu":    pt.CPU,
			"process.time_idle":   pt.Idle,
			"process.time_iowait": pt.Iowait,
		})
	}

	return dict
}
