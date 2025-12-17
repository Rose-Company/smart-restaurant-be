package prometheus

import (
	logger2 "app-noti/services/logger"
	"context"
	"fmt"
	"net/http"

	prm "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type prometheusServer struct {
	prefix string
	path   string
	port   int
	logger logger2.Loggers
}

func NewPrometheusServer(prefix string, path string, port int, logger logger2.Loggers) *prometheusServer {
	return &prometheusServer{prefix: prefix, path: path, port: port, logger: logger}
}

func (p *prometheusServer) Run() error {
	go func() {
		http.Handle(p.path, promhttp.Handler())
		p.logger.Info().Println(fmt.Sprintf("%v is running in port %v", p.prefix, p.port))
		if err := http.ListenAndServe(fmt.Sprintf(":%v", p.port), nil); err != nil {
			fmt.Println("err", err)
		}
	}()
	return nil
}

func (p *prometheusServer) RegisterMetrics(ctx context.Context, collector prm.Collector) {
	prm.Register(collector)
}

func (u *prometheusServer) GetPrefix() string {
	return u.prefix
}

func (u *prometheusServer) Stop() <-chan bool {
	stop := make(chan bool)
	go func() {
		stop <- true
	}()
	return stop
}

func (u *prometheusServer) Get() interface{} {
	return u.port
}
