package converter

import (
	model2 "github.com/SwanHtetAungPhyo/common/model"
	protos "github.com/SwanHtetAungPhyo/common/protos"
)

func FromProtoBlock(protoBlock *protos.Block) *model2.Block {
	goBlock := model2.Block{
		ID:        protoBlock.Id,
		Timestamp: protoBlock.Timestamp,
		Hash:      protoBlock.Hash,
		Parents:   protoBlock.Parents,
	}

	// Convert Transactions
	for _, protoTx := range protoBlock.Transactions {
		goBlock.Transactions = append(goBlock.Transactions, model2.Transaction{
			TransactionId:  protoTx.TransactionId,
			ActionTaker:    protoTx.ActionTaker,
			ActionReceiver: protoTx.ActionReceiver,
			Data:           protoTx.Data,
			BlockIndex:     protoTx.BlockIndex,
			Signature:      protoTx.Signature,
		})
	}

	// Convert Validator
	if protoBlock.Validators != nil {
		goBlock.Validators = model2.Validator{
			ValidatorAddress: protoBlock.Validators.ValidatorAddress,
			ValidatorPubKey:  protoBlock.Validators.ValidatorPubKey,
			Stake:            protoBlock.Validators.Stake,
		}
	}

	return &goBlock
}

func FromProtoDAG(protoDAG *protos.DAG) model2.DAG {
	goDAG := model2.DAG{
		Vertices: make(map[string]*model2.Block),
	}

	// Convert each block in the Protobuf DAG to Go Block
	for id, protoBlock := range protoDAG.Vertices {
		goDAG.Vertices[id] = FromProtoBlock(protoBlock)
	}

	return goDAG
}

// ToProtoBlock converts a Go Block to a Protobuf Block
func ToProtoBlock(goBlock model2.Block) *protos.Block {
	protoBlock := &protos.Block{
		Id:        goBlock.ID,
		Timestamp: goBlock.Timestamp,
		Hash:      goBlock.Hash,
		Parents:   goBlock.Parents,
	}

	protoBlock.Transactions = make([]*protos.Transaction, len(goBlock.Transactions))
	for i, tx := range goBlock.Transactions {
		protoBlock.Transactions[i] = &protos.Transaction{
			TransactionId:  tx.TransactionId,
			ActionTaker:    tx.ActionTaker,
			ActionReceiver: tx.ActionReceiver,
			BlockIndex:     tx.BlockIndex,
			Signature:      tx.Signature,
		}
	}

	protoBlock.Validators = &protos.Validator{
		ValidatorAddress: goBlock.Validators.ValidatorAddress,
		ValidatorPubKey:  goBlock.Validators.ValidatorPubKey,
		Stake:            goBlock.Validators.Stake,
	}

	return protoBlock
}

// ToProtoDAG converts a Go DAG to a Protobuf DAG
func ToProtoDAG(goDAG model2.DAG) *protos.DAG {
	protoDAG := &protos.DAG{
		Vertices: make(map[string]*protos.Block),
	}

	// Convert each Block in the DAG to Protobuf Block
	for id, goBlock := range goDAG.Vertices {
		protoDAG.Vertices[id] = ToProtoBlock(*goBlock)
	}

	return protoDAG
}
