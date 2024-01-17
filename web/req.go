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

type createAndSignTransactionReq struct {
	Raw            string         `json:"raw"`
	TxIn           []txin         `json:"txin"`
	PayToAddresses []payToAddress `json:"pay_to_addresses"`
}
type txin struct {
	TxId         string `json:"txid"`
	VOut         *int   `json:"vout,omitempty,-"`
	WIF          string `json:"wif,omitempty,-"`
	RedeemScript string `json:"redeem-script,omitempty,-"`
	SegWit       bool   `json:"segwit,omitempty,-"`
	Amount       int64  `json:"amount,omitempty,-"`
}
type payToAddress struct {
	Address string `json:"address"`
	Amount  int64  `json:"amount"`
}
