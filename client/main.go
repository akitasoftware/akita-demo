package main

import (
	"akitasoftware.com/demo-client/datasource"
	_ "embed"
	"flag"
	"log"
	"net/http"
)

//go:embed application.yml
var applicationYML []byte

func main() {
	flag.Parse()

	config, err := ParseConfiguration(applicationYML)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	server := datasource.NewDemoServer("http://akita-demo-server:8080", http.DefaultClient)

	app := NewApp(config, server)

	app.SendEvent("demo-client-started", map[string]any{})

	app.HandleDemoTasks()
}
