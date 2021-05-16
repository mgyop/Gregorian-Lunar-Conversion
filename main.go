package main

import (
	"Gregorian-Lunar-Conversion/abase"
	"time"
)

func main() {
	calendarConfig := abase.Init()
	now := time.Now()
	bd := abase.BaseDate{
		Year:  now.Year(),
		Month: now.Month(),
		Day:   now.Day(),
	}
	LunarDate := calendarConfig.SolarToLunar(&bd)
	println(LunarDate.String())
}
