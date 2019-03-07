package main

import (
	"log"
	"net/rpc"

	"git.tuleu.science/fort/dieu"
)

type RPCReporter struct {
	Name           string
	Conn           *rpc.Client
	ClimateReports chan dieu.ClimateReport
	AlarmReports   chan dieu.AlarmEvent
}

func (r *RPCReporter) ReportChannel() chan<- dieu.ClimateReport {
	return r.ClimateReports
}

func (r *RPCReporter) AlarmChannel() chan<- dieu.AlarmEvent {
	return r.AlarmReports
}

func (r *RPCReporter) Report() {
	var err error

	for {
		select {
		case cr, ok := <-r.ClimateReports:
			if ok == false {
				r.ClimateReports = nil
			} else {
				ncr := dieu.NamedClimateReport{cr, r.Name}
				rerr := r.Conn.Call("Hermes.ReportClimate", ncr, &err)
				if rerr != nil {
					log.Printf("[%s]: Could not transmit climate report: %s", r.Name, err)
				}
			}
		case ae, ok := <-r.AlarmReports:
			if ok == false {
				r.AlarmReports = nil
			} else {
				rerr := r.Conn.Call("Hermes.ReportAlarm", ae, &err)
				if rerr != nil {
					log.Printf("[%s]: Could not transmit climate report: %s", r.Name, err)
				}
			}
		}
		if r.AlarmReports == nil && r.ClimateReports == nil {
			break
		}
	}

	rerr := r.Conn.Call("Hermes.UnregisterZone", r.Name, &err)
	if rerr != nil {
		log.Printf("[%s]: Could not unregister zone: %s", r.Name, err)
	}
	r.Conn.Close()
}

func NewRPCReporter(name, address string) (*RPCReporter, error) {
	conn, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		return nil, err
	}

	rerr := conn.Call("Hermes.RegisterZone", name, &err)
	if rerr != nil {
		return nil, rerr
	}
	if err != nil {
		return nil, err
	}

	return &RPCReporter{
		Name:           name,
		Conn:           conn,
		ClimateReports: make(chan dieu.ClimateReport, 20),
		AlarmReports:   make(chan dieu.AlarmEvent, 20),
	}, nil
}
