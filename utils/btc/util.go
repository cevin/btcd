package btc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func serializeTx(tx *wire.MsgTx) (string, error) {

	var txBuffer bytes.Buffer
	err := tx.Serialize(&txBuffer)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(txBuffer.Bytes()), nil
}

func newLookupKeyFunc(privateKey *btcec.PrivateKey, params *chaincfg.Params) txscript.KeyDB {
	return txscript.KeyClosure(func(addr btcutil.Address) (*btcec.PrivateKey, bool, error) {
		// For validate MultiSig script
		privateKeyAddr, err := btcutil.NewAddressPubKey(privateKey.PubKey().SerializeCompressed(), params)
		if err != nil {
			return nil, false, err
		}
		if privateKeyAddr.EncodeAddress() != addr.EncodeAddress() {
			return nil, false, fmt.Errorf("private key not match")
		}
		// The normal transaction can be directly returned to the private key.
		return privateKey, true, nil
	})
}

// newScriptDbFunc only work for MultiSig
func newScriptDbFunc(pkScript []byte) txscript.ScriptDB {
	return txscript.ScriptClosure(func(address btcutil.Address) ([]byte, error) {
		return pkScript, nil
	})
}

// -- copy from btcd

// createVinList returns a slice of JSON objects for the inputs of the passed
// transaction.
func createVinList(mtx *wire.MsgTx) []btcjson.Vin {
	vinList := make([]btcjson.Vin, len(mtx.TxIn))
	if blockchain.IsCoinBaseTx(mtx) {
		txIn := mtx.TxIn[0]
		vinList[0].Coinbase = hex.EncodeToString(txIn.SignatureScript)
		vinList[0].Sequence = txIn.Sequence
		vinList[0].Witness = witnessToHex(txIn.Witness)
		return vinList
	}

	for i, txIn := range mtx.TxIn {
		// The disassembled string will contain [error] inline
		// if the script doesn't fully parse, so ignore the
		// error here.
		disbuf, _ := txscript.DisasmString(txIn.SignatureScript)

		vinEntry := &vinList[i]
		vinEntry.Txid = txIn.PreviousOutPoint.Hash.String()
		vinEntry.Vout = txIn.PreviousOutPoint.Index
		vinEntry.Sequence = txIn.Sequence
		vinEntry.ScriptSig = &btcjson.ScriptSig{
			Asm: disbuf,
			Hex: hex.EncodeToString(txIn.SignatureScript),
		}

		if mtx.HasWitness() {
			vinEntry.Witness = witnessToHex(txIn.Witness)
		}
	}

	return vinList
}

// createVoutList returns a slice of JSON objects for the outputs of the passed
// transaction.
func createVoutList(mtx *wire.MsgTx, filterAddrMap map[string]struct{}) []btcjson.Vout {
	voutList := make([]btcjson.Vout, 0, len(mtx.TxOut))
	for i, v := range mtx.TxOut {
		// The disassembled string will contain [error] inline if the
		// script doesn't fully parse, so ignore the error here.
		disbuf, _ := txscript.DisasmString(v.PkScript)

		// Ignore the error here since an error means the script
		// couldn't parse and there is no additional information about
		// it anyways.
		scriptClass, addrs, reqSigs, _ := txscript.ExtractPkScriptAddrs(
			v.PkScript, mnet)

		// Encode the addresses while checking if the address passes the
		// filter when needed.
		passesFilter := len(filterAddrMap) == 0
		encodedAddrs := make([]string, len(addrs))
		for j, addr := range addrs {
			encodedAddr := addr.EncodeAddress()
			encodedAddrs[j] = encodedAddr

			// No need to check the map again if the filter already
			// passes.
			if passesFilter {
				continue
			}
			if _, exists := filterAddrMap[encodedAddr]; exists {
				passesFilter = true
			}
		}

		if !passesFilter {
			continue
		}

		var vout btcjson.Vout
		vout.N = uint32(i)
		vout.Value = btcutil.Amount(v.Value).ToBTC()
		vout.ScriptPubKey.Addresses = encodedAddrs
		vout.ScriptPubKey.Asm = disbuf
		vout.ScriptPubKey.Hex = hex.EncodeToString(v.PkScript)
		vout.ScriptPubKey.Type = scriptClass.String()
		vout.ScriptPubKey.ReqSigs = int32(reqSigs)

		voutList = append(voutList, vout)
	}

	return voutList
}

func witnessToHex(witness wire.TxWitness) []string {
	// Ensure nil is returned when there are no entries versus an empty
	// slice so it can properly be omitted as necessary.
	if len(witness) == 0 {
		return nil
	}

	result := make([]string, 0, len(witness))
	for _, wit := range witness {
		result = append(result, hex.EncodeToString(wit))
	}

	return result
}
