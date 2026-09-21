// Package event contains the Go models for the WebSocket contract described
// in documentation/async_api.yaml.
package event

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Envelope is the common WebSocket frame format.
type Envelope struct {
	Channel       string          `json:"channel"`
	CorrelationID *uuid.UUID      `json:"correlationId,omitempty"`
	Payload       json.RawMessage `json:"payload"`
}

type Player struct {
	PlayerID uuid.UUID `json:"playerId"`
	Username string    `json:"username"`
}

type CardType string

const (
	CardTypeNormal = CardType("NORMAL")
	CardTypeIlmato = CardType("ILMATO")
	CardTypeIlbie  = CardType("ILBIE")
)

type CardSuit string

const (
	CardSuitSpades   = CardSuit("SPADES")
	CardSuitHearts   = CardSuit("HEARTS")
	CardSuitDiamonds = CardSuit("DIAMONDS")
	CardSuitClubs    = CardSuit("CLUBS")
	CardSuitTrumps   = CardSuit("TRUMPS")
)

type CardRank string

const (
	CardRankJack   = CardRank("JACK")
	CardRankKnight = CardRank("KNIGHT")
	CardRankQueen  = CardRank("QUEEN")
	CardRankKing   = CardRank("KING")
)

type Card struct {
	Type CardType `json:"type"`
	Suit CardSuit `json:"suit,omitempty"`
	Rank CardRank `json:"rank,omitempty"`
}

type PlayerJoined struct {
	GameID      uuid.UUID `json:"gameId"`
	Player      Player    `json:"player"`
	PlayerCount int       `json:"playerCount"`
	JoinedAt    time.Time `json:"joinedAt,omitempty"`
}

type PlayerLeft struct {
	GameID      uuid.UUID `json:"gameId"`
	PlayerID    uuid.UUID `json:"playerId"`
	PlayerCount int       `json:"playerCount,omitempty"`
	LeftAt      time.Time `json:"leftAt,omitempty"`
}

type PlayerDisconnected struct {
	GameID               uuid.UUID `json:"gameId"`
	PlayerID             uuid.UUID `json:"playerId"`
	DisconnectedAt       time.Time `json:"disconnectedAt"`
	ReconnectionDeadline time.Time `json:"reconnectionDeadline,omitempty"`
}

type GameStarted struct {
	GameID    uuid.UUID `json:"gameId"`
	Players   []Player  `json:"players"`
	StartedAt time.Time `json:"startedAt"`
}

type GameAbortedReason string

const (
	GameAbortedReasonPlayerDisconnected = GameAbortedReason("PLAYER_DISCONNECTED")
	GameAbortedReasonPlayerLeft         = GameAbortedReason("PLAYER_LEFT")
	GameAbortedReasonTimeout            = GameAbortedReason("TIMEOUT")
	GameAbortedReasonServerError        = GameAbortedReason("SERVER_ERROR")
)

type GameAborted struct {
	GameID    uuid.UUID         `json:"gameId"`
	Reason    GameAbortedReason `json:"reason"`
	AbortedAt time.Time         `json:"abortedAt"`
}

type GameRanking struct {
	PlayerID uuid.UUID `json:"playerId"`
	Score    int       `json:"score"`
	Rank     int       `json:"rank"`
}

type GameFinished struct {
	GameID     uuid.UUID     `json:"gameId"`
	Rankings   []GameRanking `json:"rankings"`
	FinishedAt time.Time     `json:"finishedAt"`
}

type TurnStarted struct {
	GameID        uuid.UUID `json:"gameId"`
	TurnNumber    int       `json:"turnNumber"`
	FirstBettorID uuid.UUID `json:"firstBettorId"`
	Hand          []Card    `json:"hand"`
}

type TurnFinished struct {
	GameID     uuid.UUID         `json:"gameId"`
	TurnNumber int               `json:"turnNumber"`
	Scores     map[uuid.UUID]int `json:"scores"`
}

type FoldStarted struct {
	GameID        uuid.UUID `json:"gameId"`
	TurnNumber    int       `json:"turnNumber"`
	FoldNumber    int       `json:"foldNumber"`
	FirstPlayerID uuid.UUID `json:"firstPlayerId"`
}

type PlayedCard struct {
	PlayerID uuid.UUID `json:"playerId"`
	Card     Card      `json:"card"`
}

type FoldFinished struct {
	GameID      uuid.UUID    `json:"gameId"`
	TurnNumber  int          `json:"turnNumber"`
	FoldNumber  int          `json:"foldNumber"`
	WinnerID    uuid.UUID    `json:"winnerId"`
	CardsPlayed []PlayedCard `json:"cardsPlayed"`
}

type BetAsk struct {
	GameID         uuid.UUID `json:"gameId"`
	TurnNumber     int       `json:"turnNumber"`
	PlayerID       uuid.UUID `json:"playerId"`
	MinBet         int       `json:"minBet"`
	MaxBet         int       `json:"maxBet"`
	TimeoutSeconds int       `json:"timeoutSeconds,omitempty"`
}

type BetPlaced struct {
	GameID        uuid.UUID  `json:"gameId"`
	TurnNumber    int        `json:"turnNumber"`
	PlayerID      uuid.UUID  `json:"playerId"`
	BetValue      int        `json:"betValue"`
	CorrelationID *uuid.UUID `json:"correlationId,omitempty"`
}

type BetCompleted struct {
	GameID     uuid.UUID         `json:"gameId"`
	TurnNumber int               `json:"turnNumber"`
	Bets       map[uuid.UUID]int `json:"bets"`
}

type PlayAsk struct {
	GameID         uuid.UUID `json:"gameId"`
	TurnNumber     int       `json:"turnNumber"`
	FoldNumber     int       `json:"foldNumber"`
	PlayerID       uuid.UUID `json:"playerId"`
	AllowedCards   []Card    `json:"allowedCards"`
	TimeoutSeconds int       `json:"timeoutSeconds,omitempty"`
}

type CardPlayed struct {
	GameID        uuid.UUID  `json:"gameId"`
	TurnNumber    int        `json:"turnNumber"`
	FoldNumber    int        `json:"foldNumber"`
	PlayerID      uuid.UUID  `json:"playerId"`
	Card          Card       `json:"card"`
	CorrelationID *uuid.UUID `json:"correlationId,omitempty"`
}
