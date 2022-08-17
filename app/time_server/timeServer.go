package time_server

import (
	"context"
	"fmt"
	"github.com/patrickmn/go-cache"
	"new_command/pkg/logFile"
	"time"
)

var (
	timeServerLog logFile.LogFile
)

func init() {
	timeServerLog = logFile.NewLogFile("app", "time_server")
}

type TimeServer interface {
	Start(ctx context.Context)
}

type timeServer[T any] struct {
	cache      *cache.Cache
	cacheTable string
	duration   time.Duration
}

func NewTimeServer[T any](cache *cache.Cache, tableName string, duration time.Duration) TimeServer {
	return &timeServer[T]{
		cache:      cache,
		cacheTable: tableName,
		duration:   duration,
	}
}

func (ts *timeServer[T]) Start(ctx context.Context) {
	timeServerLog.Info().Println("Time server start")
	ticker := time.NewTicker(ts.duration)
	ctxChild, cancel := context.WithCancel(ctx)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			timeServerLog.Info().Println("Time server stop gracefully")
			return
		case t := <-ticker.C:
			s := ts.refreshData()
			for _, value := range s {
				go func(value T) {
					pointer := &value
					listenToSchedule(ctxChild, pointer, t)
				}(value)
			}
			fmt.Println("Invoked at ", t)
		}
	}
}

func (ts *timeServer[T]) refreshData() (s map[int]T) {
	if x, found := ts.cache.Get(ts.cacheTable); found {
		s = x.(map[int]T)
	}
	return
}
