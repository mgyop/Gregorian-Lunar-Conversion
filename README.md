# Gregorian-Lunar-Conversion

阳历转阴历
####安装
````bigquery
go get github.com/mgyop/Gregorian-Lunar-Conversion
````
####示例
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