package btc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"sort"
)

var mnet = &chaincfg.MainNetParams

// NewPrivateKey generate a Bitcoin wallet private key
func NewPrivateKey() *btcec.PrivateKey {
	key, _ := btcec.NewPrivateKey()

	return key
}

// NewAddress generate a normal wallet
func NewAddress() Address {
	privateKey := NewPrivateKey()
	wif, _ := btcutil.NewWIF(privateKey, mnet, true)

	addr, _ := ParseWIF(wif.String())

	return *addr

}

func ParseWIF(key string) (*Address, error) {
	wif, err := btcutil.DecodeWIF(key)
	if err != nil {
		return nil, err
	}

	compressedPubKey := wif.PrivKey.PubKey().SerializeCompressed()
	addressPubKey, err := btcutil.NewAddressPubKey(compressedPubKey, mnet)
	if err != nil {
		return nil, err
	}

	// get p2pkhPkScript
	addressPubKeyHash, _ := btcutil.NewAddressPubKeyHash(btcutil.Hash160(compressedPubKey), mnet)
	p2pkhPkScript, _ := txscript.PayToAddrScript(addressPubKeyHash)
	// get witnessPkScript
	addressWitnessPubKeyHash, _ := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(compressedPubKey), mnet)
	witnessPkScript, _ := txscript.PayToAddrScript(addressWitnessPubKeyHash)

	return &Address{
		Address:       addressPubKey.EncodeAddress(),
		Bech32Address: addressWitnessPubKeyHash.EncodeAddress(),
		PrivateKey:    wif.PrivKey,
		PublicKey:     wif.PrivKey.PubKey(),
		WIF:           wif.String(),

		AddressPubKey: addressPubKey,

		p2pkhPkScript:   p2pkhPkScript,
		witnessPkScript: witnessPkScript,
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

	pubKey, err := btcec.ParsePubKey(decodeString)
	if err != nil {
		return nil, err
	}

	witnessProg := btcutil.Hash160(pubKey.SerializeCompressed())
	addressWitnessPubKeyHash, _ := btcutil.NewAddressWitnessPubKeyHash(witnessProg, mnet)

	return &Address{
		Address:       addressPubKey.EncodeAddress(),
		Bech32Address: addressWitnessPubKeyHash.EncodeAddress(),
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
		Size:     mtx.SerializeSize(),
		VSize:    mempool.GetTxVirtualSize(btcutil.NewTx(&mtx)),
		Version:  mtx.Version,
		LockTime: mtx.LockTime,
		Vins:     createVinList(&mtx),
		Vouts:    createVoutList(&mtx, make(map[string]struct{})),
	}

	return &tx, nil
}

func CreateRawTransaction(inputs Inputs, outputs Outputs) (*wire.MsgTx, error) {

	// wire.TxVersion
	tx := wire.NewMsgTx(2)

	for _, input := range inputs {
		txHash, err := chainhash.NewHashFromStr(input.TxId)
		if err != nil {
			return nil, err
		}

		prevOut := wire.NewOutPoint(txHash, uint32(input.VOut))
		txIn := wire.NewTxIn(prevOut, nil, nil)
		txIn.Sequence = wire.MaxTxInSequenceNum - 2
		tx.AddTxIn(txIn)
	}

	for _, output := range outputs {
		sendAmount := output.Amount
		address, err := btcutil.DecodeAddress(output.PayToAddress, mnet)
		if err != nil {
			return nil, err
		}

		pkScript, err := txscript.PayToAddrScript(address)
		if err != nil {
			return nil, err
		}

		out := wire.NewTxOut(int64(sendAmount), pkScript)
		tx.AddTxOut(out)
	}

	return tx, nil

}

func CreateRawStringTransaction(inputs Inputs, outputs Outputs) (string, error) {
	transaction, err := CreateRawTransaction(inputs, outputs)
	if err != nil {
		return "", err
	}

	return serializeTx(transaction)
}

func SignRawStringTransaction(script string, inputs Inputs) (string, error) {
	b, err := hex.DecodeString(script)
	if err != nil {
		return "", nil
	}

	rawTx, err := btcutil.NewTxFromBytes(b)
	if err != nil {
		return "", err
	}

	tx := rawTx.MsgTx()

	return SignRawTxTransaction(tx, inputs)
}

