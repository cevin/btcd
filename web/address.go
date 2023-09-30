package web

import (
	"btc/utils/btc"
	"encoding/hex"
	"github.com/labstack/echo/v4"
	"strings"
)

func NewAddress(c echo.Context) error {

	addr := btc.NewAddress()

	return c.JSON(200, newAddressOutput{
		Code:          200,
		Address:       addr.Address,
		Bech32Address: addr.Bech32Address,
		WIF:           addr.WIF,
		PrivateKey:    hex.EncodeToString(addr.PrivateKey.Serialize()),
		PublicKey:     hex.EncodeToString(addr.PublicKey.SerializeCompressed()),
	})
}

func NewMultiSigAddress(c echo.Context) error {
	req := new(createMultiSigReq)

	if err := c.Bind(req); err != nil {
		return err
	}

	pubKeys := strings.Split(req.PublicKeyHexes, ",")

	if len(pubKeys) < 2 {
		return c.JSON(400, errOutput{400, "at least two public key"})
	}
	if req.Required < 1 {
		return c.JSON(400, errOutput{400, "minimum signing private key at least one"})
	}

	address, err := btc.NewMultiSigAddress(pubKeys, req.Required)

	if err != nil {
		return c.JSON(400, errOutput{400, err.Error()})
	}

	return c.JSON(200, newAddressOutput{
		Code:      200,
		Type:      address.Type,
		Asm:       address.Asm,
		Addresses: address.Addresses,
		Address:   address.Address,
		ReqSigs:   uint32(address.ReqSigs),
		Script:    address.Script,
	})
}

func ParseAddress(c echo.Context) error {
	req := new(parseAddressReq)

	if err := c.Bind(req); err != nil {
		return err
	}

	if req.WIF == "" && req.PublicKeyHex == "" && req.PrivateKeyHex == "" && req.Script == "" {
		return c.JSON(400, errOutput{
			400,
			"wif,public_key_hex,private_key_hex,script one of the parameters must be passed.",
		})
	}

	if req.WIF != "" {
		if addr, err := btc.ParseWIF(req.WIF); err != nil {
			return c.JSON(400, errOutput{400, err.Error()})
		} else {
			return c.JSON(200, newAddressOutput{
				Code:          200,
				Address:       addr.Address,
				Bech32Address: addr.Bech32Address,
				WIF:           addr.WIF,
				PrivateKey:    hex.EncodeToString(addr.PrivateKey.Serialize()),
				PublicKey:     hex.EncodeToString(addr.PublicKey.SerializeCompressed()),
			})
		}
	} else if req.PublicKeyHex != "" {
		if addr, err := btc.ParsePublicKeyHex(req.PublicKeyHex); err != nil {
			return c.JSON(400, errOutput{400, err.Error()})
		} else {
			return c.JSON(200, newAddressOutput{
				Code:          200,
				Address:       addr.Address,
				Bech32Address: addr.Bech32Address,
				PublicKey:     hex.EncodeToString(addr.PublicKey.SerializeCompressed()),
			})
		}
	} else if req.Script != "" {
		if addr, err := btc.ParseMultiSigAddress(req.Script); err != nil {
			return c.JSON(400, errOutput{400, err.Error()})
		} else {
			return c.JSON(200, newAddressOutput{
				Code:      200,
				Type:      addr.Type,
				Asm:       addr.Asm,
				Addresses: addr.Addresses,
				Address:   addr.Address,
				ReqSigs:   uint32(addr.ReqSigs),
				Script:    addr.Script,
			})
		}
	}

	return c.JSON(400, errOutput{400, "invalid operation"})
}
