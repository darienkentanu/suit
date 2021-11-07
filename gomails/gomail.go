package gomails

import (
	"fmt"
	"log"

	"github.com/darienkentanu/suit/constants"
	"gopkg.in/gomail.v2"
)

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_SENDER_NAME = "PT. UANG JADI DUIT (SUIT) <rendaysjr@gmail.com>"
const CONFIG_AUTH_EMAIL = constants.EMAIL_ADDRESS
const CONFIG_AUTH_PASSWORD = constants.EMAIL_PASSWORD

func SendMail(recipient string, message string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", recipient, "only.adicipta@gmail.com")
	mailer.SetAddressHeader("Cc", "rizkakhairanii@gmail.com", "Tra Lala La")
	mailer.SetHeader("Subject", "Test mail")
	isi := fmt.Sprintf("Hello, <b>%v</b>", message)
	mailer.SetBody("text/html", isi)
	// mailer.Attach("./sample.png")

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	log.Println("Mail sent!")
	return nil
}
