package btc

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
)

type Address struct {
	Address       string
	Bech32Address string
	PrivateKey    *btcec.PrivateKey
	PublicKey     *btcec.PublicKey
	AddressPubKey *btcutil.AddressPubKey
	WIF           string

	p2pkhPkScript   []byte
	witnessPkScript []byte
}

type MultiSigAddress struct {
	Asm       string   `json:"asm,omitempty,-"`
	Type      string   `json:"type,omitempty,-"`
	ReqSigs   int32    `json:"reqSigs,omitempty,-"`
	Address   string   `json:"address,omitempty,-"`
	Addresses []string `json:"addresses,omitempty,-"`
	Script    string   `json:"script,omitempty,-"`
}

type Transaction struct {
	Txid     string         `json:"txid"`
	Hash     string         `json:"hash"`
	Size     int            `json:"size"`
	VSize    int64          `json:"vsize"`
	Weight   int            `json:"weight"`
	Version  int32          `json:"version"`
	LockTime uint32         `json:"locktime"`
	Vins     []btcjson.Vin  `json:"vin"`
	Vouts    []btcjson.Vout `json:"vout"`
}

type Input struct {
	TxId         string
	VOut         int
	WIF          string
	RedeemScript string
	SegWit       bool
	Amount       float64
}

type Output struct {
	PayToAddress string
	Amount       float64
}

type Inputs []Input
type Outputs []Output

type signInput struct {
	input        Input
	addr         *Address
	redeemScript []byte
}
