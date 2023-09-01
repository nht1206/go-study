package sse

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nht1206/go-study/subscriptions"
)

type NotificationHandler struct {
	broker *subscriptions.Broker
}

func NewNotificationHandler(broker *subscriptions.Broker) *NotificationHandler {
	return &NotificationHandler{
		broker: broker,
	}
}

func (h *NotificationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	client := subscriptions.NewDefaultClient()
	h.broker.Register(client)
	defer h.broker.Unregister(client.Id())

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	w.Write([]byte("id: " + client.Id() + "\n"))
	w.Write([]byte("event: Connect\n"))
	w.Write([]byte("data: "))
	w.Write([]byte(fmt.Sprintf(`{ "clientId": %s }`, client.Id())))
	w.Write([]byte("\n\n"))

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	log.Println("clientId: ", client.Id())

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	r = r.WithContext(ctx)

	idleTimer := time.NewTimer(5 * time.Minute)
	defer idleTimer.Stop()
	for {
		select {
		case <-idleTimer.C:
			cancel()
		case msg, ok := <-client.Channel():

			if !ok {
				continue
			}

			w.Write([]byte("id: " + client.Id() + "\n"))
			w.Write([]byte("event: " + msg.Name + "\n"))
			w.Write([]byte("data: "))
			w.Write(msg.Data)
			w.Write([]byte("\n\n"))

			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-r.Context().Done():
			return
		}
	}
}

type SubscriptionHandler struct {
	broker *subscriptions.Broker
}

func NewSubscriptionHandler(broker *subscriptions.Broker) *SubscriptionHandler {
	return &SubscriptionHandler{
		broker: broker,
	}
}

func (h *SubscriptionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clientId := r.FormValue("clientId")
	subsription := r.FormValue("subscription")
	client, err := h.broker.ClientById(clientId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Printf("Set subsription %q to %s", subsription, clientId)

	client.Subscribe(subsription)
}
