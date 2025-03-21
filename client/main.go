package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/SwanHtetAungPhyo/common/crytpography"
	"github.com/SwanHtetAungPhyo/common/model"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
	"strings" // Import the strings package for TrimSpace
)

var log = logrus.New()

const (
	saltLength  = 16
	iterations  = 3
	keyLength   = 32 // AES-256 requires 32-byte key
	nonceLength = 12 // GCM nonce length
)

// GenerateKeyPair generates a new ECDSA key pair and returns the private and public keys as hex strings.
func GenerateKeyPair() (privateKeyHex string, publicKeyHex string, err error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}

	privateKeyHex = hex.EncodeToString(privateKey.D.Bytes())
	publicKeyHex = hex.EncodeToString(privateKey.PublicKey.X.Bytes()) + hex.EncodeToString(privateKey.PublicKey.Y.Bytes())

	return privateKeyHex, publicKeyHex, nil
}

// CreateTransaction creates a new transaction with the given public key and signs it with the private key.
func CreateTransaction(privateKeyHex, publicKeyHex string) (model.Message, error) {
	tx := model.Transaction{
		TransactionId:  uuid.NewString(),
		ActionTaker:    publicKeyHex,
		ActionReceiver: publicKeyHex,
		Data: map[string]interface{}{
			"message": "Hello World!",
		},
	}

	if err := model.SignTransaction(&tx, privateKeyHex); err != nil {
		return model.Message{}, err
	}

	return model.Message{
		TransactionRequests: tx,
		PublicKey:           publicKeyHex,
	}, nil
}

func SendTransactionDemo() {
	privateKeyHex, publicKeyHex, err := GenerateKeyPair()
	if err != nil {
		log.Fatal(err)
	}
	tx := model.Transaction{
		TransactionId:  uuid.NewString(),
		ActionTaker:    publicKeyHex,
		ActionReceiver: publicKeyHex,
		Data: map[string]interface{}{
			"message": "Hello World!",
		},
	}
	if err := model.SignTransaction(&tx, privateKeyHex); err != nil {
		log.Fatal(err)
	}
	message := model.Message{
		TransactionRequests: tx,
		PublicKey:           publicKeyHex,
	}
	err = SendTransaction(message)
	if err != nil {
		return
	}
}

// SendTransaction sends the transaction to the API.
func SendTransaction(message model.Message) error {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+os.Getenv("TOKEN")).
		SetBody(message).
		Post("http://localhost:3003/chain/trans")
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		log.Errorf("API call failed with status: %v", resp.Error())
		return fmt.Errorf("API call failed with status: %v", resp.Status())
	}

	log.Info("Transaction sent successfully")
	return nil
}

// HandleCreateWallet handles the "Create Wallet" option.
func HandleCreateWallet(password string) (string, string) {
	privateKeyHex, publicKeyHex, err := GenerateKeyPair()
	if err != nil {
		log.Fatalf("Failed to generate key pair: %v", err)
	}

	log.Infof("Private Key: %s", privateKeyHex)
	log.Infof("Public Key: %s", publicKeyHex)
	encryptedPrivateKeyBytes, err := crytpography.EncryptPrivateKeyHex(publicKeyHex, []byte(password))
	if err != nil {
		log.Fatalf("Failed to encrypt private key: %v", err)
	}
	saveToFile(encryptedPrivateKeyBytes, publicKeyHex, "./private_key.pem")
	return privateKeyHex, publicKeyHex
}

func saveToFile(privateKeyHex, publicKeyHex, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Failed to close file: %v", err)
		}
	}(file)
	log.Infof("Writing private key to file: %s", privateKeyHex)
	_, err = file.WriteString("Private Key:\n" + privateKeyHex + "\n\nPublic Key:\n" + publicKeyHex)
	if err != nil {
		log.Fatalf("Failed to write private key to file: %v", err)
	}
	log.Infof("Writing public key to file: %s", publicKeyHex)
}

// HandleMakeTransaction handles the "Make Transaction" option.
func HandleMakeTransaction(privateKeyHex, publicKeyHex, password string) {
	var err error
	if err != nil {
		log.Fatalf("Failed to generate key pair: %v", err)
	}

	message, err := CreateTransaction(privateKeyHex, publicKeyHex)
	if err != nil {
		log.Fatalf("Failed to create transaction: %v", err)
	}

	if err := SendTransaction(message); err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}
}

// Survey prompts the user to choose an action and handles the selected option.
func Survey() {
	var choice string
	for {
		prompt := &survey.Select{
			Message: "Choose the wallet action:",
			Options: []string{"Create Wallet", "Make Transaction", "Simulation", "Exit"},
		}

		if err := survey.AskOne(prompt, &choice); err != nil {
			log.Fatalf("Failed to get user input: %v", err)
		}

		switch choice {
		case "Create Wallet":
			createWallet()
		case "Make Transaction":
			handleTransaction()
		case "Simulation":
			SendTransactionDemo()
		case "Exit":
			return
		default:
			log.Fatal("Invalid choice")
		}
	}
}

func handleTransaction() {
	var password string
	var privateKeyHex, publicKeyHex string

	prompt := &survey.Password{
		Message: "Password:",
	}
	if err := survey.AskOne(prompt, &password); err != nil {
		log.Fatalf("Failed to get user input: %v", err)
	}
	password = strings.TrimSpace(password) // Trim spaces from password

	prompt = &survey.Password{
		Message: "Private Key:",
	}
	if err := survey.AskOne(prompt, &privateKeyHex); err != nil {
		log.Fatalf("Failed to get user input: %v", err)
	}
	privateKeyHex = strings.TrimSpace(privateKeyHex) // Trim spaces from private key

	prompt = &survey.Password{
		Message: "Public Key:",
	}
	if err := survey.AskOne(prompt, &publicKeyHex); err != nil {
		log.Fatalf("Failed to get user input: %v", err)
	}
	publicKeyHex = strings.TrimSpace(publicKeyHex) // Trim spaces from public key

	HandleMakeTransaction(privateKeyHex, publicKeyHex, password)
}

func createWallet() {
	var userInput string
	prompt := &survey.Password{
		Message: "Enter your wallet password:",
	}
	if err := survey.AskOne(prompt, &userInput); err != nil {
		log.Fatalf("Failed to get user input: %v", err)
	}
	userInput = strings.TrimSpace(userInput)

	HandleCreateWallet(userInput)
}

func main() {
	Survey()
}
