package gomails

import (
	"fmt"
	"log"

	"github.com/darienkentanu/suit/constants"
	"gopkg.in/gomail.v2"
)

func SendMail(recipient string, message string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", constants.CONFIG_SENDER_NAME)
	mailer.SetHeader("To", recipient, "only.adicipta@gmail.com")
	mailer.SetAddressHeader("Cc", "rizkakhairanii@gmail.com", "Tra Lala La")
	mailer.SetHeader("Subject", "Test mail")
	isi := fmt.Sprintf("Hello, <b>%v</b>", message)
	mailer.SetBody("text/html", isi)
	// mailer.Attach("./sample.png")

	dialer := gomail.NewDialer(
		constants.CONFIG_SMTP_HOST,
		constants.CONFIG_SMTP_PORT,
		constants.CONFIG_AUTH_EMAIL,
		constants.CONFIG_AUTH_PASSWORD,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	log.Println("Mail sent!")
	return nil
}
