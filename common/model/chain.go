package model

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"
	"time"
)

var SwanDAG *DAG

type DAG struct {
	Vertices map[string]*Block `json:"vertices"`
}

func InitDAG() {
	if SwanDAG == nil {
		SwanDAG = &DAG{Vertices: make(map[string]*Block)}
	}

	genesisAccount, privateKey, err := NewAccount()
	if err != nil {
		log.Fatal("[ERROR] Failed to generate genesis account", err)
	}

	tx := &Transaction{
		TransactionId:  uuid.NewString(),
		ActionTaker:    genesisAccount.PubKey,
		ActionReceiver: genesisAccount.PubKey,
		Data: map[string]interface{}{
			"address": genesisAccount.Address,
			"type":    "transfer",
			"amount":  "1000",
		},
	}

	err = SignTransaction(tx, privateKey)
	if err != nil {
		panic(err)
	}

	if !VerifyTransaction(*tx) {
		panic("[ERROR] Failed to verify transaction")
	}
	genesisBlock := &Block{
		ID:           uuid.NewString(),
		Timestamp:    time.Now().Format("2006-01-02 15:04:05"),
		Parents:      []string{},
		Transactions: []Transaction{*tx},
		Validators: Validator{
			ValidatorAddress: genesisAccount.Address,
			ValidatorPubKey:  genesisAccount.PubKey,
			Stake:            100,
		},
	}
	genesisBlock.Hash = GenerateBlockHash(genesisBlock.ID, nil, nil)
	SwanDAG.Vertices[genesisBlock.ID] = genesisBlock
	jsonSerialized, _ := json.Marshal(genesisBlock)
	blockSize := len(jsonSerialized)
	blockSizeMB := float64(blockSize) / 1024 / 1024
	log.Println("[INFO] Block size:", blockSizeMB, "MB")
}

func AddNewBlock(transactions []Transaction, validator Validator) (block *Block, err error) {
	if !VerifyTransaction(transactions[0]) {
		return nil, errors.New("[ERROR] Failed to verify transaction")
	}
	if SwanDAG == nil {
		return nil, errors.New("[ERROR] Swan DAG not initialized")
	}
	parentIds := getLatestBlockIDs()
	newBlockID := uuid.New().String()
	newBlock := &Block{
		ID:           newBlockID,
		Timestamp:    time.Now().Format("2006-01-02 15:04:05"),
		Hash:         GenerateBlockHash(newBlockID, parentIds, transactions),
		Parents:      parentIds,
		Transactions: transactions,
		Validators:   validator,
	}
	SwanDAG.Vertices[newBlock.ID] = newBlock
	return newBlock, nil
}

func getLatestBlockIDs() []string {
	var latestBlockIDs []string

	for blockID, block := range SwanDAG.Vertices {
		parseTime, err := time.Parse("2006-01-02 15:04:05", block.Timestamp)
		if err != nil {
			continue
		}

		if time.Since(parseTime).Minutes() < 1 || isTipBlock(blockID) {
			latestBlockIDs = append(latestBlockIDs, blockID)
		}
		if len(latestBlockIDs) >= 3 {
			break
		}
	}

	return latestBlockIDs
}

func isTipBlock(blockID string) bool {
	for _, block := range SwanDAG.Vertices {
		for _, parentID := range block.Parents {
			if parentID == blockID {
				return false
			}
		}
	}
	return true
}

func parseTimestamp(timestamp string) int64 {
	parsed, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return 0
	}
	return parsed.Unix()
}

func GetRandomBlockID() string {
	var blocks []*Block
	//for _, block := range SwanDAG.Vertices {
	//	if len(block.Parents) > 0 {
	//		blocks = append(blocks, block)
	//	}
	//}
	//if len(blocks) == 0 {
	//	return ""
	//}
	//rand.Seed(time.Now().UnixNano())
	//index := rand.Intn(len(SwanDAG.Vertices))
	return blocks[0].ID
}
