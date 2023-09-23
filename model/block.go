package model

import (
	"horizon/structure"
)

type (
	BlockTransactionRequest struct {
		Shard uint
		Id    string
	}

	BlockTransactionResponse struct {
		Shard  uint
		Height uint //需要生成的区块的高度，是当前区块链的高度+1
		Num    int
		// TxRoot         map[uint]string
		InternalList   map[uint][]structure.InternalTransaction
		CrossShardList map[uint][]structure.CrossShardTransaction
		RelayList      map[uint][]structure.SuperTransaction
		Signature      []structure.PubKeySign
	}

	BlockPackedTransactionResponse struct {
		Shard      uint
		Height     uint
		RootString map[uint]string
		Signsmap   map[uint][]structure.PubKeySign
	}

	BlockRequest struct {
		Shard  uint
		Height uint
	}
	BlockAccountRequest struct {
		Shard uint
	}

	BlockAccountResponse struct {
		Shard       uint
		Height      uint //当前区块链的高度
		AccountList []structure.Account
		GSRoot      structure.GSRoot
	}

	BlockUploadRequest struct {
		Shard  uint
		Height uint
		Id     string
		Block  structure.Block
		// ReLayList map[uint][]structure.SuperTransaction
	}

	BlockUploadResponse struct {
		Shard   uint
		Height  uint
		Message string
	}

	TxWitnessRequest struct {
		Shard uint
	}

	TxWitnessResponse struct {
		// Id             string
		Shard          uint
		Height         uint
		Num            int
		InternalList   map[uint][]structure.InternalTransaction
		CrossShardList map[uint][]structure.CrossShardTransaction
		RelayList      map[uint][]structure.SuperTransaction
		Sign           structure.PubKeySign
	}

	TxWitnessRequest_2 struct {
		// Id             string //区块见证者的id
		Shard          uint
		Height         uint
		Num            int
		InternalList   map[uint][]structure.InternalTransaction
		CrossShardList map[uint][]structure.CrossShardTransaction
		RelayList      map[uint][]structure.SuperTransaction
		Sign           structure.PubKeySign
		// R              big.Int
	}

	TxWitnessResponse_2 struct {
		Message string
		Flag    bool
	}

	RootUploadRequest struct {
		Shard  uint
		Height uint
		TxNum  int
		Id     string
		Root   string
		SuList map[uint][]structure.SuperTransaction
	}

	RootUploadResponse struct {
		// Shard   uint
		Height  uint
		Message string
	}
	GetProposalRequest struct {
		Height   int
		Identity string
	}
	GetProposalResponse struct {
		ProposalBlocks []structure.Block
	}
)
