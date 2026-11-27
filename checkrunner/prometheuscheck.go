package checkrunner

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
	log "github.com/sirupsen/logrus"
)

type PrometheusCheckExecutor struct {
	Configuration *config.PrometheusExporter
	ResultOutput  chan *PrometheusExporterResult

	wg       sync.WaitGroup
	shutdown chan struct{}
}

func (c *PrometheusCheckExecutor) Shutdown() {
	close(c.shutdown)
	c.wg.Wait()
}

func (c *PrometheusCheckExecutor) runCheck(ctx context.Context, timeout time.Duration) {
	log.Debugln("Begin Prometheus Exporter: ", c.Configuration.Name)

	client := &http.Client{
		Timeout: timeout,
	}

	scheme := "http"
	if c.Configuration.Method == "https" {
		scheme = "https"
	}

	url := fmt.Sprintf("%s://%s:%d%s", scheme, "localhost", c.Configuration.Port, c.Configuration.Path)
	resp, err := client.Get(url)
	if err != nil {
		log.Errorf("Prometheus Exporter '%s' error: %v", c.Configuration.Name, err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Prometheus Exporter '%s' error reading response body: %v", c.Configuration.Name, err)
		return
	}

	select {
	// Return Prometheus Exporter result to Agent Instance
	case c.ResultOutput <- &PrometheusExporterResult{
		Name:   c.Configuration.Name,
		Result: string(body),
	}:
	case <-time.After(time.Second * 5):
		log.Errorln("Internal error: timeout could not save Prometheus Exporter result")
	case <-c.shutdown:
		log.Errorln("Prometheus Exporter: canceled")
		return
	case <-ctx.Done():
		log.Errorln("Prometheus Exporter: canceled")
		return
	}
	log.Debugln("Finish Prometheus Exporter: ", c.Configuration.Name)
}

func (c *PrometheusCheckExecutor) Start(parent context.Context) error {
	c.shutdown = make(chan struct{})
	timeout := time.Duration(c.Configuration.Timeout) * time.Second
	interval := time.Duration(c.Configuration.Interval) * time.Second

	if timeout > interval {
		return fmt.Errorf("prometheus exporter '%s': timeout must be lower or equal to interval", c.Configuration.Name)
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		ctx, cancel := context.WithCancel(parent)
		defer cancel()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		c.runCheck(ctx, timeout)
		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-c.shutdown:
				if !ok {
					return
				}
			case <-ticker.C:
				c.runCheck(ctx, timeout)
			}
		}
	}()

	return nil
}
