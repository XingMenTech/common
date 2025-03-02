// 触发器
// Neo

package task

import (
	"time"
)

type Trigger interface {
	// CanTrigger 可触发
	CanTrigger(now time.Time) bool
	// Trigger 触发
	Trigger()
	// CanPeriodic 可周期性
	CanPeriodic() bool
}

// AnyTrigger /任意触发器
type AnyTrigger struct {
	canTrigger  func(time.Time) bool
	canPeriodic bool
	canFilter   bool
}

// NewAnyTrigger /工厂方法
func NewAnyTrigger(canTrigger func(now time.Time) bool, canPeriodic bool) *AnyTrigger {
	object := &AnyTrigger{
		canTrigger:  canTrigger,
		canPeriodic: canPeriodic,
		canFilter:   true,
	}
	return object
}

// CanTrigger /是否可以出发
func (object *AnyTrigger) CanTrigger(now time.Time) bool {
	if nil == object.canTrigger {
		return false
	}
	return object.canTrigger(now)
}

// CanPeriodic 可周期性
func (object *AnyTrigger) CanPeriodic() bool {
	return object.canPeriodic
}

// OneMinuteTrigger /1分钟周期性触发器
type OneMinuteTrigger struct {
	*AnyTrigger
}

// NewOneMinuteTrigger /工厂方法
func NewOneMinuteTrigger() *OneMinuteTrigger {
	object := &OneMinuteTrigger{}
	object.AnyTrigger = NewAnyTrigger(func(now time.Time) bool {
		return 0 == now.Second()
	}, true)
	return object
}

// NMinutesTrigger /N分钟周期性触发器
type NMinutesTrigger struct {
	*AnyTrigger
	N          int
	lastMinute int
}

// NewNMinutesTrigger /工厂方法
func NewNMinutesTrigger(n int, canPeriodic bool) *NMinutesTrigger {
	if 0 > n || 60 < n {
		panic("n must in range (0, 60)")
	}
	object := &NMinutesTrigger{
		N:          n,
		lastMinute: -1,
	}
	object.AnyTrigger = NewAnyTrigger(func(now time.Time) bool {
		if 0 == now.Minute()%n {
			if -1 == object.lastMinute || object.lastMinute != now.Minute() {
				object.lastMinute = now.Minute()
				return true
			}
		}
		return false
	}, canPeriodic)
	return object
}

// DailyTrigger /每日周期性触发器
type DailyTrigger struct {
	*AnyTrigger
}

// NewDailyTrigger /工厂方法
func NewDailyTrigger() *DailyTrigger {
	object := &DailyTrigger{}
	object.AnyTrigger = NewAnyTrigger(func(now time.Time) bool {
		return 23 == now.Hour() && 59 == now.Minute() && 59 == now.Second()
	}, true)
	return object
}

// TimePointTrigger /时间点触发
type TimePointTrigger struct {
	*AnyTrigger
	TimePoint int64
}

// NewTimePointTrigger /工厂方法
func NewTimePointTrigger(timePoint int64, filterMaster bool) *TimePointTrigger {
	object := &TimePointTrigger{
		TimePoint: timePoint,
	}
	object.AnyTrigger = NewAnyTrigger(func(now time.Time) bool {
		return object.TimePoint == now.Unix()
	}, false)
	object.AnyTrigger.canFilter = filterMaster
	return object
}
