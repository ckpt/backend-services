package utils

import (
    "encoding/json"
    "github.com/streadway/amqp"
)

// An interface for a message queue
type AMQPQueue interface {
    Publish(CKPTEvent) error
    Consume() (<-chan amqp.Delivery, error)
}

type EventType int

const (
    NEWS_EVENT EventType = iota
    TOURNAMENT_EVENT
    CATERING_EVENT
    LOCATION_EVENT
    PLAYER_EVENT = 4
)

var TypeNames = []string{
    "news",
    "tournament",
    "catering",
    "location",
    "player",
}

type CKPTEvent struct {
    Type EventType `json:"type"`
    Subject string `json:"subject"`
    Message string `json:"message"`
}

type RMQ struct {
    url string
    queue string
}

func (rmq *RMQ) Publish(event CKPTEvent) error {
    conn, pub, err := rmq.setup()
    defer conn.Close()
    if err != nil {
        return err
    }
    msg, err := json.Marshal(event)
    if err != nil {
        return err
    }
    if err := pub.Publish("", rmq.queue, false, false, amqp.Publishing{
        ContentType: "text/json",
        ContentEncoding: "utf-8",
        Body: msg,
    }); err != nil {
        return err
    }
    return nil
}

func (rmq *RMQ) Consume() (<-chan amqp.Delivery, error) {
    _, ch, err := rmq.setup()
    if err != nil {
        return nil, err
    }
    
    deliveries, err := ch.Consume(rmq.queue, "", false, false, false, false, nil)
    return deliveries, err
}

func (rmq *RMQ) setup() (*amqp.Connection, *amqp.Channel, error) {
    conn, err := amqp.Dial(rmq.url)
    if err != nil {
        return nil, nil, err
    }
    
    ch, err := conn.Channel()
    if err != nil {
        return nil, nil, err
    }
    
    if _, err := ch.QueueDeclare(rmq.queue, false, false, false, false, nil); err != nil {
        return nil, nil, err
    }
    
    return conn, ch, nil
}

func NewRMQ(url string, queue string) *RMQ {
    rmq := new(RMQ)
    rmq.url = url
    rmq.queue = queue
    return rmq
}