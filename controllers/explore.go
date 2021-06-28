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

package controllers

import (
	"encoding/json"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"new_poly_explorer/DB"
	"new_poly_explorer/model"
)

var (
	db = DB.Db
)

type ExplorerController struct {
	beego.Controller
}

// GetExplorerInfo shows explorer information, such as current blockheight (the number of blockchain and so on) on the home page.
func (c *ExplorerController) GetExplorerInfo() {
	// get parameter
	var explorerReq model.ExplorerInfoReq
	var err error
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &explorerReq); err != nil {
		c.Data["json"] = model.MakeErrorRsp(fmt.Sprintf("request parameter is invalid!"))
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
	}

	//get all chains
	chains := make([]*model.Chain, 0)
	res := db.Find(&chains)
	if res.RowsAffected == 0 {
		c.Data["json"] = model.MakeErrorRsp(fmt.Sprintf("chain does not exist"))
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
		return
	}

	// get all chains statistic
	chainStatistics := make([]*model.ChainStatistic, 0)

	// get all tokens
	tokenBasics := make([]*model.TokenBasic, 0)
	res = db.Find(&tokenBasics)
	if res.RowsAffected == 0 {
		c.Data["json"] = model.MakeErrorRsp(fmt.Sprintf("chain does not exist"))
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
		return
	}

	c.Data["json"] = model.MakeExplorerInfoResp(chains, chainStatistics, tokenBasics)
	c.ServeJSON()
}

func (c *ExplorerController) GetTokenTxList() {
	// get parameter
	var tokenTxListReq model.TokenTxListReq
	var err error
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &tokenTxListReq); err != nil {
		c.Data["json"] = model.MakeErrorRsp(fmt.Sprintf("request parameter is invalid!"))
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
	}

	//
	transactionOnTokens := make([]*model.TransactionOnToken, 0)
	db.Raw("select a.hash, a.height, a.time, a.chain_id, b.from, b.to, b.amount, 1 as direct from src_transactions a inner join src_transfers b on a.hash = b.tx_hash where b.asset = ? and b.chain_id = ?"+
		"union select c.hash, c.height, c.time, c.chain_id, d.from, d.to, d.amount, 2 as direct from dst_transactions c inner join dst_transfers d on c.hash = d.tx_hash where d.asset = ? and d.chain_id = ?"+
		"order by height desc limit ?,?",
		tokenTxListReq.Token, tokenTxListReq.ChainId, tokenTxListReq.Token, tokenTxListReq.ChainId, tokenTxListReq.PageSize*tokenTxListReq.PageNo, tokenTxListReq.PageSize).Find(&transactionOnTokens)
	//
	tokenStatistic := new(model.TokenStatistic)
	db.Where("chain_id = ? and hash = ?", tokenTxListReq.ChainId, tokenTxListReq.Token).Find(tokenStatistic)
	//
	c.Data["json"] = model.MakeTokenTxList(transactionOnTokens, tokenStatistic)
	c.ServeJSON()
}

func (c *ExplorerController) GetAddressTxList() {
	// get parameter
	var addressTxListReq model.AddressTxListReq
	var err error
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &addressTxListReq); err != nil {
		c.Data["json"] = model.MakeErrorRsp(fmt.Sprintf("request parameter is invalid!"))
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
	}

	//
	transactionOnAddresses := make([]*model.TransactionOnAddress, 0)
	db.Raw("select a.hash, a.height, a.time, a.chain_id, b.from, b.to, b.amount, c.hash as token_hash, c.type as token_type, c.name as token_name 1 as direct from src_transactions a inner join src_transfers b on a.hash = b.tx_hash inner join tokens c on b.asset = c.hash and b.chain_id = c.chain_id where b.from = ? and b.chain_id = ?"+
		"union select d.hash, d.height, d.time, d.chain_id, e.from, e.to, e.amount, f.hash as token_hash, f.type as token_type, f.name as token_name, 2 as direct from dst_transactions d inner join dst_transfers e on d.hash = e.tx_hash inner join tokens f on e.asset = f.hash and e.chain_id = f.chain_id where e.to = ? and e.chain_id = ?"+
		"order by height desc limit ?,?",
		addressTxListReq.Address, addressTxListReq.ChainId, addressTxListReq.Address, addressTxListReq.ChainId, addressTxListReq.PageSize*addressTxListReq.PageNo, addressTxListReq.PageSize).Find(&transactionOnAddresses)
	//

	counter := new(model.Counter)
	db.Raw("select sum(cnt) as counter from (select count(*) as cnt from src_transactions a inner join src_transfers b on a.hash = b.tx_hash inner join tokens c on b.asset = c.hash and b.chain_id = c.chain_id where b.from = ? and b.chain_id = ?"+
		"union count(*) as cnt from dst_transactions d inner join dst_transfers e on d.hash = e.tx_hash inner join tokens f on e.asset = f.hash and e.chain_id = f.chain_id where e.to = ? and e.chain_id = ?) as u",
		addressTxListReq.Address, addressTxListReq.ChainId, addressTxListReq.Address, addressTxListReq.ChainId).Find(counter)
	//
	//
	c.Data["json"] = model.MakeAddressTxList(transactionOnAddresses, counter.Counter)
	c.ServeJSON()
}