func SignRawTxTransaction(tx *wire.MsgTx, inputs Inputs) (string, error) {

	errTmpl := "invalid input:[%d] : %s"

	hasSegWit := false
	hasZeroAmount := false

	prevOutFetcherMaps := make(map[wire.OutPoint]*wire.TxOut)
	prevOutFetcher := txscript.NewMultiPrevOutFetcher(prevOutFetcherMaps)

	signUseInputsMap := make(map[string]signInput)
	for idx, input := range inputs {
		addr, err := ParseWIF(input.WIF)
		if err != nil {
			return "", fmt.Errorf(errTmpl, idx, err)
		}

		if input.Amount == 0 {
			hasZeroAmount = true
		}

		if input.SegWit {
			if input.RedeemScript != "" {
				return "", fmt.Errorf(errTmpl, idx, "SegWit and RedeemScript cannot be set at the same time")
			}
			hasSegWit = true
		}

		if hasZeroAmount && hasSegWit {
			return "", fmt.Errorf("a segwit transaction was detected, but some inputs amount is zero")
		}

		signUseInputsMap[fmt.Sprintf("%s:%d", input.TxId, input.VOut)] = signInput{
			input: input,
			addr:  addr,
		}
	}

	if hasSegWit && len(signUseInputsMap) != len(tx.TxIn) {
		return "", fmt.Errorf("SegWit transaction was detected, but the number of UTXOs from the source transaction does not match the number of signed messages entered")
	}

	// calcTxSignHashes
	matched := 0
	var signHashes *txscript.TxSigHashes
	for idx, input := range tx.TxIn {
		txId := input.PreviousOutPoint.String()
		signUseInput, ok := signUseInputsMap[txId]
		if !ok {
			if hasSegWit {
				return "", fmt.Errorf("SegWit transaction was detected and %s private key is required", txId)
			}
			continue
		}

		matched++

		var pkScript []byte

		if signUseInput.input.SegWit {
			pkScript = signUseInput.addr.witnessPkScript
		} else {
			if signUseInput.input.RedeemScript != "" {
				decodeString, err := hex.DecodeString(signUseInput.input.RedeemScript)
				if err != nil {
					return "", fmt.Errorf(errTmpl, idx, err)
				}
				isMultiSignRedeemScript, _ := txscript.IsMultisigScript(decodeString)
				if !isMultiSignRedeemScript {
					return "", fmt.Errorf(errTmpl, idx, "invalid MultiSign redeem-script")
				}
				multiSignAddressPubKeyHash, err := btcutil.NewAddressScriptHash(decodeString, mnet)
				if err != nil {
					return "", fmt.Errorf(errTmpl, idx, err)
				}
				pkScript, _ = txscript.PayToAddrScript(multiSignAddressPubKeyHash)
				signUseInput.redeemScript = decodeString
			} else {
				pkScript = signUseInput.addr.p2pkhPkScript
			}
		}

		if pkScript == nil {
			return "", fmt.Errorf("calc pkScript failed for %s", txId)
		}

		amount := signUseInput.input.Amount

		signUseInputsMap[txId] = signUseInput
		prevOutFetcher.AddPrevOut(
			input.PreviousOutPoint,
			wire.NewTxOut(int64(amount), pkScript),
		)
	}

	if matched == 0 {
		return "", fmt.Errorf("no any input matched, the transaction was not signed")
	}

	if hasSegWit {
		signHashes = txscript.NewTxSigHashes(tx, prevOutFetcher)
	}

	for idx, input := range tx.TxIn {
		txId := input.PreviousOutPoint.String()
		signUseInput, ok := signUseInputsMap[txId]
		if !ok {
			continue
		}

		output := prevOutFetcher.FetchPrevOutput(input.PreviousOutPoint)

		errPrefix := fmt.Sprintf("invalid input[%d] %s : ", idx, txId)

		if !signUseInput.input.SegWit {

			pkScript := output.PkScript

			var redeemScript []byte
			if signUseInput.input.RedeemScript != "" {
				redeemScript = signUseInput.redeemScript
			}

			signScript, err := txscript.SignTxOutput(
				mnet,
				tx,
				idx,
				pkScript,
				txscript.SigHashAll,
				newLookupKeyFunc(signUseInput.addr.PrivateKey, mnet),
				newScriptDbFunc(redeemScript),
				tx.TxIn[idx].SignatureScript,
			)
			if err != nil {
				return "", fmt.Errorf("%s sign error : %s", errPrefix, err)
			}

			tx.TxIn[idx].SignatureScript = signScript
		} else {
			signHash, err := txscript.WitnessSignature(
				tx,
				signHashes,
				idx,
				output.Value,
				output.PkScript,
				txscript.SigHashAll,
				signUseInput.addr.PrivateKey,
				true,
			)
			if err != nil {
				return "", fmt.Errorf("%s sign error : %s", errPrefix, err)
			}
			tx.TxIn[idx].Witness = signHash
		}

	}

	return serializeTx(tx)
}
