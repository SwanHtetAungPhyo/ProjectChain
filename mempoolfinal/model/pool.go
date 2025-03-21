package model

import "sync"

type Pool struct {
	Transactions []Transaction `json:"transactions"`
	Mu           sync.Mutex
}

type Transaction struct {
	TransactionId  string      `json:"transactionId"`
	ActionTaker    string      `json:"actionTaker"`
	ActionReceiver string      `json:"actionReceiver"`
	Data           interface{} `json:"data"`
	BlockIndex     int64       `json:"blockIndex"`
	Signature      string      `json:"signature"`
}

func (p *Pool) AddTransaction(tx Transaction) {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	p.Transactions = append(p.Transactions, tx)
}

func (p *Pool) CleanMemPool() {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	p.Transactions = nil
}

func (p *Pool) GetTransactions() []Transaction {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	return p.Transactions
}
