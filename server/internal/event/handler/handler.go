// Package handler connects the generated event payloads to the WebSocket
// transport.
package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"fr.vamary.whist-server/internal/event"
	"github.com/google/uuid"
)

// Callbacks contains the application handlers for client-to-server events.
// Server-to-client events are sent with the Publish methods below.
type Callbacks struct {
	OnBetPlaced  func(gameID uuid.UUID, payload event.BetPlaced) error
	OnCardPlayed func(gameID uuid.UUID, payload event.CardPlayed) error
}

// Handler handles every channel declared by the AsyncAPI event contract.
type Handler struct {
	socket    *event.WebSocketServer
	callbacks Callbacks
}

func New(socket *event.WebSocketServer, callbacks Callbacks) *Handler {
	return &Handler{socket: socket, callbacks: callbacks}
}

// HandleMessage is suitable for event.NewWebSocketServer.
func (h *Handler) HandleMessage(gameID uuid.UUID, envelope event.Envelope) error {
	switch envelope.Channel {
	case "bet/placed":
		payload, err := decodePayload(gameID, envelope, func(value event.BetPlaced) uuid.UUID {
			return value.GameID
		})
		if err != nil {
			return err
		}
		if payload.CorrelationID == nil {
			payload.CorrelationID = envelope.CorrelationID
		}
		return h.HandleBetPlaced(gameID, payload)
	case "play/card":
		payload, err := decodePayload(gameID, envelope, func(value event.CardPlayed) uuid.UUID {
			return value.GameID
		})
		if err != nil {
			return err
		}
		if payload.CorrelationID == nil {
			payload.CorrelationID = envelope.CorrelationID
		}
		return h.HandleCardPlayed(gameID, payload)
	default:
		return fmt.Errorf("unsupported event channel %q", envelope.Channel)
	}
}

func (h *Handler) HandleBetPlaced(gameID uuid.UUID, payload event.BetPlaced) error {
	if payload.GameID != gameID {
		return errors.New("bet/placed payload gameId does not match the WebSocket room")
	}
	if h.callbacks.OnBetPlaced == nil {
		return errors.New("bet/placed handler is not configured")
	}
	return h.callbacks.OnBetPlaced(gameID, payload)
}

func (h *Handler) HandleCardPlayed(gameID uuid.UUID, payload event.CardPlayed) error {
	if payload.GameID != gameID {
		return errors.New("play/card payload gameId does not match the WebSocket room")
	}
	if h.callbacks.OnCardPlayed == nil {
		return errors.New("play/card handler is not configured")
	}
	return h.callbacks.OnCardPlayed(gameID, payload)
}

func decodePayload[T any](gameID uuid.UUID, envelope event.Envelope, getGameID func(T) uuid.UUID) (T, error) {
	var payload T
	if err := json.Unmarshal(envelope.Payload, &payload); err != nil {
		return payload, fmt.Errorf("invalid %s payload: %w", envelope.Channel, err)
	}
	payloadGameID := getGameID(payload)
	if payloadGameID == uuid.Nil {
		return payload, fmt.Errorf("%s payload must contain gameId", envelope.Channel)
	}
	if payloadGameID != gameID {
		return payload, fmt.Errorf("%s payload gameId does not match the WebSocket room", envelope.Channel)
	}
	return payload, nil
}

func (h *Handler) publish(channel string, payload any) error {
	if h.socket == nil {
		return errors.New("WebSocket server is not configured")
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal %s payload: %w", channel, err)
	}
	gameID := eventID(payload)
	if gameID == uuid.Nil {
		return fmt.Errorf("%s payload must contain gameId", channel)
	}
	envelope := event.Envelope{Channel: channel, Payload: data}
	return h.socket.Broadcast(gameID, envelope)
}

// eventID extracts the game ID shared by every event payload.
func eventID(payload any) uuid.UUID {
	value := reflect.ValueOf(payload)
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return uuid.Nil
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return uuid.Nil
	}
	gameID := value.FieldByName("GameID")
	if gameID.IsValid() && gameID.Type() == reflect.TypeOf(uuid.UUID{}) {
		return gameID.Interface().(uuid.UUID)
	}
	return uuid.Nil
}

func (h *Handler) PublishPlayerJoined(payload event.PlayerJoined) error {
	return h.publish("lobby/player/joined", payload)
}

func (h *Handler) PublishPlayerLeft(payload event.PlayerLeft) error {
	return h.publish("lobby/player/left", payload)
}

func (h *Handler) PublishPlayerDisconnected(payload event.PlayerDisconnected) error {
	return h.publish("lobby/player/disconnected", payload)
}

func (h *Handler) PublishGameStarted(payload event.GameStarted) error {
	return h.publish("game/started", payload)
}

func (h *Handler) PublishGameAborted(payload event.GameAborted) error {
	return h.publish("game/aborted", payload)
}

func (h *Handler) PublishGameFinished(payload event.GameFinished) error {
	return h.publish("game/finished", payload)
}

func (h *Handler) PublishTurnStarted(payload event.TurnStarted) error {
	return h.publish("turn/started", payload)
}

func (h *Handler) PublishTurnFinished(payload event.TurnFinished) error {
	return h.publish("turn/finished", payload)
}

func (h *Handler) PublishFoldStarted(payload event.FoldStarted) error {
	return h.publish("fold/started", payload)
}

func (h *Handler) PublishFoldFinished(payload event.FoldFinished) error {
	return h.publish("fold/finished", payload)
}

func (h *Handler) PublishBetAsk(payload event.BetAsk) error {
	return h.publish("bet/ask", payload)
}

func (h *Handler) PublishBetCompleted(payload event.BetCompleted) error {
	return h.publish("bet/completed", payload)
}

func (h *Handler) PublishPlayAsk(payload event.PlayAsk) error {
	return h.publish("play/ask", payload)
}
