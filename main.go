package main

import (
	"github.com/Chadius/terosGameServer/rpc/github.com/chadius/teros_game_server"
	"github.com/Chadius/terosGameServer/terosgameserver"
	"net/http"
)

func main() {
	server := terosgameserver.NewServer(nil)
	twirpHandler := teros_game_server.NewTerosGameServerServer(server)

	http.ListenAndServe(":8080", twirpHandler)
}
