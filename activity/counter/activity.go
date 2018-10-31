package counter

import (
	"sync/atomic"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

const (
	ovValue = "value"
)

var counters = make(map[string]*Counter)

type CounterFunc func() uint64

type Settings struct {
	CounterName string `md:"counterName,required"`
	Op          string `md:"op,allowed(get,increment,reset)"`
}

type Output struct {
	Value string `md:"value"`
}

func init() {
	activity.Register(&CounterActivity{}, New)
}

var activityMd = activity.ToMetadata(&Settings{}, &Output{})

// CounterActivity is a Counter CounterActivity implementation
type CounterActivity struct {
	invoke CounterFunc
}

func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	act := &CounterActivity{}

	counter, exists := counters[s.CounterName]

	if !exists {
		//log creating counter
		counter = &Counter{val: 0}
		counters[s.CounterName] = counter
	}

	switch s.Op {
	case "increment":
		act.invoke = counter.Increment
	case "reset":
		act.invoke = counter.Reset
	default:
		act.invoke = counter.Get
	}

	return act, nil
}

// Metadata implements activity.CounterActivity.Metadata
func (a *CounterActivity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements activity.CounterActivity.Eval
func (a *CounterActivity) Eval(context activity.Context) (done bool, err error) {
	val := a.invoke()

	context.SetOutput(ovValue, int(val))

	return true, nil
}

type Counter struct {
	val uint64
}

func (c *Counter) Get() uint64 {
	return atomic.LoadUint64(&c.val)
}

func (c *Counter) Increment() uint64 {
	return atomic.AddUint64(&c.val, 1)
}

func (c *Counter) Reset() uint64 {
	atomic.StoreUint64(&c.val, 0)
	return 0
}
