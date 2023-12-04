package main

import (
	"semi_systems/attendance/adopter/gateway/mysql"
	"semi_systems/chat"
	"semi_systems/cmd/api"
)

func main() {
	mysql.InitDatabase()
	hub := chat.NewHub()
	go hub.RunLoop()
	api.Execute(hub)
}
