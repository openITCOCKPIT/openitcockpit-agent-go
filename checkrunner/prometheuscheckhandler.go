package checkrunner

import (
	"context"
	"sync"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
	log "github.com/sirupsen/logrus"
)

type PrometheusExporterResult struct {
	Name   string
	Result string
}

// PrometheusCheckHandler runs prometheus exporter
type PrometheusCheckHandler struct {
	// ResultOutput channel for check results
	// Do not close before Shutdown completes
	ResultOutput  chan *PrometheusExporterResult
	Configuration []*config.PrometheusExporter

	executors []*PrometheusCheckExecutor
	shutdown  chan struct{}
	wg        sync.WaitGroup
}

// stop all prometheus exporter executors in parallel
// the cancel of the context should cause all executors to stop almost immediatly
func (c *PrometheusCheckHandler) stopExecutors() {
	if len(c.executors) < 1 {
		return
	}

	stopC := make(chan *PrometheusCheckExecutor)

	for i := 0; i < len(c.executors); i++ {
		go func() {
			for e := range stopC {
				e.Shutdown()
				log.Infoln("Prometheus Exporter ", e.Configuration.Name, " stopped")
				c.wg.Done()
			}
		}()
	}

	for _, executor := range c.executors {
		stopC <- executor
	}

	close(stopC)
}

// Run the prometheus exporters in background (DO NOT RUN IN GO ROUTINE)
func (c *PrometheusCheckHandler) Start(parentCtx context.Context) {
	c.shutdown = make(chan struct{})
	c.executors = make([]*PrometheusCheckExecutor, len(c.Configuration))

	for i, checkConfig := range c.Configuration {
		c.executors[i] = &PrometheusCheckExecutor{
			Configuration: checkConfig,
			ResultOutput:  c.ResultOutput,
		}
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		ctx, cancel := context.WithCancel(parentCtx)
		defer cancel()

		for _, executor := range c.executors {
			log.Infoln("Prometheus Exporter ", executor.Configuration.Name, " starting")
			c.wg.Add(1)
			if err := executor.Start(ctx); err != nil {
				log.Errorln(err)
			}
		}

		defer c.stopExecutors()

		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-c.shutdown:
				if !ok {
					return
				}
			}
		}

	}()
}

// Shutdown prometheus exporter runner, waits for completion
func (c *PrometheusCheckHandler) Shutdown() {
	close(c.shutdown)
	c.wg.Wait()
}
