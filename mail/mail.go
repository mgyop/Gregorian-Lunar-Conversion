package mail

import (
	"calendar/abase"
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
	"os"
	"sync"
	"time"
)

func SendMail(config abase.SendMailConfig, wg *sync.WaitGroup) {
	defer wg.Done() //邮件发送出去后关闭同步计数器
	content := "阳历: " + time.Now().Format("2006-01-02") + ", 阴历: " + config.Person.Date + "\n" +
		"是" + config.Person.Name + "的生日, 不要忘记了奥"
	e := email.NewEmail()
	//设置发送方的邮箱
	e.From = "桃 <969060233@qq.com>"
	// 设置接收方的邮箱
	e.To = []string{config.ToMail}
	//设置主题
	e.Subject = "生日邮件提醒"
	//设置文件发送的内容
	e.Text = []byte(content)
	//设置服务器相关的配置
	err := e.Send("smtp.qq.com:25", smtp.PlainAuth("", "969060233@qq.com", "arpjkkcgmvsrbbif", "smtp.qq.com"))
	logContents := time.Now().String() + "\n"
	logContents += content + "\n"
	fileName := "./runtime/"
	if err != nil {
		fileName += "err_log.log"
	} else {
		fileName += "suc_log.log"
	}
	Logging(fileName, []byte(logContents))
}

func Logging(fileName string, contents []byte) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		log.Println("log error: ", err)
	}
	_, err = file.Write(contents)
	if err != nil {
		log.Println(err)
		//log.Println(write)
	}
	defer file.Close()

}
