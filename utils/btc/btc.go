package btc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"sort"
)

var mnet = &chaincfg.MainNetParams

type Address struct {
	Address       string
	PrivateKey    *btcec.PrivateKey
	PublicKey     *btcec.PublicKey
	AddressPubKey *btcutil.AddressPubKey
	WIF           string
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
	TxId string
	VOut int
}

type Output struct {
	PayToAddress string
	Amount       float64
}

type Inputs []Input
type Outputs []Output

// NewPrivateKey generate a Bitcoin wallet private key
func NewPrivateKey() *btcec.PrivateKey {
	key, _ := btcec.NewPrivateKey()

	return key
}

// NewAddress generate a normal wallet
func NewAddress() Address {
	privateKey := NewPrivateKey()
	wif, _ := btcutil.NewWIF(privateKey, &chaincfg.MainNetParams, true)

	addressPubKey, _ := btcutil.NewAddressPubKey(privateKey.PubKey().SerializeCompressed(), mnet)

	return Address{
		Address:       addressPubKey.EncodeAddress(),
		PrivateKey:    privateKey,
		PublicKey:     privateKey.PubKey(),
		WIF:           wif.String(),
		AddressPubKey: addressPubKey,
	}

}

func ParseWIF(key string) (*Address, error) {
	wif, err := btcutil.DecodeWIF(key)
	if err != nil {
		return nil, err
	}

	addressPubKey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), mnet)
	if err != nil {
		return nil, err
	}

	return &Address{
		Address:       addressPubKey.EncodeAddress(),
		PrivateKey:    wif.PrivKey,
		PublicKey:     wif.PrivKey.PubKey(),
		WIF:           wif.String(),
		AddressPubKey: addressPubKey,
	}, nil
}

func ParsePublicKeyHex(key string) (*Address, error) {
	decodeString, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}

	addressPubKey, err := btcutil.NewAddressPubKey(decodeString, mnet)
	if err != nil {
		return nil, err
	}

	return &Address{
		Address:       addressPubKey.EncodeAddress(),
		PublicKey:     addressPubKey.PubKey(),
		AddressPubKey: addressPubKey,
	}, nil
}

func NewMultiSigAddress(pubKeys []string, required int) (*MultiSigAddress, error) {
	var usedKeys []*btcutil.AddressPubKey
	addresses := make([]string, len(pubKeys))

	sort.Strings(pubKeys)

	for i, pubKey := range pubKeys {
		address, err := ParsePublicKeyHex(pubKey)
		if err != nil {
			return nil, err
		}
		usedKeys = append(usedKeys, address.AddressPubKey)
		addresses[i] = address.Address
	}

	script, err := txscript.MultiSigScript(usedKeys, required)
	if err != nil {
		return nil, err
	}
	addressHash, err := btcutil.NewAddressScriptHashFromHash(btcutil.Hash160(script), mnet)
	if err != nil {
		return nil, err
	}

	disBuff, _ := txscript.DisasmString(script)

	return &MultiSigAddress{
		Asm:       disBuff,
		Type:      "multisig",
		Addresses: addresses,
		Address:   addressHash.EncodeAddress(),
		ReqSigs:   int32(required),
		Script:    hex.EncodeToString(script),
	}, nil
}

func ParseMultiSigAddress(script string) (*MultiSigAddress, error) {
	b, err := hex.DecodeString(script)
	if err != nil {
		return nil, err
	}

	addressHash, err := btcutil.NewAddressScriptHashFromHash(btcutil.Hash160(b), mnet)
	if err != nil {
		return nil, err
	}

	disBuff, _ := txscript.DisasmString(b)
	scriptClass, addrs, reqSigs, err := txscript.ExtractPkScriptAddrs(b, mnet)
	if err != nil {
		return nil, err
	}
	addresses := make([]string, len(addrs))
	for i, addr := range addrs {
		addresses[i] = addr.EncodeAddress()
	}

	return &MultiSigAddress{
		Asm:       disBuff,
		Type:      scriptClass.String(),
		Address:   addressHash.EncodeAddress(),
		Addresses: addresses,
		ReqSigs:   int32(reqSigs),
		Script:    script,
	}, nil
}

