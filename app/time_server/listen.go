package time_server

import (
	"context"
	"time"
)

type listener interface {
	TimeActive(t time.Time) bool
	GetID() int
	Execute() error
}

// listenToSchedule 判斷不同的value type決定如何處理
func listenToSchedule(ctx context.Context, pointer any, t time.Time) {
	select {
	case <-ctx.Done():
		return
	default:
		s := pointer.(listener)
		if s.TimeActive(t) {
			if err := s.Execute(); err != nil {
				timeServerLog.Error().Printf("Error!, Schedule id: %d err: %v\n",
					s.GetID(), err)
			} else {
				timeServerLog.Info().Printf("Schedule id: %d, done\n", s.GetID())
			}
		}
	}
}
