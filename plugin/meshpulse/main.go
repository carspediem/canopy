package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/canopy-network/go-plugin/contract"
)

func main() {
	log.Println("🚀 MeshPulse DePIN Node starting…")

	// Load config — reads chain.json if present, otherwise uses defaults
	cfg := contract.DefaultConfig()
	if len(os.Args) > 1 {
		var err error
		cfg, err = contract.NewConfigFromFile(os.Args[1])
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}
	}

	// Start the Canopy plugin (connects to FSM socket, starts UI server on :8080)
	contract.StartPlugin(cfg)

	log.Println("✅ MeshPulse ready — UI at http://localhost:8080")

	// Block until killed
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	log.Println("MeshPulse shutting down…")
}
