package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bibaroc/paladinswr/app/wrsvc"
	"github.com/go-kit/kit/log"
	"github.com/oklog/oklog/pkg/group"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	inflixDBConfig := wrsvc.GetWriteAPIConfig()
	assetHandler, cancelFunc := wrsvc.CachedHTTPHandler(logger, inflixDBConfig)

	defer cancelFunc()

	var g group.Group
	{ // hello service
		httpListener, err := net.Listen("tcp", ":8080")
		if err != nil {
			_ = logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			return
		}

		g.Add(func() error {
			app := http.NewServeMux()
			app.HandleFunc("/stats", assetHandler.GetStats)

			_ = logger.Log("transport", "HTTP", "addr", ":8080")
			return http.Serve(httpListener, app)
		}, func(error) {
			httpListener.Close()
		})
	}

	{ // metrics
		httpListener, err := net.Listen("tcp", ":9100")
		if err != nil {
			_ = logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			return
		}
		g.Add(func() error {
			app := http.NewServeMux()
			app.Handle("/metrics", promhttp.Handler())

			_ = logger.Log("transport", "HTTP", "addr", ":9100")
			return http.Serve(httpListener, app)
		}, func(error) {
			httpListener.Close()
		})
	}
	{ // graceful shutdown
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	_ = logger.Log("exit", g.Run())
}
