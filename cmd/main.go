package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/mrumpel/go-img-resizer/internal/cache"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mrumpel/go-img-resizer/internal/app"
	"github.com/mrumpel/go-img-resizer/internal/grabber"
	"github.com/mrumpel/go-img-resizer/internal/logger"
	"github.com/mrumpel/go-img-resizer/internal/resizer"
	"github.com/mrumpel/go-img-resizer/internal/server"
)

type Configuration struct {
	Host      string
	Port      string
	CacheDir  string
	CacheSize int
	LogLevel  string
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stdout, "service fail: ", err)
		os.Exit(1)
	}
}

func run() error {
	// configuration reading
	config := Configuration{}
	flag.StringVar(&(config.Host), "host", "", "service host")
	flag.StringVar(&(config.Port), "port", "8080", "service port")
	flag.StringVar(&(config.CacheDir), "cachedir", "", "directory for chached images")
	flag.IntVar(&(config.CacheSize), "cachesize", 10, "max size of the cache")
	flag.StringVar(&(config.LogLevel), "loglevel", "trace", "logging level")
	flag.Parse()

	// service parts initialize
	l := logger.NewLogger(config.LogLevel)
	g := grabber.NewGrabber()
	r := resizer.NewResizer()
	c, err := cache.NewCache(config.CacheSize, config.CacheDir)
	if err != nil {
		return fmt.Errorf("error in chache creating %w", err)
	}
	defer func() {
		err = c.Clear()
		if err != nil {
			l.Error("clearing cache error", err)
		}
	}()

	a := application.NewApp(l, g, r, c)

	// server start
	s := server.NewServer(net.JoinHostPort(config.Host, config.Port), l.GetRequestLoggingHandler()(a.GetServiceHandler()), l)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.Stop(ctx)
	}()

	s.Start()

	return nil
}
