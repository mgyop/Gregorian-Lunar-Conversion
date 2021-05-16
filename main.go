package main

import (
	"calendar/abase"
	"calendar/mail"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	for {
		taskRunning()
		time.Sleep(24 * time.Hour)
	}
}

func taskRunning() {
	calendarConfig := abase.Init()
	now := time.Now()
	bd := abase.BaseDate{
		Year:  now.Year(),
		Month: now.Month(),
		Day:   now.Day(),
	}
	LunarDate := calendarConfig.SolarToLunar(&bd)
	//读取配置文件
	persons := readConfig()
	var sendTask []abase.SendMailConfig
	for _, config := range persons {
		toMails := strings.Split(config.ReceiveMail, ",")
		var compaireDate string
		if config.TypeInt == 2 { //公历生日
			compaireDate = bd.StringMonthDay()
		} else {
			compaireDate = LunarDate.StringMonthDay()
		}
		if compaireDate == config.Date {
			for _, mail := range toMails {
				if len(mail) > 0 {
					var configItem = abase.SendMailConfig{
						Person: config,
						ToMail: mail,
					}
					sendTask = append(sendTask, configItem)
				}

			}
		}
	}
	if len(sendTask) > 0 {
		var wg sync.WaitGroup
		ch := make(chan abase.SendMailConfig, 10)
		for _, configItem := range sendTask {
			ch <- configItem
		}
		wg.Add(len(sendTask))
		for i := 0; i < len(sendTask); i++ {
			go func() {
				for receiveConfig := range ch {
					mail.SendMail(receiveConfig, &wg)
				}
			}()
		}
		wg.Wait()
		close(ch)
	} else {
		//今天没有谁过生日
		mail.Logging("./runtime/none_log_"+time.Now().Format("2006-01-02"), []byte("今天没有发出去邮件"))
	}

}

func readConfig() []abase.Person {
	file, err := os.Open("./config/config.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err != nil {
		fmt.Println("配置文件错误")
	}
	var persons []abase.Person
	err = decoder.Decode(&persons)
	if err != nil {
		fmt.Println("配置文件错误")
	}
	return persons
}
