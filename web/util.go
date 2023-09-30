package web

import (
	"btc/utils/btc"
	"fmt"
)

func convertReqToInput(reqInputs []txin) (*btc.Inputs, error) {

	if len(reqInputs) == 0 {
		return nil, fmt.Errorf("empty txin")
	}

	var inputs btc.Inputs
	for idx, reqInput := range reqInputs {
		if reqInput.TxId == "" {
			return nil, inputErr(idx, "txid is empty")
		}
		if reqInput.VOut == nil {
			return nil, inputErr(idx, "vout required")
		}

		if reqInput.RedeemScript != "" && reqInput.SegWit {
			return nil, inputErr(idx, "SegWit and RedeemScript cannot be set at the same time")
		}

		if reqInput.SegWit && reqInput.Amount == 0 {
			return nil, inputErr(idx, "amount is required when the input is SegWit transaction")
		}

		inputs = append(inputs, btc.Input{
			TxId:         reqInput.TxId,
			VOut:         *reqInput.VOut,
			WIF:          reqInput.WIF,
			RedeemScript: reqInput.RedeemScript,
			SegWit:       reqInput.SegWit,
			Amount:       reqInput.Amount,
		})
	}

	return &inputs, nil
}

func inputErr(idx int, err string) error {
	return fmt.Errorf("invalid input[%d], %s", idx, err)
}

func convertReqToOutput(reqOutputs []payToAddress) btc.Outputs {

	var outputs btc.Outputs
	for _, out := range reqOutputs {
		outputs = append(outputs, btc.Output{
			PayToAddress: out.Address,
			Amount:       out.Amount,
		})
	}

	return outputs
}
