package model

type Message struct {
	TransactionRequests Transaction `json:"transactionRequests"`
	PublicKey           string      `json:"publicKey"`
}

type TransactionRequest struct {
	ActionTaker    string      `json:"actionTaker"`
	ActionReceiver string      `json:"actionReceiver"`
	Data           interface{} `json:"data"`
	Signature      string      `json:"signature"`
}

type PublicKeyJSON struct {
	X string `json:"x"`
	Y string `json:"y"`
}
