package main

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
)

type AlarmClient interface {
	StartAlarm() error
	StopAlarm() error
}

func ConsumeMessages(sub message.Subscriber, alarmClient AlarmClient) {
	messages, err := sub.Subscribe(context.Background(), "smoke_sensor")
	if err != nil {
		panic(err)
	}

	for msg := range messages {
		switch string(msg.Payload) {
		case "1":
			handleSmokeDetected(msg, alarmClient)
			continue

		case "0":
			handleSmokeCleared(msg, alarmClient)
			continue

		}
	}
}

func handleSmokeDetected(msg *message.Message, alarmClient AlarmClient) {
	err := alarmClient.StartAlarm()
	if err != nil {
		msg.Nack()
		return
	}

	msg.Ack()
}

func handleSmokeCleared(msg *message.Message, alarmClient AlarmClient) {
	err := alarmClient.StopAlarm()
	if err != nil {
		msg.Nack()
		return
	}

	msg.Ack()
}
