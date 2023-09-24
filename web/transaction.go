package web

import (
	"btc/utils/btc"
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
	req := new(createTransactionReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(400, errOutput{400, err.Error()})
	}

	if len(req.TxIn) == 0 {
		return c.JSON(400, errOutput{400, "invalid txin"})
	}

	var inputs btc.Inputs
	var outputs btc.Outputs

	for _, in := range req.TxIn {
		if in.TxId == "" {
			return c.JSON(400, errOutput{400, "invalid input item, txid is empty"})
		}
		inputs = append(inputs, btc.Input{
			TxId: in.TxId,
			VOut: in.VOut,
		})
	}
	for _, out := range req.PayToAddress {
		outputs = append(outputs, btc.Output{
			PayToAddress: out.Address,
			Amount:       out.Amount,
		})
	}

	transaction, err := btc.CreateRawTransaction(inputs, outputs)
	if err != nil {
		return c.JSON(400, errOutput{400, err.Error()})
	}

	return c.JSON(200, createTransactionOutput{200, transaction})
}

func SignTransaction(c echo.Context) error {
	req := new(signTransactionReq)

	if err := c.Bind(req); err != nil {
		return c.JSON(400, errOutput{400, err.Error()})
	}

	raw, err := btc.SignRawTransaction(req.Tx, req.Wif, req.RedeemScript)
	if err != nil {
		return c.JSON(400, errOutput{400, err.Error()})
	}

	return c.JSON(200, signTransactionOutput{200, raw})
}
