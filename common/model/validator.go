package model

type Validator struct {
	ValidatorAddress string `json:"validator_address"`
	ValidatorPubKey  string `json:"validator_pub_key"`
	Stake            int64  `json:"stake"`
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

type ValidatorMetaData struct {
	Address               string `json:"address"`
	SolanaAddress         string `json:"SolanaAddress"`
	Stake                 uint64 `json:"Stake"`
	JoinTime              string `json:"JoinTime"`
	InNetwork             bool   `json:"InNetwork"`
	CompletedVerification int    `json:"CompletedVerification"`
}

type WorkLoad struct {
	TransactionsAmount int  `json:"TransactionsAmount"`
	CanHandle          bool `json:"CanHandle"`
}