// TODO GetCrossTxList gets Cross transaction list from start to end (to be optimized)
func (c *ExplorerController) GetCrossTxList() {
	// get parameter
	var crossTxListReq model.CrossTxListReq
	var err error
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &crossTxListReq); err != nil {
		c.Data["json"] = model.MakeErrorRsp(fmt.Sprintf("request parameter is invalid!"))
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
	}

	srcPolyDstRelations := make([]*model.SrcPolyDstRelation, 0)
	db.Model(&model.PolyTransaction{}).
		Select("src_transactions.hash as src_hash, poly_transactions.hash as poly_hash, dst_transactions.hash as dst_hash").
		Where("src_transactions.standard = ?", 0).
		Joins("left join src_transactions on src_transactions.hash = poly_transactions.src_hash").
		Joins("left join dst_transactions on poly_transactions.hash = dst_transactions.poly_hash").
		Preload("SrcTransaction").
		Preload("SrcTransaction.SrcTransfer").
		Preload("PolyTransaction").
		Preload("DstTransaction").
		Preload("DstTransaction.DstTransfer").
		Limit(crossTxListReq.PageSize).Offset(crossTxListReq.PageSize * crossTxListReq.PageNo).
		Find(&srcPolyDstRelations)

	var transactionNum int64
	db.Model(&model.PolyTransaction{}).Where("src_transactions.standard = ?", 0).
		Joins("left join src_transactions on src_transactions.hash = poly_transactions.src_hash").Count(&transactionNum)

	c.Data["json"] = model.MakeCrossTxListResp(srcPolyDstRelations)
	c.ServeJSON()
}

// GetCrossTx gets cross tx by Tx
func (c *ExplorerController) GetCrossTx() {
	var crossTxReq model.CrossTxReq
	var err error
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &crossTxReq); err != nil {
		c.Data["json"] = model.MakeErrorRsp(fmt.Sprintf("request parameter is invalid!"))
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
	}
	srcPolyDstRelations := make([]*model.SrcPolyDstRelation, 0)
	db.Model(&model.SrcTransaction{}).
		Select("src_transactions.hash as src_hash, poly_transactions.hash as poly_hash, dst_transactions.hash as dst_hash, src_transactions.chain_id as chain_id, src_transactions.asset as token_hash").
		Where("src_transactions.standard = ? and (src_transactions.hash = ? or poly_transactions.hash = ? or dst_transactions.hash = ?)", 0, crossTxReq.TxHash, crossTxReq.TxHash, crossTxReq.TxHash).
		Joins("left join src_transfers on src_transactions.hash = src_transfers.tx_hash").
		Joins("left join poly_transactions on src_transactions.hash = poly_transactions.src_hash").
		Joins("left join dst_transactions on poly_transactions.hash = dst_transactions.poly_hash").
		Preload("SrcTransaction").
		Preload("SrcTransaction.SrcTransfer").
		Preload("PolyTransaction").
		Preload("DstTransaction").
		Preload("DstTransaction.DstTransfer").
		Preload("Token").
		Preload("Token.TokenBasic").
		Find(&srcPolyDstRelations)
	c.Data["json"] = model.MakeCrossTxResp(srcPolyDstRelations)
	c.ServeJSON()
}
