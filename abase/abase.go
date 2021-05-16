package abase

import (
	"errors"
	"fmt"
	"github.com/imroc/biu"
	"github.com/patrickmn/go-cache"
	"strings"
	"time"
)

type CalendarConfig struct {
	cache                                                              *cache.Cache
	MinYear, MaxYear                                                   int
	HeavenlyStems, EarthlyBranches, Zodiac, SolarTerm, MonthCn, DateCn []string
	LunarFestival                                                      map[string]string
	LunarInfo                                                          [][4]int
}

func Init() CalendarConfig {
	var cc CalendarConfig
	//实例化缓存
	cc.cache = cache.New(10*time.Minute, 20*time.Minute)

	cc.MinYear = 1890
	cc.MaxYear = 2100
	heavenlyStems := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"} //天干
	cc.HeavenlyStems = heavenlyStems
	earthlyBranches := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"} //地支
	cc.EarthlyBranches = earthlyBranches
	Zodiac := []string{"鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"} //对应地支十二生肖
	cc.Zodiac = Zodiac
	//二十四节气
	SolarTerm := []string{"小寒", "大寒", "立春", "雨水", "惊蛰", "春分", "清明", "谷雨", "立夏", "小满", "芒种", "夏至", "小暑", "大暑", "立秋", "处暑", "白露", "秋分", "寒露", "霜降", "立冬", "小雪", "大雪", "冬至"}
	cc.SolarTerm = SolarTerm
	MonthCn := []string{"正", "二", "三", "四", "五", "六", "七", "八", "九", "十", "十一", "十二"}
	cc.MonthCn = MonthCn
	DateCn := []string{"初一", "初二", "初三", "初四", "初五", "初六", "初七", "初八", "初九", "初十", "十一", "十二", "十三", "十四", "十五", "十六", "十七", "十八", "十九", "二十", "廿一", "廿二", "廿三", "廿四", "廿五", "廿六", "廿七", "廿八", "廿九", "三十", "卅一"}
	cc.DateCn = DateCn
	//农历节日
	lunarFestival := map[string]string{
		"d0101": "春节",
		"d0115": "元宵节",
		"d0202": "龙抬头节",
		"d0323": "妈祖生辰",
		"d0505": "端午节",
		"d0707": "七夕情人节",
		"d0715": "中元节",
		"d0815": "中秋节",
		"d0909": "重阳节",
		"d1015": "下元节",
		"d1208": "腊八节",
		"d1223": "小年",
		"d0100": "除夕"}
	cc.LunarFestival = lunarFestival
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
	yearDays, monthDays := cc.getLunarYearDays(lunarYear, days)
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
func (cc *CalendarConfig) getLunarYearDays(lunarYear int, days int) (yearDays int, monthDays []int) {
	yearData := cc.LunarInfo[lunarYear-cc.MinYear]
	//闰月所在月 0为没有
	leapMonth := yearData[0]
	monthData := biu.ToBinaryString(yearData[3]) //转为二进制
	monthData = monthData[1 : len(monthData)-1]
	bytes := strings.Split(monthData, " ")
	monthData = bytes[len(bytes)-2] + bytes[len(bytes)-1]
	monthDataArr := strings.Split(monthData, "")
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
func (cc *CalendarConfig) SolarToLunar(bd *BaseDate) BaseDate {
	bd, err := cc.formatDate(bd)
	if err != nil {
		panic(err)
	}
	bdResult := cc.getLunarByBetween(bd)
	return bdResult
}

//公历某月日历
//fill 是否用上下月数据补齐首尾空缺，首例数据从周日开始
func (cc *CalendarConfig) SolarCalendar(bd *BaseDate, fill bool) MonthInfo {
	bd, err := cc.formatDate(bd)
	if err != nil {
		panic(err)
	}
	date := time.Date(bd.Year, bd.Month, 1, 0, 0, 0, 0, time.Local)
	nextDate := date.AddDate(0, 1, 0)
	var monthInfo MonthInfo
	//该月一号星期几
	monthInfo.WeekOf1st = int(date.Weekday())
	monthInfo.MonthDays = int((nextDate.Unix() - date.Unix()) / 86400)
	monthInfo.DayInfo = cc.creatMonthInfo(bd, monthInfo.MonthDays)
	fmt.Println(date)
	fmt.Println(nextDate)
	fmt.Println(monthInfo)
	return monthInfo

}

//func (cc *CalendarConfig) Calendar (year int, month int, fill bool)  {
//	fmt.Println(cc.MinYear)
//}
//func (cc *CalendarConfig) Formatter (year int, month int, fill bool) map[string]int  {
//	return map[string]int{
//		"a": 1,
//	}
//}
