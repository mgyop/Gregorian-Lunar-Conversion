# Gregorian-Lunar-Conversion

Gregorian-Lunar Conversion

````
func main() {
    ##初始化配置
	calendarConfig := abase.Init()
	##构建当前时间
	now := time.Now()
	bd := abase.BaseDate{
		Year:  now.Year(),
		Month: now.Month(),
		Day:   now.Day(),
	}
	##获取农历日期
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(LunarDate.String())
}
````