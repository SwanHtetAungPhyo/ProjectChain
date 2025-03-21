package http

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/json"
	"github.com/SwanHtetAungPhyo/common/model"
	"github.com/SwanHtetAungPhyo/ledchain/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"math/big"
	"net"
	"os"
)

// Handlers interface defines the account and transaction methods.
type Handlers interface {
	CreateAccount(ctx *fiber.Ctx) error
	ExecuteTransaction(ctx *fiber.Ctx) error
}

// HandlerImpl struct implements the Handlers interface.
type HandlerImpl struct {
	AccountService     *services.AccountService
	TransactionService *services.TransactionService
}

// NewHandler initializes a new HandlerImpl with dependencies.
func NewHandler(accountService *services.AccountService, transactionService *services.TransactionService) Handlers {
	return &HandlerImpl{
		AccountService:     accountService,
		TransactionService: transactionService,
	}
}

// CreateAccount godoc
// @Summary Create a new wallet account
// @Description Creates a new wallet and stores the private key in a cookie
// @Tags Account
// @Accept  json
// @Produce  json
// @Success 200 {object} ApiResponse "Successfully created account"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /chain/wallet [post]
func (h HandlerImpl) CreateAccount(ctx *fiber.Ctx) error {
	wallet, privateKey, err := h.AccountService.NewAccount()
	if err != nil {
		return SendErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	ctx.Cookie(&fiber.Cookie{
		Name:  "wallet_private_key",
		Value: privateKey,
	})

	return ctx.Status(fiber.StatusOK).JSON(ApiResponse{
		Success: true,
		Message: "Store the private key securely. This is the Beta mode, and security risk is at the highest level.",
		Data:    wallet,
	})
}

// ExecuteTransaction godoc
// @Summary Execute a blockchain transaction
// @Description Verifies, signs, and sends a transaction to the mempool
// @Tags Transaction
// @Accept  json
// @Produce  json
// @Success 200 {object} ApiResponse "Transaction successfully executed"
// @Failure 400 {object} ApiResponse "Invalid request data"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /chain/trans [post]
func (h HandlerImpl) ExecuteTransaction(ctx *fiber.Ctx) error {
	var message model.Message
	if err := json.Unmarshal(ctx.Body(), &message); err != nil {
		return SendErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	log.Info("Received a message: %s", message.TransactionRequests)
	log.Info("Received From : ", ctx.IP())

	tx := model.Transaction{
		TransactionId:  uuid.New().String(),
		ActionTaker:    message.TransactionRequests.ActionTaker,
		ActionReceiver: message.TransactionRequests.ActionReceiver,
		Signature:      message.TransactionRequests.Signature,
		Data:           message.TransactionRequests.Data,
	}
	validator := model.Validator{
		ValidatorAddress: ctx.IP(),
		ValidatorPubKey:  ctx.IP(),
		Stake:            1000,
	}
	block, err := model.AddNewBlock([]model.Transaction{tx}, validator)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	//TODO : Send to the  validator
	validAddrees := os.Getenv("VALIDATOR_ADDRESS")
	sendTransactionToValidator(&message, validAddrees)
	return ctx.Status(fiber.StatusOK).JSON(ApiResponse{
		Success: true,
		Message: "Transaction executed successfully and wait in the validator to be added to the DAG",
		Data:    block,
	})
}

func sendTransactionToValidator(transaction *model.Message, serverAddress string) error {
	data, err := json.Marshal(transaction)
	if err != nil {
		log.Error("Error marshaling transaction:", err)
		return err
	}

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Error("Error connecting to validator server:", err)
		return err
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Error("Error closing connection to validator server:", err)
		}
	}(conn)

	_, err = conn.Write(data)
	if err != nil {
		log.Error("Error sending transaction to validator server:", err)
		return err
	}

	log.Info("Transaction sent to validator server successfully")
	return nil
}

func JSONToPublicKey(publicKeyJSON model.PublicKeyJSON) (*ecdsa.PublicKey, error) {

	xBytes, err := base64.StdEncoding.DecodeString(publicKeyJSON.X)
	if err != nil {
		return nil, err
	}
	yBytes, err := base64.StdEncoding.DecodeString(publicKeyJSON.Y)
	if err != nil {
		return nil, err
	}
	publicKey := &ecdsa.PublicKey{
		Curve: elliptic.P384(),
		X:     new(big.Int).SetBytes(xBytes),
		Y:     new(big.Int).SetBytes(yBytes),
	}

	return publicKey, nil
}
