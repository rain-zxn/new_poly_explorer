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

// Package classification User API.
//
// The purpose of this service is to provide an application
// that is using plain go code to define an API
//
//      Host: localhost
//      Version: 0.0.1
//
// swagger:meta

package model

import (
	"strconv"
	"strings"
)

type ExplorerInfoReq struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type ErrorRsp struct {
	Message string
}

func MakeErrorRsp(messgae string) *ErrorRsp {
	errorRsp := &ErrorRsp{
		Message: messgae,
	}
	return errorRsp
}

type ExplorerInfoResp struct {
	Chains        []*ChainInfoResp       `json:"chains"`
	CrossTxNumber int64                  `json:"crosstxnumber"`
	Tokens        []*CrossChainTokenResp `json:"tokens"`
}

func getChainStatistic(chainId uint64, statistics []*ChainStatistic) *ChainStatistic {
	for _, statistic := range statistics {
		if statistic.ChainId == chainId {
			return statistic
		}
	}
	return nil
}

func MakeExplorerInfoResp(chains []*Chain, statistics []*ChainStatistic, tokenBasics []*TokenBasic) *ExplorerInfoResp {
	chainInfoResps := make([]*ChainInfoResp, 0)
	for _, chain := range chains {
		chainInfoResp := MakeChainInfoResp(chain)
		for _, statistic := range statistics {
			if statistic.ChainId == chain.ChainId {
				chainInfoResp.Addresses = statistic.Addresses
				chainInfoResp.In = statistic.In
				chainInfoResp.Out = statistic.Out
			}
		}
		for _, tokenBasic := range tokenBasics {
			for _, token := range tokenBasic.Tokens {
				if token.ChainId == chain.ChainId {
					chainInfoResp.Tokens = append(chainInfoResp.Tokens, MakeChainTokenResp(token))
				}
			}
		}
		chainInfoResps = append(chainInfoResps, chainInfoResp)
	}
	crossTxNumber := getChainStatistic(CHAIN_POLY, statistics).In
	crossChainTokenResp := make([]*CrossChainTokenResp, 0)
	for _, tokenBasic := range tokenBasics {
		crossChainTokenResp = append(crossChainTokenResp, MakeTokenBasicResp(tokenBasic))
	}
	explorerInfoResp := &ExplorerInfoResp{
		Chains:        chainInfoResps,
		CrossTxNumber: crossTxNumber,
		Tokens:        crossChainTokenResp,
	}
	return explorerInfoResp
}

type ChainInfoResp struct {
	Id     uint32 `json:"chainid"`
	Name   string `json:"chainname"`
	Height uint32 `json:"blockheight"`
	In     int64  `json:"in"`
	//InCrossChainTxStatus []*CrossChainTxStatus    `json:"incrosschaintxstatus"`
	Out int64 `json:"out"`
	//OutCrossChainTxStatus []*CrossChainTxStatus    `json:"outcrosschaintxstatus"`
	Addresses int64 `json:"addresses"`
	//Contracts []*ChainContractResp `json:"contracts"`
	Tokens []*ChainTokenResp `json:"tokens"`
}

func MakeChainInfoResp(chain *Chain) *ChainInfoResp {
	chainInfoResp := &ChainInfoResp{
		Id:        0,
		Name:      "",
		Height:    0,
		In:        0,
		Out:       0,
		Addresses: 0,
		Tokens:    nil,
	}
	return chainInfoResp
}

type CrossChainTxStatus struct {
	TT       uint32 `json:"timestamp"`
	TxNumber uint32 `json:"txnumber"`
}

type ChainContractResp struct {
	Id       uint32 `json:"chainid"`
	Contract string `json:"contract"`
}

