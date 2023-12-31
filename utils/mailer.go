package utils

import (
	"log"
	"net/smtp"

	"github.com/mhafidk/ngartos/config"
)

func SendVerificationEmail(verificationToken, email string) {
	from := config.Config("MAIL_FROM")
	pass := config.Config("MAIL_PASS")
	to := email

	var verifUrl string
	if config.Config("ENVIRONMENT") == "DEV" {
		verifUrl = config.Config("VERIF_URL_DEV")
	} else {
		verifUrl = config.Config("VERIF_URL")
	}

	body := "Terima kasih telah bergabung dengan Ngartos!\nLangkah selanjutnya yaitu verifikasi email, klik link berikut untuk dapat memverifikasi email Anda.\n\n" + verifUrl + verificationToken
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Verifikasi Email\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, pass, "smtp.gmail.com"), from, []string{to}, []byte(msg))
	if err != nil {
		log.Println(err)
	}
}

func SendForgotPasswordEmail(forgotPasswordToken, email string) {
	from := config.Config("MAIL_FROM")
	pass := config.Config("MAIL_PASS")
	to := email

	var forgotPasswordUrl string
	if config.Config("ENVIRONMENT") == "DEV" {
		forgotPasswordUrl = config.Config("FORGOT_PASSWORD_URL_DEV")
	} else {
		forgotPasswordUrl = config.Config("FORGOT_PASSWORD_URL")
	}

	body := "Klik link di bawah ini untuk mereset password Anda.\n\n" + forgotPasswordUrl + forgotPasswordToken
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Reset Password\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, pass, "smtp.gmail.com"), from, []string{to}, []byte(msg))
	if err != nil {
		log.Println(err)
	}
}
