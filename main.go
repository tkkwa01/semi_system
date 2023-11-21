package main

import (
	"semi_systems/attendance/adopter/gateway/mysql"
	"semi_systems/cmd/api"
)

func main() {
	mysql.InitDatabase()
	api.Execute()
}
