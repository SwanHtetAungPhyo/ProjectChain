package model

import "sync"

type Pool struct {
	Transactions []Transaction `json:"transactions"`
	mu           sync.Mutex
}

func (p *Pool) AddTransaction(tx Transaction) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Transactions = append(p.Transactions, tx)
}

func (p *Pool) CleanMemPool() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Transactions = nil
}

func (p *Pool) GetTransactions() []Transaction {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Transactions
}
