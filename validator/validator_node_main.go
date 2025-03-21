package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/SwanHtetAungPhyo/common/model"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-resty/resty/v2"
	"github.com/goccy/go-json"
	"github.com/sirupsen/logrus"
	"github.com/tyler-smith/go-bip39"
	"os"
	"os/signal"
	"time"
)

type Account struct {
	PrivateKey string
	PublicKey  string
}

type SignableValidatorMessage struct {
	IP       string `json:"IP"`
	Password string `json:"Password"`
	Version  string `json:"Version"`
}

type ValidatorMessage struct {
	IP        string `json:"IP"`
	Password  string `json:"Password"`
	Version   string `json:"Version"`
	Signature string `json:"Signature"`
	PubKey    string `json:"PubKey"`
}

var validatorAccount = new(Account)
var log = logrus.New()

func WalletForValidator() (string, string) {
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	seed := bip39.NewSeed(mnemonic, "")
	privateKey, _ := crypto.ToECDSA(seed[:32])

	privateKeyHex := hex.EncodeToString(crypto.FromECDSA(privateKey))
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	publicKeyHex := hex.EncodeToString(crypto.FromECDSAPub(publicKey))

	return privateKeyHex, publicKeyHex
}

func SignTheValidatorMessage(privateKeyHex string, msg ValidatorMessage) (string, bool) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Errorf("Error decoding private key: %s", err)
		return "", false
	}

	signableMsg := SignableValidatorMessage{
		IP:       msg.IP,
		Password: msg.Password,
		Version:  msg.Version,
	}

	jsonData, _ := json.Marshal(signableMsg)
	message := fmt.Sprintf("Ethereum Signed Message:\n%d%s", len(jsonData), jsonData)
	hash := crypto.Keccak256Hash([]byte(message))

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Errorf("Error signing message: %s", err)
		return "", false
	}

	return hex.EncodeToString(signature), true
}

func UI() {
	for {
		fmt.Print("\033[H\033[2J")
		fmt.Println("\n=== Validator Node Setup ===")

		var choice string
		prompt := &survey.Select{
			Message: "Select an option:",
			Options: []string{"🔑 Create Wallet", "📝 Request QueueID", "❌ Exit"},
		}
		survey.AskOne(prompt, &choice)

		switch choice {
		case "🔑 Create Wallet":
			privateKey, publicKey := WalletForValidator()
			fmt.Printf("\nPrivate Key: %s\nPublic Key: %s\n", privateKey, publicKey)
			validatorAccount.PrivateKey = privateKey
			validatorAccount.PublicKey = publicKey

		case "📝 Request QueueID":
			RequestToParticipate()

		case "❌ Exit":
			os.Exit(0)
		}
	}
}

func RequestToParticipate() {
	privateKey, publicKey := WalletForValidator()
	msg := ValidatorMessage{
		IP:       "0.0.0.0",
		Password: "SwanHtet12@",
		Version:  "0.0.1",
		PubKey:   publicKey,
	}
	infoMeta := model.ValidatorMetaData{
		Address:       validatorAccount.PublicKey,
		SolanaAddress: validatorAccount.PrivateKey,
		Stake:         100,
	}

	signature, _ := SignTheValidatorMessage(privateKey, msg)
	msg.Signature = signature

	client := resty.New().
		SetTimeout(5 * time.Second).
		SetDebug(true). // Enable HTTP traffic logging
		EnableTrace()   // Enable request tracing
	resp, err := client.R().SetBody(msg).Post("http://localhost:4000/registry/Validator_one")
	if err != nil || resp.StatusCode() != 200 {
		log.Errorf("Request failed: %v", err)
		return
	}

	resp, err = client.R().SetBody(msg).
		SetHeader("Content-Type", "application/json").
		SetHeader("signature", signature).
		SetBody(infoMeta).
		Post("http://localhost:4000/infoset/Validator_one")
	if err != nil || resp.StatusCode() != 200 {
		log.Errorf("Request failed: %v", err)
		return
	}
	os.WriteFile("response.json", resp.Body(), 0644)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}

func main() {
	UI()
}
