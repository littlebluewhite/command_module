package schedule

import "time"

func (s *Schedule) TimeActive(t time.Time) (result bool) {
	result = s.TimeData.AccessTimeDate(t) && s.Enabled
	return
}

func (s *Schedule) GetID() (n int) {
	n = s.ID
	return
}

func (s *Schedule) Execute() (err error) {
	err = s.Command.Execute()
	return
}
