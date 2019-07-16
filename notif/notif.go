package notif

import (
	email "github.com/Varunram/essentials/email"
)

// footerString is a common footer string that is used by all emails
var footerString = "Have a nice day!\n\nWarm Regards, \nThe OpenClimate Team\n\n\n\n" +
	"You're receiving this email because your contact was given" +
	" on the OpenClimate platform for receiving notifications on reports in which you're a party.\n\n\n"

// SendReminder sends a reminder to platform participants to submit their data
func SendReminder(to string) error {
	// this is sent to the recipient on investment from an investor
	body := "Greetings from the OpenClimate platform! \n\n" +
		"This is a reminder to let you know that your reports for the next month are due\n\n" +
		"Please submit them within the submission window in order to enable analysis on the platform dashboard\n" +
		footerString
	return email.SendMail(body, to)
}
