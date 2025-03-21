package main

//import (
//	"encoding/json"
//	"fmt"
//	"github.com/SwanHtetAungPhyo/common/model"
//	"github.com/SwanHtetAungPhyo/common/protos"
//	"github.com/golang/protobuf/proto"
//	"github.com/sirupsen/logrus"
//	"io"
//	"net"
//	"os"
//	"os/signal"
//	"strings"
//	"syscall"
//	"time"
//)
//
//const (
//	serverAddress = ":9002"
//	outputFile    = "received_file.bin"
//)
//
//var log = logrus.New()
//
//// Animated loading effect
//func showLoadingAnimation(message string, duration time.Duration) {
//	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
//	start := time.Now()
//
//	for {
//		for _, frame := range frames {
//			fmt.Printf("\r%s %s", frame, message)
//			time.Sleep(100 * time.Millisecond)
//
//			// Stop after the given duration
//			if time.Since(start) > duration {
//				fmt.Printf("\r✅ %s\n", message) // ✅ Show completion message
//				return
//			}
//		}
//	}
//}
//
//func receiveFile(con net.Conn) error {
//	fmt.Println("\n📡 Connecting to main node...")
//	showLoadingAnimation("Downloading DAG data...", 3*time.Second)
//
//	file, err := os.Create(outputFile)
//	if err != nil {
//		log.Fatal("Error creating output file:", err)
//		return err
//	}
//	defer file.Close()
//
//	buffer := make([]byte, 2048)
//
//	n, err := con.Read(buffer)
//	if err != nil {
//		log.Fatal("Error reading from connection:", err)
//		return err
//	}
//
//	if n == 0 {
//		log.Error("No data received from the server.")
//		return fmt.Errorf("no data received")
//	}
//
//	_, err = file.Write(buffer[:n])
//	if err != nil {
//		log.Fatal("Error writing to file:", err)
//		return err
//	}
//	log.Info("📥 File received successfully, written to ", outputFile)
//	return nil
//}
//
//func DownloadTheMainNodeData() {
//	conn, err := net.Dial("tcp", serverAddress)
//	if err != nil {
//		log.Fatal("Error connecting to main node: ", err)
//	}
//	defer conn.Close()
//
//	err = receiveFile(conn)
//	if err != nil {
//		log.Fatal("Error receiving file: ", err)
//	} else {
//		log.Info("✅ DAG Data successfully downloaded!")
//	}
//}
//
//func LoadDAGFromFile() *protos.DAG {
//	fmt.Println("\n📂 Loading DAG structure from file...")
//	showLoadingAnimation("Parsing DAG file...", 3*time.Second)
//
//	file, err := os.Open("./received_file.bin")
//	if err != nil {
//		log.Fatalf("❌ Error opening file: %v", err)
//		return nil
//	}
//	defer file.Close()
//
//	receivedFileInfo, err := os.Stat("received_file.bin")
//	if err != nil {
//		log.Fatalf("Failed to get received file info: %v", err)
//	}
//	log.Infof("Received file size: %d bytes", receivedFileInfo.Size())
//
//	data, err := io.ReadAll(file)
//	if err != nil {
//		log.Fatalf("❌ Error reading file: %v", err)
//		return nil
//	}
//
//	var protoDag protos.DAG
//	err = proto.Unmarshal(data, &protoDag)
//	if err != nil {
//		log.Fatalf("❌ Error unmarshalling file: %v", err)
//		return nil
//	}
//
//	log.Println("📂 DAG data loaded successfully!")
//	return &protoDag
//}
//func handleConnection(conn net.Conn) {
//	defer func() {
//		if err := conn.Close(); err != nil {
//			log.Error("Error closing connection: ", err)
//		} else {
//			log.Info("Connection closed successfully.")
//		}
//	}()
//
//	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
//
//	buffer := make([]byte, 3000)
//
//	n, err := conn.Read(buffer)
//	if err != nil {
//		log.Error("Error reading from connection: ", err)
//		if err.Error() == "EOF" {
//			log.Warn("Connection closed by client (EOF).")
//		}
//		return
//	}
//
//	log.Infof("Received %d bytes of data from client", n)
//
//	var txs []model.Message
//	err = json.Unmarshal(buffer[:n], &txs)
//	if err != nil {
//		log.Error("Error unmarshalling transactions: ", err)
//		return
//	}
//}
//
//// ValidatorServer listens for transactions and validates them.
//func ValidatorServer(port string) error {
//	listen, err := net.Listen("tcp", ":"+port)
//	if err != nil {
//		log.Error("Error starting validator server: ", err)
//		return err
//	}
//	defer listen.Close()
//
//	log.Infof("Validator server listening on port: %s", port)
//
//	sigChan := make(chan os.Signal, 1)
//	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
//
//	go func() {
//		for {
//			conn, err := listen.Accept()
//			if err != nil {
//				log.Error("Error accepting connection: ", err)
//				continue
//			}
//
//			go handleConnection(conn)
//		}
//	}()
//
//	<-sigChan
//	log.Info("Received shutdown signal, shutting down the server gracefully...")
//	return nil
//}
//
//func main() {
//
//	validatorPort := os.Getenv("VALIDATOR_PORT")
//	if len(validatorPort) == 0 {
//		validatorPort = "9007"
//	}
//	DownloadTheMainNodeData()
//
//	go func() {
//		if err := ValidatorServer(validatorPort); err != nil {
//			log.Fatal(err)
//		}
//	}()
//
//	//go func() {
//	dag := LoadDAGFromFile()
//	fmt.Println("\n🚀 ===== DAG Structure ===== 🚀\n")
//	for _, block := range dag.Vertices {
//		fmt.Printf("🆔 Block ID: %v\n", block.Id)
//		fmt.Printf("🔗 Hash: %v\n", block.Hash)
//		fmt.Printf("⏰ Timestamp: %v\n", block.Timestamp)
//		fmt.Printf("🔗 Parents: %v\n", block.Parents)
//		fmt.Printf("🛡️ Validators: %v\n", block.Validators)
//
//		fmt.Println("📜 Transactions:")
//		for _, tx := range block.Transactions {
//			fmt.Printf("  🔹 Tx ID: %v\n", tx.TransactionId)
//			fmt.Printf("  ✍️ Signature: %v\n", tx.Signature)
//			fmt.Printf("  📦 Data: %v\n", tx.Data)
//			fmt.Printf("  👤 Action Taker: %v\n", tx.ActionTaker)
//			fmt.Printf("  🎯 Action Receiver: %v\n", tx.ActionReceiver)
//			fmt.Println("  " + strings.Repeat("⚡", 50))
//		}
//		fmt.Println(strings.Repeat("🟠", 50))
//	}
//	//}()
//
//	sigChan := make(chan os.Signal, 1)
//	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
//	<-sigChan
//}
