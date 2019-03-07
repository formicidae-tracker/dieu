package main

import (
	"fmt"
	"os"
	"time"

	"git.tuleu.science/fort/dieu"
)

type ClimateReporter interface {
	C() chan<- dieu.ClimateReport
	Report()
}

type fileClimateReporter struct {
	File  *os.File
	Start time.Time
	Chan  chan dieu.ClimateReport
}

func (n *fileClimateReporter) C() chan<- dieu.ClimateReport {
	return n.Chan
}

func (n *fileClimateReporter) Report() {
	for cr := range n.Chan {
		fmt.Fprintf(n.File,
			"%d %.2f %.2f %.2f %.2f %.2f\n",
			cr.Time.Sub(n.Start).Nanoseconds()/1e6,
			cr.Humidity,
			cr.Temperatures[0],
			cr.Temperatures[1],
			cr.Temperatures[2],
			cr.Temperatures[3])
	}
	n.File.Close()
}

func NewFileClimateReporter(filename string) (ClimateReporter, string, error) {
	res := &fileClimateReporter{
		Chan:  make(chan dieu.ClimateReport, 10),
		Start: time.Now(),
	}

	var err error
	var fname string
	res.File, fname, err = dieu.CreateFileWithoutOverwrite(filename)
	if err != nil {
		return nil, "", err
	}
	fmt.Fprintf(res.File, "# Starting date %s\n# Time(ms) Relative Humidity (%%) Temperature (°C) Temperature (°C) Temperature (°C) Temperature (°C)\n", res.Start)

	return res, fname, nil
}

func NewTCPClimateReporterNotifier(address string) (ClimateReporter, error) {
	return nil, nil
}