type ChainTokenResp struct {
	Chain     int32  `json:"chainid"`
	ChainName string `json:"chainname"`
	Hash      string `json:"hash"`
	Token     string `json:"token"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Precision uint64 `json:"precision"`
	Desc      string `json:"desc"`
}

func MakeChainTokenResp(token *Token) *ChainTokenResp {
	chainTokenResp := &ChainTokenResp{}
	return chainTokenResp
}

type CrossChainTokenResp struct {
	Name   string            `json:"name"`
	Tokens []*ChainTokenResp `json:"tokens"`
}

func MakeTokenBasicResp(tokenBasic *TokenBasic) *CrossChainTokenResp {
	crossChainTokenResp := &CrossChainTokenResp{
		Name: tokenBasic.Name,
	}
	for _, token := range tokenBasic.Tokens {
		crossChainTokenResp.Tokens = append(crossChainTokenResp.Tokens, MakeChainTokenResp(token))
	}
	return crossChainTokenResp
}

type FChainTxResp struct {
	ChainId    uint32              `json:"chainid"`
	ChainName  string              `json:"chainname"`
	TxHash     string              `json:"txhash"`
	State      byte                `json:"state"`
	TT         uint32              `json:"timestamp"`
	Fee        string              `json:"fee"`
	Height     uint32              `json:"blockheight"`
	User       string              `json:"user"`
	TChainId   uint32              `json:"tchainid"`
	TChainName string              `json:"tchainname"`
	Contract   string              `json:"contract"`
	Key        string              `json:"key"`
	Param      string              `json:"param"`
	Transfer   *FChainTransferResp `json:"transfer"`
}

func makeFChainTxResp(fChainTx *SrcTransaction, token *Token) *FChainTxResp {
	fChainTxResp := &FChainTxResp{
		ChainId:    uint32(fChainTx.ChainId),
		ChainName:  ChainId2Name(fChainTx.ChainId),
		TxHash:     fChainTx.Hash,
		State:      byte(fChainTx.State),
		TT:         uint32(fChainTx.Time),
		Fee:        FormatFee(fChainTx.ChainId, fChainTx.Fee),
		Height:     uint32(fChainTx.Height),
		User:       Hash2Address(fChainTx.ChainId, fChainTx.User),
		TChainId:   uint32(fChainTx.DstChainId),
		TChainName: ChainId2Name(fChainTx.DstChainId),
		Contract:   fChainTx.Contract,
		Key:        fChainTx.Key,
		Param:      fChainTx.Param,
	}
	if fChainTx.SrcTransfer != nil {
		fChainTxResp.Transfer = &FChainTransferResp{
			From:        Hash2Address(fChainTx.SrcTransfer.ChainId, fChainTx.SrcTransfer.From),
			To:          Hash2Address(fChainTx.SrcTransfer.DstChainId, fChainTx.SrcTransfer.To),
			Amount:      strconv.FormatUint(fChainTx.SrcTransfer.Amount.Uint64(), 10),
			ToChain:     uint32(fChainTx.SrcTransfer.DstChainId),
			ToChainName: ChainId2Name(fChainTx.SrcTransfer.DstChainId),
			ToUser:      Hash2Address(fChainTx.SrcTransfer.DstChainId, fChainTx.SrcTransfer.DstUser),
		}
		fChainTxResp.Transfer.TokenHash = fChainTx.SrcTransfer.Asset
		if token != nil {
			fChainTxResp.Transfer.TokenHash = token.Hash
			fChainTxResp.Transfer.TokenName = token.Name
			fChainTxResp.Transfer.TokenType = token.TokenType
			fChainTxResp.Transfer.Amount = FormatAmount(token.Precision, fChainTx.SrcTransfer.Amount)
		}
		totoken := GetToken(fChainTx.SrcTransfer.DstAsset, fChainTx.SrcTransfer.DstChainId)
		fChainTxResp.Transfer.ToTokenHash = fChainTx.SrcTransfer.DstAsset
		if totoken != nil {
			fChainTxResp.Transfer.ToTokenHash = totoken.Hash
			fChainTxResp.Transfer.ToTokenName = totoken.Name
			fChainTxResp.Transfer.ToTokenType = totoken.TokenType
		}
	}
	if fChainTx.ChainId == CHAIN_ETH {
		fChainTxResp.TxHash = "0x" + fChainTx.Key
	} else if fChainTx.ChainId == CHAIN_SWITCHEO {
		fChainTxResp.TxHash = strings.ToUpper(fChainTxResp.TxHash)
	}
	return fChainTxResp
}

type FChainTransferResp struct {
	TokenHash   string `json:"tokenhash"`
	TokenName   string `json:"tokenname"`
	TokenType   string `json:"tokentype"`
	From        string `json:"from"`
	To          string `json:"to"`
	Amount      string `json:"amount"`
	ToChain     uint32 `json:"tchainid"`
	ToChainName string `json:"tchainname"`
	ToTokenHash string `json:"totokenhash"`
	ToTokenName string `json:"totokenname"`
	ToTokenType string `json:"totokentype"`
	ToUser      string `json:"tuser"`
}

type MChainTxResp struct {
	ChainId    uint32 `json:"chainid"`
	ChainName  string `json:"chainname"`
	TxHash     string `json:"txhash"`
	State      byte   `json:"state"`
	TT         uint32 `json:"timestamp"`
	Fee        string `json:"fee"`
	Height     uint32 `json:"blockheight"`
	FChainId   uint32 `json:"fchainid"`
	FChainName string `json:"fchainname"`
	FTxHash    string `json:"ftxhash"`
	TChainId   uint32 `json:"tchainid"`
	TChainName string `json:"tchainname"`
	Key        string `json:"key"`
}

func makeMChainTxResp(mChainTx *PolyTransaction) *MChainTxResp {
	mChainTxResp := &MChainTxResp{
		ChainId:    uint32(mChainTx.ChainId),
		ChainName:  ChainId2Name(mChainTx.ChainId),
		TxHash:     mChainTx.Hash,
		State:      byte(mChainTx.State),
		TT:         uint32(mChainTx.Time),
		Fee:        FormatFee(mChainTx.ChainId, mChainTx.Fee),
		Height:     uint32(mChainTx.Height),
		FChainId:   uint32(mChainTx.SrcChainId),
		FChainName: ChainId2Name(mChainTx.SrcChainId),
		FTxHash:    mChainTx.SrcHash,
		TChainId:   uint32(mChainTx.DstChainId),
		TChainName: ChainId2Name(mChainTx.DstChainId),
		Key:        mChainTx.Key,
	}
	return mChainTxResp
}

type TChainTxResp struct {
	ChainId    uint32              `json:"chainid"`
	ChainName  string              `json:"chainname"`
	TxHash     string              `json:"txhash"`
	State      byte                `json:"state"`
	TT         uint32              `json:"timestamp"`
	Fee        string              `json:"fee"`
	Height     uint32              `json:"blockheight"`
	FChainId   uint32              `json:"fchainid"`
	FChainName string              `json:"fchainname"`
	Contract   string              `json:"contract"`
	RTxHash    string              `json:"mtxhash"`
	Transfer   *TChainTransferResp `json:"transfer"`
}

func makeTChainTxResp(tChainTx *DstTransaction) *TChainTxResp {
	tChainTxResp := &TChainTxResp{
		ChainId:    uint32(tChainTx.ChainId),
		ChainName:  ChainId2Name(tChainTx.ChainId),
		TxHash:     tChainTx.Hash,
		State:      byte(tChainTx.State),
		TT:         uint32(tChainTx.Time),
		Fee:        FormatFee(tChainTx.ChainId, tChainTx.Fee),
		Height:     uint32(tChainTx.Height),
		FChainId:   uint32(tChainTx.SrcChainId),
		FChainName: ChainId2Name(tChainTx.SrcChainId),
		Contract:   tChainTx.Contract,
		RTxHash:    tChainTx.PolyHash,
	}
	if tChainTx.DstTransfer != nil {
		tChainTxResp.Transfer = &TChainTransferResp{
			From:   tChainTx.DstTransfer.From,
			To:     tChainTx.DstTransfer.To,
			Amount: strconv.FormatUint(tChainTx.DstTransfer.Amount.Uint64(), 10),
		}
		token := GetToken(tChainTx.DstTransfer.Asset, tChainTx.DstTransfer.ChainId)
		tChainTxResp.Transfer.TokenHash = tChainTx.DstTransfer.Asset
		if token != nil {
			tChainTxResp.Transfer.TokenHash = token.Hash
			tChainTxResp.Transfer.TokenName = token.Name
			tChainTxResp.Transfer.TokenType = token.TokenType
			tChainTxResp.Transfer.Amount = FormatAmount(token.Precision, tChainTx.DstTransfer.Amount)
		}
	}
	if tChainTx.ChainId == CHAIN_ETH {
		tChainTxResp.TxHash = "0x" + tChainTxResp.TxHash
	} else if tChainTx.ChainId == CHAIN_SWITCHEO {
		tChainTxResp.TxHash = strings.ToUpper(tChainTxResp.TxHash)
	}
	return tChainTxResp
}

type TChainTransferResp struct {
	TokenHash string `json:"tokenhash"`
	TokenName string `json:"tokenname"`
	TokenType string `json:"tokentype"`
	From      string `json:"from"`
	To        string `json:"to"`
	Amount    string `json:"amount"`
}

type CrossTransferResp struct {
	CrossTxType uint32 `json:"crosstxtype"`
	CrossTxName string `json:"crosstxname"`
	FromChainId uint32 `json:"fromchainid"`
	FromChain   string `json:"fromchainname"`
	FromAddress string `json:"fromaddress"`
	ToChainId   uint32 `json:"tochainid"`
	ToChain     string `json:"tochainname"`
	ToAddress   string `json:"toaddress"`
	TokenHash   string `json:"tokenhash"`
	TokenName   string `json:"tokenname"`
	TokenType   string `json:"tokentype"`
	Amount      string `json:"amount"`
}

func makeCrossTransfer(chainid uint64, user string, transfer *SrcTransfer) *CrossTransferResp {
	if transfer == nil {
		return nil
	}
	crossTransfer := new(CrossTransferResp)
	crossTransfer.CrossTxType = 1
	crossTransfer.CrossTxName = TxType2Name(crossTransfer.CrossTxType)
	crossTransfer.FromChainId = uint32(chainid)
	crossTransfer.FromChain = ChainId2Name(uint64(crossTransfer.FromChainId))
	crossTransfer.FromAddress = Hash2Address(chainid, user)
	crossTransfer.ToChainId = uint32(transfer.DstChainId)
	crossTransfer.ToChain = ChainId2Name(uint64(crossTransfer.ToChainId))
	crossTransfer.ToAddress = Hash2Address(transfer.DstChainId, transfer.DstUser)
	token := GetToken(transfer.Asset, transfer.ChainId)
	if token != nil {
		crossTransfer.TokenHash = token.Hash
		crossTransfer.TokenName = token.Name
		crossTransfer.TokenType = token.TokenType
		crossTransfer.Amount = FormatAmount(token.Precision, transfer.Amount)
	}
	return crossTransfer
}

// swagger:parameters CrossTxReq
type CrossTxReq struct {
	// in: query
	TxHash string `json:"txhash"`
}

type CrossTxResp struct {
	Transfer       *CrossTransferResp `json:"crosstransfer"`
	Fchaintx       *FChainTxResp      `json:"fchaintx"`
	Fchaintx_valid bool               `json:"fchaintx_valid"`
	Mchaintx       *MChainTxResp      `json:"mchaintx"`
	Mchaintx_valid bool               `json:"mchaintx_valid"`
	Tchaintx       *TChainTxResp      `json:"tchaintx"`
	Tchaintx_valid bool               `json:"tchaintx_valid"`
}

func MakeCrossTxResp(srcPolyDst []*SrcPolyDstRelation) *CrossTxResp {
	crosstx := &CrossTxResp{
		Fchaintx_valid: false,
		Mchaintx_valid: false,
		Tchaintx_valid: false,
		Transfer: &CrossTransferResp{
			CrossTxType: 0,
		},
	}
	tx := srcPolyDst[0]

	if tx.SrcTransaction != nil {
		crosstx.Fchaintx_valid = true
		crosstx.Fchaintx = makeFChainTxResp(tx.SrcTransaction, tx.Token)
		crosstx.Transfer = makeCrossTransfer(tx.SrcTransaction.ChainId, tx.SrcTransaction.User, tx.SrcTransaction.SrcTransfer)
	}
	if tx.PolyTransaction != nil && crosstx.Fchaintx_valid == true {
		crosstx.Mchaintx_valid = true
		crosstx.Mchaintx = makeMChainTxResp(tx.PolyTransaction)
	}
	if tx.DstTransaction != nil && crosstx.Mchaintx_valid == true {
		crosstx.Tchaintx_valid = true
		crosstx.Tchaintx = makeTChainTxResp(tx.DstTransaction)
	}
	return crosstx
}

type CrossTxListReq struct {
	PageSize int
	PageNo   int
}

type CrossTxOutlineResp struct {
	TxHash     string `json:"txhash"`
	State      byte   `json:"state"`
	TT         uint32 `json:"timestamp"`
	Fee        uint64 `json:"fee"`
	Height     uint32 `json:"blockheight"`
	FChainId   uint32 `json:"fchainid"`
	FChainName string `json:"fchainname"`
	TChainId   uint32 `json:"tchainid"`
	TChainName string `json:"tchainname"`
}

type CrossTxListResp struct {
	CrossTxList []*CrossTxOutlineResp `json:"crosstxs"`
}

func MakeCrossTxListResp(txs []*SrcPolyDstRelation) *CrossTxListResp {
	crossTxListResp := &CrossTxListResp{}
	crossTxListResp.CrossTxList = make([]*CrossTxOutlineResp, 0)
	for _, tx := range txs {
		crossTxListResp.CrossTxList = append(crossTxListResp.CrossTxList, &CrossTxOutlineResp{
			TxHash: tx.PolyHash,
		})
	}
	return crossTxListResp
}

type TokenTxListReq struct {
	PageSize int
	PageNo   int
	ChainId  uint64
	Token    string `json:"token"`
}

type TokenTxResp struct {
	TxHash string `json:"txhash"`
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
	TT     uint32 `json:"timestamp"`
	Height uint32 `json:"blockheight"`
	Direct uint32 `json:"direct"`
}

type TokenTxListResp struct {
	TokenTxList []*TokenTxResp `json:"tokentxs"`
	Total       int64          `json:"total"`
}

func MakeTokenTxList(transactoins []*TransactionOnToken, tokenStatistic *TokenStatistic) *TokenTxListResp {
	tokenTxListResp := &TokenTxListResp{}
	tokenTxListResp.Total = tokenStatistic.InCounter + tokenStatistic.OutCounter
	tokenTxListResp.TokenTxList = make([]*TokenTxResp, 0)
	for _, transactoin := range transactoins {
		tokenTxListResp.TokenTxList = append(tokenTxListResp.TokenTxList, &TokenTxResp{
			TxHash: transactoin.Hash,
		})
	}
	return tokenTxListResp
}

type AddressTxListReq struct {
	PageSize int
	PageNo   int
	Address  string `json:"address"`
	ChainId  string `json:"chain"`
}

type AddressTxResp struct {
	TxHash    string `json:"txhash"`
	From      string `json:"from"`
	To        string `json:"to"`
	Amount    string `json:"amount"`
	TT        uint32 `json:"timestamp"`
	Height    uint32 `json:"blockheight"`
	TokenHash string `json:"tokenhash"`
	TokenName string `json:"tokenname"`
	TokenType string `json:"tokentype"`
	Direct    uint32 `json:"direct"`
}

type AddressTxListResp struct {
	AddressTxList []*AddressTxResp `json:"addresstxs"`
	Total         int64            `json:"total"`
}

func MakeAddressTxList(transactoins []*TransactionOnAddress, counter int64) *AddressTxListResp {
	addressTxListResp := &AddressTxListResp{}
	addressTxListResp.Total = counter
	addressTxListResp.AddressTxList = make([]*AddressTxResp, 0)
	for _, transactoin := range transactoins {
		addressTxListResp.AddressTxList = append(addressTxListResp.AddressTxList, &AddressTxResp{
			TxHash: transactoin.Hash,
		})
	}
	return addressTxListResp
}
