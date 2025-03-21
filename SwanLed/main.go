package main

import (
	"context"
	"fmt"
	"github.com/SwanHtetAungPhyo/common/model"
	"github.com/SwanHtetAungPhyo/common/services"
	"github.com/SwanHtetAungPhyo/ledchain/cmd/http"
	"github.com/SwanHtetAungPhyo/ledchain/cmd/p2p"
	"github.com/SwanHtetAungPhyo/ledchain/cmd/validation_server"
	crons "github.com/SwanHtetAungPhyo/ledchain/internal/cron"
	"github.com/sirupsen/logrus"

	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var log = logrus.New()

func PrintDAG() {
	fmt.Println("\n===== DAG Structure =====\n")

	for _, block := range model.SwanDAG.Vertices {
		fmt.Printf("🔹 Block ID: %s\n", block.ID)
		fmt.Printf("   📌 Hash: %s\n", block.Hash)
		fmt.Printf("   ⏳ Timestamp: %s\n", block.Timestamp)
		fmt.Printf("   🏗 Parents: %v\n", block.Parents)
		fmt.Printf("   🔑 Validator: %s (Stake: %d)\n", block.Validators.ValidatorAddress, block.Validators.Stake)
		fmt.Println("   📜 Transactions:")
		for i, tx := range block.Transactions {
			fmt.Printf("      🔹 [%d] Transaction ID: %s\n", i+1, tx.TransactionId)
			fmt.Printf("         📍 Block Index: %d\n", tx.BlockIndex)
			fmt.Printf("         💰 Data: %v\n", tx.Data)
			fmt.Printf("         ✍ Signature: %s\n", tx.Signature)
		}
		fmt.Println(strings.Repeat("-", 50))
	}
}

func Main() {
	model.InitDAG()
	crons.SaveDAGToFile()
	dag := crons.LoadDAGFromFile()
	fmt.Println("\n🚀 ===== DAG Structure ===== 🚀\n")
	for _, block := range dag.Vertices {
		fmt.Printf("🆔 Block ID: %v\n", block.Id)
		fmt.Printf("🔗 Hash: %v\n", block.Hash)
		fmt.Printf("⏰ Timestamp: %v\n", block.Timestamp)
		fmt.Printf("🔗 Parents: %v\n", block.Parents)
		fmt.Printf("🛡️ Validators: %v\n", block.Validators)

		fmt.Println("📜 Transactions:")
		for _, tx := range block.Transactions {
			fmt.Printf("  🔹 Tx ID: %v\n", tx.TransactionId)
			fmt.Printf("  ✍️ Signature: %v\n", tx.Signature)
			fmt.Printf("  📦 Data: %v\n", tx.Data)
			fmt.Printf("  👤 Action Taker: %v\n", tx.ActionTaker)
			fmt.Printf("  🎯 Action Receiver: %v\n", tx.ActionReceiver)
			fmt.Println("  " + strings.Repeat("⚡", 50))
		}
		fmt.Println(strings.Repeat("🟠", 50))
	}

}

func main() {
	model.InitDAG()

	PORT := os.Getenv("PORT")
	HTTP_PORT := os.Getenv("HTTP_PORT")
	PEERS := os.Getenv("PEERS")
	VALIDATOR_ADDRESS := os.Getenv("VALIDATOR_ADDRESS")
	REGISTRY_PORT := os.Getenv("REGISTRY_PORT")

	if PORT == "" || HTTP_PORT == "" {
		panic("[ERROR] PORT and HTTP_PORT environment variables not set")
	}

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
	services.AnimationLoop("💻 Full node is trying to start", 4*time.Second)
	go func() {
		log.Info("[INFO] Starting up the Http Server to interact with the DAG ")
		http.Start(HTTP_PORT)
	}()

	time.Sleep(1 * time.Second) // Give some time for servers to initialize
	go func() {
		log.Info("[INFO] Starting up the P2P Server to interact with the Peer")
		log.Info("[INFO] P2P protocol is currently the WebSocket")
		p2p.Start(PORT, PEERS)

	}()

	log.Info("[INFO] Starting up the MemPool Server to cache the transactions")
	log.Info("[INFO] MemPool Protocol is currently TCP.")
	time.Sleep(1 * time.Second)
	log.Info("[INFO] Cron Jobs to take snapshots of the blockchain are running every 1 minute in the background")
	go crons.SetupCronJob()

	go func() {
		validation_server.Start(REGISTRY_PORT)
	}()
	go func() {
		printInfo(PORT, HTTP_PORT, PEERS, VALIDATOR_ADDRESS)
	}()
	sigReceived := <-osSignal
	log.Info(fmt.Sprintf("[INFO] Received signal: %s, initiating graceful shutdown...", sigReceived))

	log.Info("[INFO] Shutting down...")
	time.Sleep(1 * time.Second)
	log.Info("[INFO] Shutdown complete.")
}

func printInfo(PORT, HTTP_PORT, PEERS, VALIDATORADDRESS string) {
	lines := []string{
		"🌐 Application Layer Port:         http://localhost:" + HTTP_PORT,
		"🔗 Node Communication Peer Port:   " + PORT,
		"📡 Peer Nodes in Network:          " + PEERS,
		"✅ Validator Node Port:            " + VALIDATORADDRESS,
	}
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}
	border := "┌" + strings.Repeat("─", maxLen+2) + "┐"
	fmt.Println(border)
	for _, line := range lines {
		fmt.Printf("│ %-*s │\n", maxLen, line)
	}
	fmt.Println("└" + strings.Repeat("─", maxLen+2) + "┘")
}
