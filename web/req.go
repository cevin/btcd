package web

type parseAddressReq struct {
	PublicKeyHex  string `json:"public_key_hex" form:"public_key_hex" query:"public_key_hex"`
	PrivateKeyHex string `json:"private_key_hex" form:"private_key_hex" query:"private_key_hex"`
	WIF           string `json:"wif" form:"wif" query:"wif"`
	Script        string `json:"script" from:"script" query:"script"`
}

type createMultiSigReq struct {
	PublicKeyHexes string `json:"public_key_hexes" form:"public_key_hexes" query:"public_key_hexes"`
	Required       int    `json:"required" form:"required" query:"required"`
}

type parseTransactionReq struct {
	Tx string `json:"tx" form:"tx" query:"tx"`
}

type createTransactionReq struct {
	TxIn         []txin         `json:"txin"`
	PayToAddress []payToAddress `json:"pay_to_address"`
}
type txin struct {
	TxId string `json:"txid"`
	VOut int    `json:"vout,omitempty,-"`
}
type payToAddress struct {
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
}

type signTransactionReq struct {
	Tx           string `json:"tx" form:"tx" query:"tx"`
	Wif          string `json:"wif" form:"wif" query:"wif"`
	RedeemScript string `json:"redeem-script" form:"redeem-script" query:"redeem-script"`
}
