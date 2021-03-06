package abase

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CalendarConfig struct {
	MinYear, MaxYear                                                   int
	HeavenlyStems, EarthlyBranches, Zodiac, SolarTerm, MonthCn, DateCn []string
	LunarFestival                                                      map[string]string
	LunarInfo                                                          [][4]int
}

func Init() CalendarConfig {
	var cc CalendarConfig
	cc.MinYear = 1890
	cc.MaxYear = 2100
	//格式化农历时间月
	MonthCn := []string{"正", "二", "三", "四", "五", "六", "七", "八", "九", "十", "十一", "十二"}
	cc.MonthCn = MonthCn
	//格式化农历时间天
	DateCn := []string{"初一", "初二", "初三", "初四", "初五", "初六", "初七", "初八", "初九", "初十", "十一", "十二", "十三", "十四", "十五", "十六", "十七", "十八", "十九", "二十", "廿一", "廿二", "廿三", "廿四", "廿五", "廿六", "廿七", "廿八", "廿九", "三十", "卅一"}
	cc.DateCn = DateCn
	//初始化农历数据
	InitLunarInfo(&cc)
	return cc
}

//统一的日期格式化入口 输入月份从1开始，内部月份统一从0开始
func (cc *CalendarConfig) formatDate(bd *BaseDate) (*BaseDate, error) {
	if bd.Year < (cc.MinYear+1) || bd.Year > cc.MaxYear {
		return nil, errors.New("输入的年份超过了可查询范围，仅支持1891至2100年")
	}
	bd.Month = time.Month(int(bd.Month) - 1)
	return bd, nil
}

//统一的日期格式化入口
func (cc *CalendarConfig) creatMonthInfo(bd *BaseDate, days int) []BaseDate {
	var monthInfo []BaseDate
	for i := 0; i < days; i++ {
		bdItem := BaseDate{
			Year:  bd.Year,
			Month: bd.Month,
			Day:   i,
		}
		monthInfo = append(monthInfo, bdItem)
	}
	return monthInfo
}

//根据距离正月初一的天数计算农历日期
func (cc *CalendarConfig) getLunarByBetween(bd *BaseDate) BaseDate {
	yearDetail := cc.LunarInfo[bd.Year-cc.MinYear]
	var bdSource = BaseDate{
		Year:  bd.Year,
		Month: time.Month(yearDetail[1] - 1),
		Day:   yearDetail[2],
	}
	days := cc.getDaysBetweenSolar(bd, &bdSource)
	var bdResult BaseDate
	if days == 0 { //正月初一
		bdResult.Year = bd.Year
		bdResult.Month = 0
		bdResult.Day = 1
	} else {
		var lunarYear int
		if days > 0 {
			lunarYear = bd.Year
		} else {
			lunarYear = bd.Year - 1
		}
		bdResult = cc.getLunarDateByBetween(lunarYear, days)
	}
	return bdResult
}

//通过间隔天数查找农历日期
func (cc *CalendarConfig) getLunarDateByBetween(lunarYear int, days int) BaseDate {
	yearDays, monthDays := cc.getLunarYearDays(lunarYear)
	var end, tempDays, month int
	if days > 0 {
		end = days
	} else {
		end = yearDays + days
	}
	for i := 0; i < len(monthDays); i++ {
		tempDays += monthDays[i]
		if tempDays > end {
			month = i
			tempDays = tempDays - monthDays[i]
			break
		}
	}
	bd := BaseDate{
		Year:  lunarYear,
		Month: time.Month(month),
		Day:   end - tempDays + 1,
	}
	return bd
}

//获取农历年份一年的每月的天数及一年的总天数
func (cc *CalendarConfig) getLunarYearDays(lunarYear int) (yearDays int, monthDays []int) {
	yearData := cc.LunarInfo[lunarYear-cc.MinYear]
	//闰月所在月 0为没有
	leapMonth := yearData[0]
	formatInt := strconv.FormatInt(int64(yearData[3]), 2)
	formatIntString := fmt.Sprintf("%016s", formatInt)
	monthDataArr := strings.Split(formatIntString, "")
	yearMonth := 12
	if leapMonth > 0 {
		yearMonth = 13
	}
	for i := 0; i < yearMonth; i++ {
		days := 30
		if monthDataArr[i] == "0" {
			days = 29
		}
		yearDays += days
		monthDays = append(monthDays, days)

	}
	return
}

//两个公历日期之间的天数
func (cc *CalendarConfig) getDaysBetweenSolar(bdTarget, bdSource *BaseDate) int {
	bdTarget.Month = time.Month(int(bdTarget.Month) + 1)
	bdSource.Month = time.Month(int(bdSource.Month) + 1)
	targetDate := time.Date(bdTarget.Year, bdTarget.Month, bdTarget.Day, 0, 0, 0, 0, time.Local)
	sourceDate := time.Date(bdSource.Year, bdSource.Month, bdSource.Day, 0, 0, 0, 0, time.Local)
	return int((targetDate.Unix() - sourceDate.Unix()) / 86400)
}

//将公历转换为农历
//month 从0开始 显示需要 +1
func (cc *CalendarConfig) SolarToLunar(bd *BaseDate) (BaseDate, error) {
	bd, err := cc.formatDate(bd)
	if err != nil {
		return BaseDate{}, errors.New("date format error: " + err.Error())
	}
	bdResult := cc.getLunarByBetween(bd)
	return bdResult, nil
}

//公历某月日历
//fill 是否用上下月数据补齐首尾空缺，首例数据从周日开始
func (cc *CalendarConfig) SolarCalendar(bd *BaseDate) (MonthInfo, error) {
	bd, err := cc.formatDate(bd)
	if err != nil {
		return MonthInfo{}, errors.New("date format error: " + err.Error())
	}
	date := time.Date(bd.Year, bd.Month, 1, 0, 0, 0, 0, time.Local)
	nextDate := date.AddDate(0, 1, 0)
	var monthInfo MonthInfo
	//该月一号星期几
	monthInfo.WeekOf1st = int(date.Weekday())
	monthInfo.MonthDays = int((nextDate.Unix() - date.Unix()) / 86400)
	monthInfo.DayInfo = cc.creatMonthInfo(bd, monthInfo.MonthDays)
	return monthInfo, nil
}
