package main

import (
	"Gregorian-Lunar-Conversion/abase"
	"fmt"
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
	LunarDate, err := calendarConfig.SolarToLunar(&bd)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(LunarDate.String())
}
