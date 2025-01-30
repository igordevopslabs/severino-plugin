package main

import (
	"github.com/Kong/go-pdk/server"
	"github.com/igordevopslabs/severino-plugin/handler"
)

func main() {
	//Criar um server do KONG pdk e também registrar o plugin
	server.StartServer(handler.New, "0.0.1", 1000)
}
