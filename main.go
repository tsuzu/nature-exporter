package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tenntenn/natureremo"
)

var (
	temperatureGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "tsuzu_room_temperature",
		Help: "The current temperature in Tsuzu's room",
	})
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cli := natureremo.NewClient(os.Getenv("TOKEN"))

	go func() {
		for {
			time.Sleep(10 * time.Second)
			select {
			case <-ctx.Done():
				return
			default:
			}

			devices, err := cli.DeviceService.GetAll(ctx)

			if err != nil {
				log.Println(err)

				continue
			}

			for _, dev := range devices {
				temperatureGauge.Set(dev.NewestEvents[natureremo.SensorTypeTemperature].Value)
			}
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	server := http.Server{
		Addr: ":2112",
	}
	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	server.ListenAndServe()
}
