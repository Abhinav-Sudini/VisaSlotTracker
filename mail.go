package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

func SendMail(To string, avl_slots int, loc string) {
	from := "sudiniabhinav@gmail.com"
	appPassword := os.Getenv("EMAIL_PASS")

	// to := []string{"abhinavsai5205418@gmail.com"}
	to := []string{To}

	msg := []byte(
		fmt.Sprintf("Subject: Visa slots available at %v \r\n", loc) +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\";\r\n\r\n" +

		"Hi,\n\n" +

		"update regarding visa appointment availability.\n\n" +
		fmt.Sprintf(
			"Currently, %d appointment slots are available at the location := %s \n",
			avl_slots,
			loc,
		) +
		"\nit might be a good time to check the portal.\n\n" +

		"Best of luck\n" +
		"abhinav sudini\n",
	)

	auth := smtp.PlainAuth("", from, appPassword, "smtp.gmail.com")

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		from,
		to,
		msg,
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Email sent successfully to - ", to)
}
