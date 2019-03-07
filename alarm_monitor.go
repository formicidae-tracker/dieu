package main

import (
	"os"
	"path"
	"time"
)

type AlarmStatus int

const (
	AlarmOn AlarmStatus = iota
	AlarmOff
)

type AlarmEvent struct {
	Zone   string
	Alarm  Alarm
	Status AlarmStatus
	Time   time.Time
}

type AlarmMonitor interface {
	Name() string
	Monitor()
	Inbound() chan<- Alarm
	Outbound() <-chan AlarmEvent
}

type alarmMonitor struct {
	inbound  chan Alarm
	outbound chan AlarmEvent
	name     string
}

func (m *alarmMonitor) Name() string {
	return m.name
}

func wakeupAfter(wakeup chan<- string, quit <-chan struct{}, reason string, after time.Duration) {
	time.Sleep(after)
	select {
	case <-quit:
		return
	default:
	}
	wakeup <- reason
}

func (m *alarmMonitor) Monitor() {
	trigged := make(map[string]int)
	alarms := make(map[string]Alarm)
	clearWakeup := make(chan string)

	defer func() {
		close(m.outbound)
		close(clearWakeup)
	}()

	quit := make(chan struct{})

	for {
		select {
		case a, ok := <-m.inbound:
			if ok == false {
				close(quit)
				return
			}
			if t, ok := trigged[a.Reason()]; ok == false || t <= 0 {
				go func() {
					m.outbound <- AlarmEvent{
						Alarm:  a,
						Status: AlarmOn,
						Time:   time.Now(),
						Zone:   m.name,
					}
				}()
				trigged[a.Reason()] = 1
				alarms[a.Reason()] = a
			} else {
				trigged[a.Reason()] += 1
			}
			go wakeupAfter(clearWakeup, quit, a.Reason(), 3*a.RepeatPeriod())
		case r := <-clearWakeup:
			t, ok := trigged[r]
			if ok == false {
				// should not happen but lets says it does
				continue
			}

			if t == 1 {
				go func() {
					m.outbound <- AlarmEvent{
						Alarm:  alarms[r],
						Status: AlarmOff,
						Time:   time.Now(),
						Zone:   m.name,
					}
				}()
			}
			if t != 0 {
				trigged[r] = t - 1
			}
		}
	}
}

func (m *alarmMonitor) Inbound() chan<- Alarm {
	return m.inbound
}

func (m *alarmMonitor) Outbound() <-chan AlarmEvent {
	return m.outbound
}

func NewAlarmMonitor(zoneName string) (AlarmMonitor, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &alarmMonitor{
		inbound:  make(chan Alarm, 30),
		outbound: make(chan AlarmEvent, 60),
		name:     path.Join(hostname, "zones", zoneName),
	}, nil
}
