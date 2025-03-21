package main

import (
	"encoding/json"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/SwanHtetAungPhyo/mempool/model"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()
var pool = model.Pool{}
var poolMutex sync.Mutex

func SendToValidator(port string, txs []model.Transaction) {
	conn, err := net.Dial("tcp", ":"+port)
	if err != nil {
		log.Error("Error connecting to validator: ", err)
		return
	}
	defer conn.Close()

	txJSON, err := json.Marshal(txs)
	if err != nil {
		log.Error("Error marshalling transactions: ", err)
		return
	}

	_, err = conn.Write(txJSON)
	if err != nil {
		log.Error("Error sending transactions to validator: ", err)
		return
	}
	log.Info("Sent transactions to validator")
}

func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Error("Error closing connection: ", err)
		}
	}()

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	buffer := make([]byte, 2048)

	n, err := conn.Read(buffer)
	if err != nil {
		log.Error("Error reading from connection: ", err)
		return
	}

	var tx model.Transaction
	err = json.Unmarshal(buffer[:n], &tx)
	if err != nil {
		log.Error("Error unmarshalling transaction: ", err)
		return
	}

	poolMutex.Lock()
	pool.AddTransaction(tx)
	currentSize := len(pool.GetTransactions())

	if currentSize == 3 {
		txsToSend := make([]model.Transaction, 3)
		copy(txsToSend, pool.GetTransactions())
		pool.CleanMemPool()
		poolMutex.Unlock()

		go SendToValidator("3001", txsToSend)
	} else {
		poolMutex.Unlock()
	}

	log.Infof("Received transaction: %+v", tx)
	log.Info("Transaction added to the mempool.")
}

func MempoolServer(port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Error("Error starting server: ", err)
		return err
	}
	defer listen.Close()

	log.Infof("Mempool server listening on port: %s", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Error("Error accepting connection: ", err)
				continue
			}

			go handleConnection(conn)

			poolMutex.Lock()
			log.Infof("Mempool size before check: %d", len(pool.GetTransactions()))
			poolMutex.Unlock()

			poolMutex.Lock()
			if len(pool.GetTransactions()) == 3 {
				pool.CleanMemPool()
				log.Info("Mempool contains 3 transactions, sending to validator...")
				SendToValidator("3001", pool.GetTransactions())
			}
			poolMutex.Unlock()
		}
	}()

	<-sigChan
	log.Info("Received shutdown signal, shutting down the server gracefully...")
	return nil
}

func main() {
	go func() {
		if err := MempoolServer("8080"); err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Info("Received shutdown signal, shutting down the server gracefully...")
}
