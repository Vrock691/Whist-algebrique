package main

import (
	"log"
	"net/http"

	event "fr.vamary.whist-server/internal/event"
	eventhandler "fr.vamary.whist-server/internal/event/handler"
	api "fr.vamary.whist-server/internal/server"
	"github.com/google/uuid"
)

func main() {
	handler := newHTTPHandler()

	log.Println("REST API and WebSocket server listening on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}

func newHTTPHandler() http.Handler {
	var eventHandler *eventhandler.Handler
	webSocketServer := event.NewWebSocketServer(func(gameID uuid.UUID, envelope event.Envelope) error {
		return eventHandler.HandleMessage(gameID, envelope)
	})
	eventHandler = eventhandler.New(webSocketServer, eventhandler.Callbacks{})

	mux := http.NewServeMux()
	mux.HandleFunc("GET /ws", webSocketServer.Handle)
	return api.HandlerFromMux(api.NewServer(), mux)
}
