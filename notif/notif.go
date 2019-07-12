package notif

import (
	"github.com/pkg/errors"
	"log"
	"net/smtp"

	"github.com/spf13/viper"
)

// footerString is a common footer string that is used by all emails
var footerString = "Have a nice day!\n\nWarm Regards, \nThe Openclimate Team\n\n\n\n" +
	"You're receiving this email because your contact was given" +
	" on the openclimate platform for receiving notifications on reports in which you're a party.\n\n\n"

// sendMail is a handler for sending out an email to an entity, reading required params
// from the config file
func sendMail(body string, to string) error {
	var err error
	// read from config.yaml in the working directory
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "error while reading email values from config file")
	}
	log.Println("SENDING EMAIL: ", viper.Get("email"), viper.Get("password"))
	from := viper.Get("email").(string)    // interface to string
	pass := viper.Get("password").(string) // interface to string
	auth := smtp.PlainAuth("", from, pass, "smtp.gmail.com")
	// to can also be an array of addresses if needed
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Openclimate Notification\n\n" + body

	err = smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(msg))
	if err != nil {
		return errors.Wrap(err, "smtp error")
	}
	return nil
}

// SendInvestmentNotifToRecipient sends a notification to the recipient when an investor
// invests in an order he's the recipient of
func SendReminder(to string) error {
	// this is sent to the recipient on investment from an investor
	body := "Greetings from the openclimate platform! \n\n" +
		"This is a reminder to let you know that your reports for the next month are due\n\n" +
		"Please submit them within the submission window in order to enable analysis on the platform dashboard\n" +
		footerString
	return sendMail(body, to)
}
