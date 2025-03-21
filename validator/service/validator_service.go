package service

import (
	"errors"
	model2 "github.com/SwanHtetAungPhyo/common/model"
)

func ExecuteTransaction(transaction *model2.Transaction, privateKey string) error {

	if err := model2.SignTransaction(transaction, privateKey); err != nil {
		return err
	}

	if !model2.VerifyTransaction(*transaction) {
		return errors.New("transaction verification failed")
	}

	validator := model2.Validator{
		ValidatorPubKey:  privateKey,
		ValidatorAddress: privateKey,
	}
	if _, err := model2.AddNewBlock([]model2.Transaction{*transaction}, validator); err != nil {
		return err
	}

	return nil
}
