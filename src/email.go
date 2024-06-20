package main

import (
	"fmt"
	"os"

	"github.com/wneessen/go-mail"
)

func create_email_client() *mail.Client {
	client, _ := mail.NewClient(
		"smtp.gmail.com",
		mail.WithPort(587),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(os.Getenv("EMAIL_ADDRESS")),
		mail.WithPassword(os.Getenv("EMAIL_APP_PASSWORD")),
	)
	return client
}

func get_email_message(problems []map[string]string) *mail.Msg {
	//call database to get the problems for a specific user
	email_message := mail.NewMsg()
	email_message.FromFormat("LeetCode Spaced Repitition", os.Getenv("EMAIL_ADDRESS"))
	email_message.Subject("Testing Bulk Email!")

	mystr := "Here are the problems you should complete today:\n"
	for _, problem := range problems {
		mystr += (problem["test"] + "\n")
	}

	email_message.SetBodyString(mail.TypeTextPlain, mystr)

	email_message.SetBulk()
	return email_message
}

func send_email(email_recipients []string, problems []map[string]string) {
	client := create_email_client()

	var email_messages []*mail.Msg

	for _, email := range email_recipients {
		email_message := get_email_message(problems)

		if err := email_message.To(string(email)); err != nil {
			fmt.Printf("failed to set To address: %s", err)
			// delete this email from list (aka not a valid email)
			continue
		}

		email_messages = append(email_messages, email_message)
	}

	if err := client.DialAndSend(email_messages...); err != nil {
		fmt.Printf("failed to send mail: %s", err)
	}

	fmt.Println("success")
}