func ParseRawTransaction(str string) (*Transaction, error) {
	decodeString, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}

	var mtx wire.MsgTx
	err = mtx.Deserialize(bytes.NewReader(decodeString))
	if err != nil {
		return nil, err
	}

	tx := Transaction{
		Txid:     mtx.TxHash().String(),
		Hash:     mtx.TxHash().String(),
		Size:     mtx.SerializeSize(),
		VSize:    mempool.GetTxVirtualSize(btcutil.NewTx(&mtx)),
		Version:  mtx.Version,
		LockTime: mtx.LockTime,
		Vins:     createVinList(&mtx),
		Vouts:    createVoutList(&mtx, make(map[string]struct{})),
	}

	return &tx, nil
}

func CreateRawTransaction(inputs Inputs, outputs Outputs) (string, error) {

	tx := wire.NewMsgTx(wire.TxVersion)

	for _, input := range inputs {
		txHash, err := chainhash.NewHashFromStr(input.TxId)
		if err != nil {
			return "", err
		}

		prevOut := wire.NewOutPoint(txHash, uint32(input.VOut))
		txIn := wire.NewTxIn(prevOut, nil, nil)
		tx.AddTxIn(txIn)
	}

	for _, output := range outputs {

		sendAmount, err := btcutil.NewAmount(output.Amount)
		if err != nil {
			return "", err
		}
		address, err := btcutil.DecodeAddress(output.PayToAddress, mnet)
		if err != nil {
			return "", err
		}

		pkScript, err := txscript.PayToAddrScript(address)
		if err != nil {
			return "", err
		}

		out := wire.NewTxOut(int64(sendAmount), pkScript)
		tx.AddTxOut(out)
	}

	return serializeTx(tx)

}

func SignRawTransaction(script string, key string, redeemScript string) (string, error) {

	b, err := hex.DecodeString(script)
	if err != nil {
		return "", nil
	}

	wif, err := btcutil.DecodeWIF(key)
	if err != nil {
		return "", err
	}

	var pkScript, scriptPkScript []byte
	if redeemScript == "" {
		// qus
		address, err := btcutil.NewAddressPubKeyHash(
			btcutil.Hash160(wif.PrivKey.PubKey().SerializeCompressed()),
			mnet,
		)
		if err != nil {
			return "", fmt.Errorf("parseing address pubkey error, %s", err)
		}
		scriptPkScript, err = txscript.PayToAddrScript(address)
		if err != nil {
			return "", err
		}
		pkScript = nil
	} else {
		pkScript, err = hex.DecodeString(redeemScript)
		if err != nil {
			return "", err
		}
		ok, _ := txscript.IsMultisigScript(pkScript)
		if !ok {
			return "", fmt.Errorf("invalid redeem script")
		}
		scriptAddr, err := btcutil.NewAddressScriptHash(pkScript, mnet)
		if err != nil {
			return "", err
		}
		scriptPkScript, err = txscript.PayToAddrScript(scriptAddr)
		if err != nil {
			return "", err
		}
	}

	rawTx, err := btcutil.NewTxFromBytes(b)
	if err != nil {
		return "", err
	}

	tx := rawTx.MsgTx()

	for i := 0; i < len(tx.TxIn); i++ {
		var signScript []byte
		signScript, err = txscript.SignTxOutput(
			mnet,
			tx,
			i,
			scriptPkScript,
			txscript.SigHashAll,
			newLookupKeyFunc(wif.PrivKey, mnet),
			newScriptDbFunc(pkScript),
			tx.TxIn[i].SignatureScript,
		)
		if err != nil {
			return "", nil
		}

		tx.TxIn[i].SignatureScript = signScript
	}

	signedTx, _ := serializeTx(tx)

	return signedTx, nil
}
