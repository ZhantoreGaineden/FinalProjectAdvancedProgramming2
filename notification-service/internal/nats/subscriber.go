package natslistener

import (
	"encoding/json"
	"log"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/notification-service/internal/email"
	"github.com/nats-io/nats.go"
)

type Subscriber struct {
	conn   *nats.Conn
	sender *email.Sender
}

type UserRegisteredEvent struct {
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type OrderCreatedEvent struct {
	OrderID    string  `json:"order_id"`
	UserID     string  `json:"user_id"`
	UserEmail  string  `json:"user_email"`
	TotalPrice float64 `json:"total_price"`
}

type OrderStatusUpdatedEvent struct {
	OrderID   string `json:"order_id"`
	UserEmail string `json:"user_email"`
	Status    string `json:"status"`
}

func NewSubscriber(conn *nats.Conn, sender *email.Sender) *Subscriber {
	return &Subscriber{
		conn:   conn,
		sender: sender,
	}
}

func (s *Subscriber) Subscribe() error {
	if _, err := s.conn.Subscribe("user.registered", s.handleUserRegistered); err != nil {
		return err
	}

	if _, err := s.conn.Subscribe("order.created", s.handleOrderCreated); err != nil {
		return err
	}

	if _, err := s.conn.Subscribe("order.status_updated", s.handleOrderStatusUpdated); err != nil {
		return err
	}

	log.Println("Notification Service subscribed to NATS subjects")
	return nil
}

func (s *Subscriber) handleUserRegistered(msg *nats.Msg) {
	var event UserRegisteredEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("failed to parse user.registered event: %v", err)
		return
	}

	subject := "Welcome to Pet Store"
	body := "Hello " + event.FullName + ", welcome to our Pet Store!"

	if err := s.sender.Send(event.Email, subject, body); err != nil {
		log.Printf("failed to send welcome email: %v", err)
		return
	}

	log.Printf("welcome email sent to %s", event.Email)
}

func (s *Subscriber) handleOrderCreated(msg *nats.Msg) {
	var event OrderCreatedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("failed to parse order.created event: %v", err)
		return
	}

	subject := "Order Created"
	body := "Your order " + event.OrderID + " has been created successfully."

	if err := s.sender.Send(event.UserEmail, subject, body); err != nil {
		log.Printf("failed to send order created email: %v", err)
		return
	}

	log.Printf("order created email sent to %s", event.UserEmail)
}

func (s *Subscriber) handleOrderStatusUpdated(msg *nats.Msg) {
	var event OrderStatusUpdatedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("failed to parse order.status_updated event: %v", err)
		return
	}

	subject := "Order Status Updated"
	body := "Your order " + event.OrderID + " status changed to: " + event.Status

	if err := s.sender.Send(event.UserEmail, subject, body); err != nil {
		log.Printf("failed to send order status email: %v", err)
		return
	}

	log.Printf("order status email sent to %s", event.UserEmail)
}
