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

func get_welcome_email_message(username string) *mail.Msg {
	email_message := mail.NewMsg()
	email_message.SetBulk()
	email_message.FromFormat("LeetCode Spaced Repetition", os.Getenv("EMAIL_ADDRESS"))
	email_message.Subject("Thank you for subscribing to LeetCode Spaced Repetition!")
	email_message.SetBodyString(mail.TypeTextPlain, `Hello `+username+`, 
	the developer (hyperlink to my socials) appreciates you for subscribing to LeetCode Spaced Repetition.
	When it is time to review one or more LeetCode questions you will automatically be sent an email
	reminding you of which problems should be completed on that particular day.
	
	So what is spaced repetition anyways?
	Spaced repetition is an evidence-based learning technique that is usually performed with flashcards, but in
	this case LeetCode problems. Newly introduced and more difficult flashcards are shown more frequently, while 
	older and less difficult flashcards are shown less frequently in order to exploit the psychological spacing effect. 
	The use of spaced repetition has been proven to increase the rate of learning.`)
	return email_message
}

func get_daily_email_message(username string, problems []LeetCodeProblem) *mail.Msg {
	//call database to get the problems for a specific user
	email_message := mail.NewMsg()
	email_message.SetBulk()
	email_message.FromFormat("LeetCode Spaced Repetition", os.Getenv("EMAIL_ADDRESS"))
	email_message.Subject("Testing Bulk Email!")

	mystr := "Hello " + username + ", here are the problems you should complete today:\n"
	for _, problem := range problems {
		mystr += (problem.Link + "\n")
	}

	email_message.SetBodyString(mail.TypeTextPlain, mystr)
	return email_message
}

func send_daily_email(subscriber Subscriber) {
	client := create_email_client()

	email_message := get_daily_email_message(subscriber.Username, subscriber.Problems)

	if err := email_message.To(string(subscriber.Email)); err != nil {
		fmt.Printf("failed to set To address: %s", err)
		// delete this email from list (aka not a valid email)
	}

	if err := client.DialAndSend(email_message); err != nil {
		fmt.Printf("failed to send mail: %s", err)
		return
	}

	fmt.Println("success")
}

func send_welcome_email(username string) {
	client := create_email_client()
	email_message := get_welcome_email_message(username)

	if err := client.DialAndSend(email_message); err != nil {
		fmt.Printf("failed to send mail: %s", err)
		return
	}
	fmt.Println("success")
}
