package server

import (
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

type ServerImplementation struct{}

func NewServer() *ServerImplementation {
	return &ServerImplementation{}
}

func (s *ServerImplementation) CreateLobby(w http.ResponseWriter, r *http.Request) {

}

func (s *ServerImplementation) GetLobby(w http.ResponseWriter, r *http.Request, lobbyID openapi_types.UUID) {

}

func (s *ServerImplementation) DestroyLobby(w http.ResponseWriter, r *http.Request, lobbyID openapi_types.UUID) {

}

func (s *ServerImplementation) JoinLobby(w http.ResponseWriter, r *http.Request, lobbyID openapi_types.UUID) {

}

func (s *ServerImplementation) LeaveLobby(w http.ResponseWriter, r *http.Request, lobbyID openapi_types.UUID) {

}

func (s *ServerImplementation) StartGame(w http.ResponseWriter, r *http.Request, lobbyID openapi_types.UUID) {

}
