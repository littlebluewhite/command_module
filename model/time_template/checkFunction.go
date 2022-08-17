package time_template

func (tt *TimeTemplate) CheckTimeTemplate() (err error) {
	err = tt.TimeData.CheckTimeData()
	return
}
