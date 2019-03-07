package main

import (
	"context"
	"net/http"
	"net/rpc"
	"sync"

	"git.tuleu.science/fort/dieu"
	. "gopkg.in/check.v1"
)

type Hermes struct {
	C chan *C
}

func (h *Hermes) UnregisterZone(name *string, err *error) error {
	c := <-h.C
	c.Check(*name, Equals, "myself/zones/test-zone")
	*err = nil
	return nil
}

func (h *Hermes) RegisterZone(name *string, err *error) error {
	c := <-h.C
	c.Check(*name, Equals, "myself/zones/test-zone")
	*err = nil
	return nil
}

func (h *Hermes) ReportClimate(cr *dieu.NamedClimateReport, err *error) error {
	c := <-h.C
	c.Check(cr.ZoneIdentifier, Equals, "myself/zones/test-zone")
	c.Check(cr.Humidity, Equals, dieu.Humidity(50.0))
	for i := 0; i < 4; i++ {
		c.Check(cr.Temperatures[i], Equals, dieu.Temperature(21))
	}
	*err = nil
	return nil
}

const testAddress = "localhost:12345"

type RPCClimateReporterSuite struct {
	Http   *http.Server
	Rpc    *rpc.Server
	H      *Hermes
	Errors chan error
}

var _ = Suite(&RPCClimateReporterSuite{})

func (s *RPCClimateReporterSuite) SetUpSuite(c *C) {
	s.Http = &http.Server{Addr: testAddress}
	s.Rpc = rpc.NewServer()
	s.H = &Hermes{make(chan *C)}
	s.Rpc.Register(s.H)
	s.Rpc.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	s.Errors = make(chan error)
	go func() {
		err := s.Http.ListenAndServe()
		if err != http.ErrServerClosed {
			s.Errors <- err
		}
		close(s.Errors)
	}()
}

func (s *RPCClimateReporterSuite) TearDownSuite(c *C) {
	s.Http.Shutdown(context.Background())
	err, ok := <-s.Errors
	c.Check(ok, Equals, false)
	c.Check(err, IsNil)
}

func (s *RPCClimateReporterSuite) TestClimateReport(c *C) {
	go func() { s.H.C <- c }()
	n, err := NewRPCReporter("myself/zones/test-zone", testAddress)
	c.Assert(err, IsNil)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		n.Report()
		wg.Done()
	}()

	go func() { s.H.C <- c }()
	n.ReportChannel() <- dieu.ClimateReport{Humidity: 50, Temperatures: [4]dieu.Temperature{21, 21, 21, 21}}

	go func() { s.H.C <- c }()
	close(n.ReportChannel())
	close(n.AlarmChannel())
	wg.Wait()
}
