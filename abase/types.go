package abase

import (
	"fmt"
	"time"
)

type BaseDate struct {
	Year,
	Day int
	Month time.Month
}

func (b *BaseDate) StringMonthDay() string {
	return fmt.Sprintf("%02d-%02d", b.Month+1, b.Day)
}
func (b *BaseDate) String() string {
	return fmt.Sprintf("%d-%02d-%02d", b.Year, b.Month+1, b.Day)
}

type MonthInfo struct {
	WeekOf1st int //该月1号星期几
	MonthDays int //该月天数
	DayInfo   []BaseDate
}

type Person struct {
	TypeInt     int    `json:"TypeInt"`
	Name        string `json:"Name"`
	Date        string `json:"Date"`
	ReceiveMail string `json:"ReceiveMail"`
}

type SendMailConfig struct {
	Person Person
	ToMail string
}
