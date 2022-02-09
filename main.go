package main

import (
	"github.com/chadius/terosgameserver/rpc/github.com/chadius/teros_game_server"
	"github.com/chadius/terosgameserver/terosgameserver"
	"net/http"
)

func main() {
	server := terosgameserver.NewServer(nil)
	twirpHandler := teros_game_server.NewTerosGameServerServer(server)

	http.ListenAndServe(":8080", twirpHandler)
}
