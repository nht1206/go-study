package main

import (
	"log"
	"net/http"
	"time"

	"github.com/nht1206/go-study/sse"
	"github.com/nht1206/go-study/subscriptions"
	"goji.io/v3"
	"goji.io/v3/pat"
)

func main() {

	broker := subscriptions.NewBroker()

	notificationHander := sse.NewNotificationHandler(broker)
	subscriptionHandler := sse.NewSubscriptionHandler(broker)

	mux := goji.NewMux()

	mux.Handle(pat.Get("/notification"), notificationHander)
	mux.Handle(pat.Post("/notification"), subscriptionHandler)

	go sendNotification(broker)

	http.ListenAndServe(":8080", mux)
}

func sendNotification(broker *subscriptions.Broker) {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		for _, client := range broker.Clients() {
			if client.HasSubscription("Notification") {
				log.Println("Send notification to ", client.Id())
				client.Send(subscriptions.Message{
					Name: "Notification",
					Data: []byte("This is a demo nofitication"),
				})
			}
		}
	}
}
