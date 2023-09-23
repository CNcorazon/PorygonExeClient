package structure

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"time"
)

type (
	Block struct {
		Header BlockHeader
		Body   BlockBody
	}

	Root        string
	BlockHeader struct {
		// Shard                uint   //表示是第几号分片中的区块
		Height          uint   //当前区块的高度
		Time            int64  //区块产生的时候的Unix时间戳
		Vote            uint   //本区块收到的移动节点的票数
		TransactionRoot string //修改了本分片状态的交易区块的SHA256值
		// SuperTransactionRoot string //产生的超级交易区块的SHA256值
		StateRoot GSRoot //当前执行完本交易之后，当前区块链账本的世界状态
	}

	BlockBody struct {
		// Shard            uint
		Height           uint
		LockedAccount    []int
		TransactionRoots map[uint]string
		// SuperTransaction SuperTransactionBlock_Consensus
	}
	Proposal struct {
		Shard         uint
		Height        uint
		InternalBatch []Batchtrans
		CrossBatch    []Batchtrans
		SuperBatch    []Batchtrans
	}
	Batchtrans struct {
		Id             int
		Shard          int
		Abstract       string
		RelatedAccount []int // 第一维是内部交易的，第二维是跨分片from方的，第三维是跨分片to方的
		PubIndex       []int
		Sig            string
		TransType      int
	}
	GSRoot struct {
		StateRoot string
		Vote      map[uint]map[string]int //记录每个执行分片计算出的subTreeRoot以及对应的票数
	}

	TransactionBlock struct {
		// Id             string //区块见证人的id，这边是因为签名问题（sign的东西要不同）补充的
		Height         uint
		InternalList   map[uint][]InternalTransaction
		CrossShardList map[uint][]CrossShardTransaction
		SuperList      map[uint][]SuperTransaction //需要被打包进这个区块内部的SueprList
		Sign           PubKeySign
	}

	// ProposalBlock 共识阶段：排序委员会提出的proposal block
	ProposalBlock struct {
		Id            string
		IdList        []string
		LockedAccount []int
		Height        int
		Hash          string     //前一个区块的hash
		Root          GSRoot     //状态树树根
		ProposalList  []Proposal //交易列表
		Vrf           int        //vrf
	}

	// TransactionBlock_Consensus struct {
	// 	InternalList   map[uint][]string
	// 	CrossShardList map[uint][]string
	// 	SuperList      map[uint][]string //需要被打包进这个区块内部的SueprList
	// }

	SuperTransactionBlock struct {
		SuperList map[uint][]SuperTransaction //执行完成TransactionList之后生成的一个SuperList
	}

	// SuperTransactionBlock_Consensus struct {
	// 	SuperList map[uint]string //执行完成TransactionList之后生成的一个SuperList
	// }
)

func (t *TransactionBlock) CalculateRoot() string {
	jsonString, err := json.Marshal(t)
	if err != nil {
		log.Fatalln("计算交易区块Root失败")
	}
	byte32 := sha256.Sum256(jsonString)
	return hex.EncodeToString(byte32[:])
}

func (r *SuperTransactionBlock) CalculateRoot() string {
	jsonString, err := json.Marshal(r)
	if err != nil {
		log.Fatalln("计算接力交易区块Root失败")
	}
	byte32 := sha256.Sum256(jsonString)
	return hex.EncodeToString(byte32[:])
}

// func MakeTransactionBlock(IntTraList map[uint][]string, CroList map[uint][]string, SuList map[uint][]string) TransactionBlock_Consensus {
// 	res := TransactionBlock_Consensus{
// 		InternalList:   IntTraList,
// 		CrossShardList: CroList,
// 		SuperList:      SuList,
// 	}
// 	return res
// }

func MakeBlock(TxRoot map[uint]string, s *State, height uint, gsroot GSRoot) Block {
	//首先打包生成本区快的交易区块
	//transBlock := MakeTransactionBlock(IntTraList, CroList, SuList)
	//根据打包好的交易区块，记录生成的接力交易区块
	//只执行跨分片交易且不修改状态
	// SuList_new := make(map[uint][]SuperTransaction, ShardNum)
	// for i := 1; i <= ShardNum; i++ {
	// 	for _, tran := range CroList[uint(i)] {
	// 		res := ExcuteCross(tran, height, s, i)
	// 		SuList_new[uint(i)] = append(SuList_new[uint(i)], *res)
	// 	}
	// }

	// SuperBlock := SuperTransactionBlock{
	// 	SuperList: SuList_new,
	// }

	body := BlockBody{
		// Shard:            s.Shard,
		Height:           height,
		TransactionRoots: TxRoot,
		// SuperTransaction: SuperBlock,
	}

	root, _ := json.Marshal(TxRoot)
	byte32 := sha256.Sum256(root)
	tr := hex.EncodeToString(byte32[:])
	header := BlockHeader{

		// Shard:                s.Shard,
		Height:          height,
		Time:            time.Now().Unix(),
		Vote:            1,
		TransactionRoot: tr,
		StateRoot:       gsroot,
	}

	block := Block{
		Header: header,
		Body:   body,
	}

	return block
}

func CompareBlocks(b1 Block, b2 Block) bool {
	//只考虑BlockHeader中的三个root,time无所谓
	if b1.Header.Height != b2.Header.Height {
		log.Printf("两个区块的高度不同")
		return false
	} else if b1.Header.TransactionRoot != b2.Header.TransactionRoot {
		log.Printf("两个区块的交易列表不同")
		log.Println(b1.Header.TransactionRoot)
		log.Println(b2.Header.TransactionRoot)
		return false
		// } else if b1.Header.SuperTransactionRoot != b2.Header.SuperTransactionRoot {
		// 	log.Printf("两个区块生成的接力交易列表不同")
		// 	return false
	} else if b1.Header.StateRoot.StateRoot != b2.Header.StateRoot.StateRoot {
		log.Printf("两个区块的state root不同")
		return false
	}
	log.Printf("高度为%v的区块经过检验", b1.Header.Height)
	return true
}
