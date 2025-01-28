package email_sender

import (
	"fmt"
	"net/smtp"
)

const (
	smtpHost = "smtp.mail.ru"
	smtpPort = "587"
)

func main() {
	var arr = []string{"mihail_yermolayev@mail.ru"}
	SendMailFunc("", arr, "", "", "")

}

func SendMailFunc(from string, to []string, subject string, content string, password string) error {
	message := []byte("Subject: Test\r\n" +
		"From: " + from + "\r\n" +
		"To: mihail_yermolayev@mail.ru\r\n\r\n" +
		content)
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Отправка почты
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("Email sent successfully!")
	return err
}
