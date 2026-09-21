package server

import (
	"net/http"

	"fr.vamary.whist-server/internal/server"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (j server.JoinLobbyRequestObject) CreateLobby(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (j server.JoinLobbyJSONRequestBody) DestroyLobby(w http.ResponseWriter, r *http.Request, lobbyID openapi_types.UUID) {
	//TODO implement me
	panic("implement me")
}

func (j server.JoinLobbyJSONRequestBody) GetLobby(w http.ResponseWriter, r *http.Request, lobbyID openapi_types.UUID) {
	//TODO implement me
	panic("implement me")
}

func (j server.JoinLobbyJSONRequestBody) StartGame(w http.ResponseWriter, r *http.Request, lobbyID openapi_types.UUID) {
	//TODO implement me
	panic("implement me")
}

func (j server.JoinLobbyJSONRequestBody) JoinLobby(w http.ResponseWriter, r *http.Request, lobbyID openapi_types.UUID) {
	//TODO implement me
	panic("implement me")
}

func (j server.JoinLobbyJSONRequestBody) LeaveLobby(w http.ResponseWriter, r *http.Request, lobbyID openapi_types.UUID) {
	//TODO implement me
	panic("implement me")
}
