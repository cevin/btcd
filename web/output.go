package web

import "btc/utils/btc"

type errOutput struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type newAddressOutput struct {
	Code          int      `json:"code"`
	Address       string   `json:"address,omitempty,-"`
	Bech32Address string   `json:"bech32_address,omitempty,-"`
	WIF           string   `json:"wif,omitempty,-"`
	PrivateKey    string   `json:"private_key_hex,omitempty,-"`
	PublicKey     string   `json:"public_key_hex,omitempty,-"`
	Script        string   `json:"script,omitempty,-"`
	Asm           string   `json:"asm,omitempty,-"`
	Type          string   `json:"type,omitempty,-"`
	ReqSigs       uint32   `json:"reqSigs,omitempty,-"`
	Addresses     []string `json:"addresses,omitempty,-"`
}

type decodeTransactionOutput struct {
	Code        int             `json:"code"`
	Transaction btc.Transaction `json:"transaction"`
}

type createTransactionOutput struct {
	Code int    `json:"code"`
	Raw  string `json:"raw"`
}

type signTransactionOutput struct {
	Code int    `json:"code"`
	Raw  string `json:"raw"`
}
