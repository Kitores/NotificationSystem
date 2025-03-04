package email_sender

import (
	"fmt"
	"net/smtp"
)

func SendMailFunc(to []string, subject, content, mailFrom, mailPass, mailHost, mailPort string) error {
	message := []byte("Subject:" + subject + "\r\n" +
		"From: " + mailFrom + "\r\n\r\n" +
		content)
	auth := smtp.PlainAuth("", mailFrom, mailPass, mailHost)

	err := smtp.SendMail(mailHost+":"+mailPort, auth, mailFrom, to, message)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("Email sent successfully!")
	return err
}
