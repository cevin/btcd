package web

import (
	"btc/utils/btc"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"github.com/labstack/echo/v4"
)

func DecodeTransaction(c echo.Context) error {

	req := new(parseTransactionReq)

	if err := c.Bind(req); err != nil {
		return c.JSON(400, errOutput{400, err.Error()})
	}

	if req.Tx == "" {
		return c.JSON(400, errOutput{400, "invalid request"})
	}

	transaction, err := btc.ParseRawTransaction(req.Tx)
	if err != nil {
		return c.JSON(400, errOutput{400, err.Error()})
	}

	return c.JSON(200, decodeTransactionOutput{Code: 200, Transaction: *transaction})
}

func CreateRawTransaction(c echo.Context) error {
	req := new(createAndSignTransactionReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(400, errOutput{400, err.Error()})
	}

	inputs, outputs, err := performCreate(req)
	if err != nil {
		return c.JSON(400, errOutput{400, err.Error()})
	}

	transaction, err := btc.CreateRawStringTransaction(*inputs, *outputs)
	if err != nil {
		return c.JSON(400, errOutput{400, err.Error()})
	}

	return c.JSON(200, createTransactionOutput{200, transaction})
}

func SignTransaction(c echo.Context) error {
	req := new(createAndSignTransactionReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(400, errOutput{400, err.Error()})
	}

	if req.Raw == "" {
		return c.JSON(400, errOutput{400, "missing required parameter : raw"})
	}

	inputs, _, err := performCreate(req)
	if err != nil {
		return c.JSON(400, errOutput{400, err.Error()})
	}

	raw, err := btc.SignRawStringTransaction(req.Raw, *inputs)

	return c.JSON(200, signTransactionOutput{200, raw})
}

func CreateAndSignTransaction(c echo.Context) error {
	req := new(createAndSignTransactionReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(400, errOutput{400, err.Error()})
	}

	inputs, outputs, err := performCreate(req)
	if err != nil {
		return c.JSON(400, errOutput{400, err.Error()})
	}

	var unsignedTx *wire.MsgTx

	if req.Raw != "" {
		decodedTx, err := hex.DecodeString(req.Raw)
		if err != nil {
		}
		tx, err := btcutil.NewTxFromBytes(decodedTx)
		unsignedTx = tx.MsgTx()
	} else {
		unsignedTx, err = btc.CreateRawTransaction(*inputs, *outputs)
		if err != nil {
			return c.JSON(400, errOutput{400, err.Error()})
		}
	}

	signedTx, err := btc.SignRawTxTransaction(unsignedTx, *inputs)
	if err != nil {
		return c.JSON(400, errOutput{400, err.Error()})
	}

	return c.JSON(200, signTransactionOutput{200, signedTx})
}

func performCreate(req *createAndSignTransactionReq) (*btc.Inputs, *btc.Outputs, error) {

	inputs, err := convertReqToInput(req.TxIn)
	if err != nil {
		return nil, nil, err
	}
	outputs := convertReqToOutput(req.PayToAddresses)

	return inputs, &outputs, nil
}
