/*
 * Copyright (C) 2020 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package model

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	cosmos_types "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/joeqian10/neo-gogogo/helper"
	ontcommon "github.com/ontio/ontology/common"
	"github.com/shopspring/decimal"
	"math/big"
	"new_poly_explorer/DB"
	"new_poly_explorer/basedef"
	"strconv"
	"strings"
)

var (
	db = DB.Db
)

const (
	BTC_TOKEN_NAME = "btc"
	BTC_TOKEN_HASH = "0000000000000000000000000000000000000011"
)

type BigInt struct {
	big.Int
}

func NewBigIntFromInt(value int64) *BigInt {
	x := new(big.Int).SetInt64(value)
	return NewBigInt(x)
}

func NewBigInt(value *big.Int) *BigInt {
	return &BigInt{Int: *value}
}

func (bigInt *BigInt) Value() (driver.Value, error) {
	if bigInt == nil {
		return "null", nil
	}
	return bigInt.String(), nil
}

func (bigInt *BigInt) Scan(v interface{}) error {
	value, ok := v.([]byte)
	if !ok {
		return fmt.Errorf("type error, %v", v)
	}
	if string(value) == "null" {
		return nil
	}
	data, ok := new(big.Int).SetString(string(value), 10)
	if !ok {
		return fmt.Errorf("not a valid big integer: %s", value)
	}
	bigInt.Int = *data
	return nil
}

func HexString2Base58Address(address string) string {
	addr, err := ontcommon.AddressFromHexString(address)
	if err != nil {
		return ""
	}
	return addr.ToBase58()
}

func HexBytes2Base58Address(address []byte) string {
	addr, err := ontcommon.AddressParseFromBytes(address)
	if err != nil {
		return ""
	}
	return addr.ToBase58()
}

func String2Float64(value string) float64 {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return v
}

func HexReverse(arr []byte) []byte {
	l := len(arr)
	x := make([]byte, 0)
	for i := l - 1; i >= 0; i-- {
		x = append(x, arr[i])
	}
	return x
}

func HexStringReverse(value string) string {
	aa, _ := hex.DecodeString(value)
	bb := HexReverse(aa)
	return hex.EncodeToString(bb)
}

func AssetInfo(tokenHash string) (string, string) {
	token := new(Token)
	db.Where("hash = ?", tokenHash).Find(token)
	return token.Name, token.TokenType
}

func GetChains() []*Chain {
	// get all chains
	chains := make([]*Chain, 0)
	db.Find(&chains)
	if chains == nil {
		panic("GetExplorerInfo: can't get AllChainInfos")
	}
	return chains
}

func ChainId2Name(chainId uint64) string {
	chains := GetChains()
	for _, chain := range chains {
		if chain.ChainId == chainId {
			return chain.Name
		}
	}
	return "unknow chain"
}

func FormatFee(chain uint64, fee *BigInt) string {
	if chain == basedef.BTC_CROSSCHAIN_ID {
		precision_new := decimal.New(int64(100000000), 0)
		fee_new := decimal.New(fee.Int64(), 0)
		return fee_new.Div(precision_new).String()
	} else if chain == basedef.ONT_CROSSCHAIN_ID {
		precision_new := decimal.New(int64(1000000000), 0)
		fee_new := decimal.New(fee.Int64(), 0)
		return fee_new.Div(precision_new).String()
	} else if chain == basedef.ETHEREUM_CROSSCHAIN_ID {
		precision_new := decimal.New(int64(1000000000000000000), 0)
		fee_new := decimal.New(fee.Int64(), 0)
		return fee_new.Div(precision_new).String()
	} else {
		precision_new := decimal.New(int64(1), 0)
		fee_new := decimal.New(fee.Int64(), 0)
		return fee_new.Div(precision_new).String()
	}
}

func GetToken(tokenHash string, chainId uint64) *Token {
	token := new(Token)
	db.Where("hash = ? and chain_id = ?", tokenHash, chainId).First(token)
	return token
}

func FormatAmount(precision uint64, amount *BigInt) string {
	precision_new := decimal.New(int64(precision), 0)
	amount_new := decimal.New(amount.Int64(), 0)
	return amount_new.Div(precision_new).String()
}

func TxType2Name(ttype uint32) string {
	return "cross chain transfer"
}

func Hash2Address(chainId uint64, value string) string {
	if chainId == basedef.ETHEREUM_CROSSCHAIN_ID {
		addr := ethcommon.HexToAddress(value)
		return strings.ToLower(addr.String()[2:])
	} else if chainId == basedef.SWITCHEO_CROSSCHAIN_ID {
		addr, _ := cosmos_types.AccAddressFromHex(value)
		return addr.String()
	} else if chainId == basedef.BTC_CROSSCHAIN_ID {
		addrHex, _ := hex.DecodeString(value)
		return string(addrHex)
	} else if chainId == basedef.NEO_CROSSCHAIN_ID {
		addrHex, _ := hex.DecodeString(value)
		addr, _ := helper.UInt160FromBytes(addrHex)
		return helper.ScriptHashToAddress(addr)
	} else if chainId == basedef.ONT_CROSSCHAIN_ID {
		value = HexStringReverse(value)
		addr, _ := ontcommon.AddressFromHexString(value)
		return addr.ToBase58()
	}
	return value
}
