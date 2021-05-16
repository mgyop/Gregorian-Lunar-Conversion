package abase

type CalendarInterface interface {
	formatDate(bd *BaseDate) (*BaseDate, error)
	SolarCalendar(bd *BaseDate) (MonthInfo, error)
	creatMonthInfo(bd *BaseDate, days int) []BaseDate
	//将公历转换为农历
	SolarToLunar(bd *BaseDate) (BaseDate, error)
	//根据距离正月初一的天数计算农历日期
	getLunarByBetween(date *BaseDate) BaseDate
}
