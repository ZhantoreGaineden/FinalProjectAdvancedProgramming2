package main

import (
	"log"
	"time"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/notification-service/internal/config"
	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/notification-service/internal/email"
	natslistener "github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/notification-service/internal/nats"
	"github.com/nats-io/nats.go"
)

func main() {
	cfg := config.Load()

	var natsConn *nats.Conn
	var err error

	for i := 1; i <= 10; i++ {
		natsConn, err = nats.Connect(cfg.NATSURL)
		if err == nil {
			break
		}

		log.Printf("failed to connect to nats, retry %d/10: %v", i, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("failed to connect to nats after retries: %v", err)
	}
	defer natsConn.Close()

	sender := email.NewSender(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
		cfg.FromEmail,
	)

	subscriber := natslistener.NewSubscriber(natsConn, sender)

	if err := subscriber.Subscribe(); err != nil {
		log.Fatalf("failed to subscribe to nats: %v", err)
	}

	log.Println("Notification Service started")

	select {}
}
