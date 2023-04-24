package jobCounter

import (
	"github.com/3th1nk/easygo/util/strUtil"
	"sync"
	"sync/atomic"
	"time"
)

func New(name string) *Counter {
	lock.Lock()
	defer lock.Unlock()
	obj := &Counter{name: name, cid: strUtil.Rand(6)}
	allCounter = append(allCounter, obj)
	return obj
}

type Counter struct {
	name string
	cid  string
	jobs sync.Map
	size int32
}

func AllCounter() []*Counter {
	return allCounter
}

var (
	allCounter = make([]*Counter, 0, 16)
	lock       = sync.Mutex{}
)

func (this *Counter) ClientId() string {
	return this.cid
}

func (this *Counter) Name() string {
	return this.name
}

func (this *Counter) Size() int {
	return int(this.size)
}

func (this *Counter) Add(jobId string, state ...string) {
	var theState string
	if len(state) != 0 {
		theState = state[0]
	}
	atomic.AddInt32(&this.size, 1)
	this.jobs.Store(jobId, &jobData{Time: time.Now(), State: theState})
}

func (this *Counter) SetState(jobId, state string) {
	if val, _ := this.jobs.Load(jobId); val != nil {
		val.(*jobData).State = state
	}
}

func (this *Counter) Done(jobId string) {
	this.jobs.Delete(jobId)
	atomic.AddInt32(&this.size, -1)
}

func (this *Counter) RunningJobs() (jobs []*Job) {
	jobs = make([]*Job, 0, this.size)
	this.jobs.Range(func(key, val interface{}) bool {
		jobs = append(jobs, &Job{JobId: key.(string), Time: val.(*jobData).Time, State: val.(*jobData).State})
		return true
	})
	return
}

type jobData struct {
	Time  time.Time
	State string
}

type Job struct {
	JobId string    `json:"job_id,omitempty"`
	Time  time.Time `json:"time,omitempty"`
	State string    `json:"state,omitempty"`
}
