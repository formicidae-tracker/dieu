package main

import (
	"fmt"
	"os"
	"time"

	"github.com/formicidae-tracker/zeus"
)

type ClimateReporter interface {
	ReportChannel() chan<- zeus.ClimateReport
}

type FileClimateReporter struct {
	File  *os.File
	Start time.Time
	Chan  chan zeus.ClimateReport
}

func (n *FileClimateReporter) ReportChannel() chan<- zeus.ClimateReport {
	return n.Chan
}

func (n *FileClimateReporter) Report() {
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

func NewFileClimateReporter(filename string) (*FileClimateReporter, string, error) {
	res := &FileClimateReporter{
		Chan:  make(chan zeus.ClimateReport, 10),
		Start: time.Now(),
	}

	var err error
	var fname string
	res.File, fname, err = zeus.CreateFileWithoutOverwrite(filename)
	if err != nil {
		return nil, "", err
	}
	fmt.Fprintf(res.File, "# Starting date %s\n# Time(ms) Relative Humidity (%%) Temperature (°C) Temperature (°C) Temperature (°C) Temperature (°C)\n", res.Start)

	return res, fname, nil
}