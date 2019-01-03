package players

import (
	//	"errors"
	"encoding/json"
	"os"

	"github.com/ckpt/backend-services/utils"
	mailgun "github.com/mailgun/mailgun-go"
	//	"github.com/m4rw3r/uuid"
)

import "fmt"

func StartEventProcessor() error {
	events, err := eventqueue.Consume()
	if err != nil {
		fmt.Printf("Could not consume events:\nError was:\n%v\n", err)
		return err
	}
	go func() {
		fmt.Printf("Entering main event processor loop\n")
		for msg := range events {

			// Fetch event
			var event utils.CKPTEvent
			_ = json.Unmarshal(msg.Body, &event)
			fmt.Printf("Found new event of type: %d\n", event.Type)

			// Notify subscribers
			allPlayers, err := AllPlayers()
			if err != nil {
				msg.Nack(false, true)
				continue
			}
			for _, p := range allPlayers {
				notifyPlayer := false
				fmt.Printf("Checking if user %s should be notified\n", p.Nick)
				if p.User.SubscribedTo(utils.TypeNames[event.Type]) {
					notifyPlayer = true
				}
				for _, rp := range event.RestrictedTo {
					notifyPlayer = false
					if rp == p.UUID {
						notifyPlayer = true
						break
					}
				}
				if notifyPlayer {
					fmt.Printf("Notifying user for event\n")
					NotifyUser(p.Nick, p.Profile.Email, event.Subject, event.Message)
				}
			}
			msg.Ack(false)
		}
	}()

	return nil
}

func NotifyUser(name, email, subject, message string) {
	fmt.Printf("Notifying %s with subject:\n", email)
	fmt.Printf("%s\n", subject)

	mailto := fmt.Sprintf("%s <%s>", name, email)
	gun := mailgun.NewMailgun("mail.ckpt.no", os.Getenv("CKPT_MAILGUN_KEY"), "pubkey-b3e133632123a0da24d1e2c5842039b6")
	m := mailgun.NewMessage("CKPT <notifications@mail.ckpt.no>", subject, message, mailto)
	m.AddHeader("Content-Type", "text/plain; charset=\"utf-8\"")
	response, id, err := gun.Send(m)

	if err != nil {
		fmt.Printf("Error:\n%+v\n", err.Error())
	}
	fmt.Printf("Response ID: %s\n", id)
	fmt.Printf("Message from server: %s\n", response)
}
