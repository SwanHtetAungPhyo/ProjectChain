package services

import (
	"encoding/json"
	"github.com/SwanHtetAungPhyo/common/model"
	"github.com/sirupsen/logrus"
	"net"
)

var log = logrus.New()

// SendToMemPool sends the transaction to the mempool via a TCP connection
func SendToMemPool(Tx *model.Transaction, port string) error {
	log.Infof("Connecting to mempool on port: %s", port)

	conn, err := net.Dial("tcp", ":"+port)
	if err != nil {
		log.Error("Error connecting to mempool: ", err)
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Error("Error closing connection: ", err)
		}
	}()

	encodedJson, err := json.MarshalIndent(Tx, "", "  ")
	if err != nil {
		log.Error("Error marshalling transaction: ", err)
		return err
	}

	_, err = conn.Write(encodedJson)
	if err != nil {
		log.Error("Error writing to mempool: ", err)
		return err
	}

	log.Infof("Transaction successfully sent to mempool on port %s", port)
	return nil
}
