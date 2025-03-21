package validation_server

import (
	"encoding/hex"
	"fmt"
	"github.com/SwanHtetAungPhyo/common/model"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	log              = logrus.New()
	ValidatorHashMap = make(map[string]bool)
	validatorInfoMap = make(map[string]model.ValidatorMetaData)
	registryMap      = make(map[string]model.ValidatorMessage)
	validatorLock    sync.Mutex
)

// logWithPrefix returns a logrus.Entry with a given prefix field.
func logWithPrefix(name string) *logrus.Entry {
	return log.WithFields(logrus.Fields{
		"prefix": name,
	})
}

// Start initializes and starts the Fiber server.
func Start(port string) {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	setupValidatorRegistry(app)
	logWithPrefix("Validator Registry").Infoln("Starting validator Registry")
	if err := app.Listen(":" + port); err != nil {
		panic(err)
	}
}

// setupValidatorRegistry creates the validator registry endpoint.
func setupValidatorRegistry(app *fiber.App) {
	registry := app.Group("registry")
	registry.Post("/:id", validatorRegistrationHandler, limiter.New(limiter.Config{
		Expiration: 5 * time.Minute,
		Max:        5,
	}))
	infoSetUp := registry.Group("infoset")
	infoSetUp.Post("/:id", setUpTheMetaData, limiter.New(limiter.Config{}))
}

func setUpTheMetaData(ctx *fiber.Ctx) error {
	signature := ctx.Get("signature") // Extract the signature from headers

	if signature == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Missing signature header",
		})
	}
	var setUpData model.ValidatorMetaData
	id := ctx.Params("id")

	if err := ctx.BodyParser(&setUpData); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	validatorLock.Lock()
	validatorInfoMap[id] = setUpData
	validatorLock.Unlock()
	return ctx.Status(200).JSON(fiber.Map{
		"success": true,
		"data":    setUpData,
	})
}

// validatorRegistrationHandler handles new validator registration.
func validatorRegistrationHandler(ctx *fiber.Ctx) error {
	var validator model.ValidatorMessage
	if err := ctx.BodyParser(&validator); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	validatorId := ctx.Params("id")
	if validatorId == "" {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	log.WithFields(logrus.Fields{
		"ip":        validator.IP,
		"pubkey":    validator.PubKey,
		"signature": validator.Signature,
	}).Info("Validation attempt")

	if _, exists := ValidatorHashMap[validatorId]; exists {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Validator ID already in use",
		})
	}

	if !VerifyValidatorMessage(validator.PubKey, validator, validator.Signature) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid signature",
		})
	}

	ValidatorHashMap[validatorId] = true
	registryMap[validatorId] = validator
	logWithPrefix("Validator Registry").Infof("New validator registered: %s", validatorId)
	return ctx.SendFile("dag_data.json")
}

// VerifyValidatorMessage verifies the signature of a ValidatorMessage.
// It uses only the signable fields (IP, Password, Version) to construct the message.
func VerifyValidatorMessage(publicKeyHex string, msg model.ValidatorMessage, signatureHex string) bool {
	// Decode the signature from hex.
	signatureBytes, err := hex.DecodeString(signatureHex)
	if err != nil {
		log.Errorf("Error decoding signature: %s", err)
		return false
	}

	// Decode the public key.
	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		log.Errorf("Error decoding public key: %s", err)
		return false
	}

	pubKey, err := crypto.UnmarshalPubkey(publicKeyBytes)
	if err != nil {
		log.Errorf("Error unmarshalling public key: %s", err)
		return false
	}

	// Build the signable message from the common model.
	signableMsg := model.SignableValidatorMessage{
		IP:       msg.IP,
		Password: msg.Password,
		Version:  msg.Version,
	}

	jsonData, err := json.Marshal(signableMsg)
	if err != nil {
		log.Errorf("Error encoding signable message: %s", err)
		return false
	}

	// Construct the Ethereum signed message.
	message := fmt.Sprintf("Ethereum Signed Message:\n%d%s", len(jsonData), jsonData)
	hash := crypto.Keccak256Hash([]byte(message))

	return crypto.VerifySignature(
		crypto.FromECDSAPub(pubKey),
		hash.Bytes(),
		signatureBytes[:64],
	)
}
