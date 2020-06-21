package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/shukla2112/counter-app/config"
	"github.com/shukla2112/counter-app/server"
)

// app details
const (
	Name    = "counter-app"
	Usage   = "Service to keep the counter for the given key - it's redis based counter"
	Version = "0.01"
)

func main() {
	// Read in CLI args
	cliArgs, jsonStr := config.ReadCliArgs()
	log.Printf("COUNTER_1001: Starting %s service, version %s, mode %s with cli args %s\n", Name, Version, cliArgs.Mode, jsonStr)

	// Create new app config
	appC, err := config.NewAppConfig(&cliArgs)
	if err != nil {
		panic(fmt.Errorf("COUNTER_1002EF: Error creating new app object: %v", err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go server.StartServer(ctx, appC)

	sigInt := make(chan os.Signal, 1)
	signal.Notify(sigInt, os.Interrupt)
	for range sigInt {
		log.Println("COUNTER_1003: Received SigInt. stopping all services...")
		cancel()
		time.Sleep(5 * time.Second)
		break
	}
}
