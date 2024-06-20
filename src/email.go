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

func get_email_message(username string, problems []map[string]string) *mail.Msg {
	//call database to get the problems for a specific user
	email_message := mail.NewMsg()
	email_message.FromFormat("LeetCode Spaced Repitition", os.Getenv("EMAIL_ADDRESS"))
	email_message.Subject("Testing Bulk Email!")

	mystr := "Hello " + username + ", here are the problems you should complete today:\n"
	for _, problem := range problems {
		mystr += (problem["test"] + "\n")
	}

	email_message.SetBodyString(mail.TypeTextPlain, mystr)
	email_message.SetBulk()
	return email_message
}

func send_email(row Row) {
	client := create_email_client()

	email_message := get_email_message(row.Username, row.Problems)

	if err := email_message.To(string(row.Email)); err != nil {
		fmt.Printf("failed to set To address: %s", err)
		// delete this email from list (aka not a valid email)
	}

	if err := client.DialAndSend(email_message); err != nil {
		fmt.Printf("failed to send mail: %s", err)
		return
	}

	fmt.Println("success")
}
