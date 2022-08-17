package schedule

func (s *Schedule) CheckSchedule() (err error) {
	err = s.TimeData.CheckTimeData()
	return
}
