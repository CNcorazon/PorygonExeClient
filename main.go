package main

import (
	"fmt"
	"horizon/model"
	"horizon/request"
	"horizon/structure"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	HTTPURL = "http://127.0.0.1:8088"
	WSURL   = "ws://127.0.0.1:8088"
	// HTTPURL           = "http://172.18.166.60:8800"
	// WSURL             = "ws://http://172.18.166.60:8800"
	blockTransaction      = "/block/transaction"
	blockAccount          = "/block/account"
	blockUpload           = "/block/upload"
	blockWitness          = "/block/witness"
	blockWitness_2        = "/block/witness_2"
	blockTxValidation     = "/block/validate"
	blockUploadRoot       = "/block/uploadroot"
	blockGetProposalBlock = "/block/proposalBlock"
	shardNum              = "/shard/shardNum"
	consensusflag         = "/shard/flag"
	register              = "/shard/register"
	muliticastconn        = "/shard/multicast"
	multicastblock        = "/shard/block"
	sendtvote             = "/shard/vote"
	heightNum             = "/shard/height"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	for {
		log.Printf("---------------开始交易验证---------------")
		// 交易验证哈希
		height_old := request.HeightRequest(HTTPURL, heightNum)
		for {
			time.Sleep(1 * time.Second)
			height_new := request.HeightRequest(HTTPURL, heightNum)
			if height_new != height_old {
				break
			}
		}
		validatestart := time.Now()
		// 初始化随机数生成器的种子
		rand.Seed(time.Now().UnixNano())
		// 生成一个1到shardnum（包含1和shardnum）之间的随机整数
		shard := uint(structure.SelfShardNum)
		transaction := request.RequestBlock(shard, HTTPURL, blockTxValidation)
		log.Printf("获得了第%v个区块的交易", transaction.Height)
		//time.Sleep(4 * time.Second)
		log.Printf("交易下载完成,有%v条", transaction.Num)
		DownloadTxTime := time.Since(validatestart)
		DownloadAccStart := time.Now()
		accList := request.RequestAccount(shard, HTTPURL, blockAccount)
		//time.Sleep(560 * time.Millisecond)
		log.Println("账户信息下载完成")
		DownloadAccTime := time.Since(DownloadAccStart)
		state := structure.MakeStateWithAccount(shard, accList.AccountList)
		txlist := structure.TransactionBlock{
			InternalList:   transaction.InternalList,
			CrossShardList: transaction.CrossShardList,
			SuperList:      transaction.RelayList,
		}
		validatesignstart := time.Now()
		//time.Sleep(time.Duration(40000*structure.SIGN_VERIFY_TIME) * time.Microsecond)
		req := model.GetProposalRequest{
			Height:   int(transaction.Height),
			Identity: "execute",
		}
		log.Printf("请求第%v个区块中的lockedAccount", int(transaction.Height))
		problock := model.GetProposalResponse{
			ProposalBlocks: make([]structure.Block, 0),
		}
		if transaction.Height > 0 {
			problock = request.GetProposalBlock(HTTPURL, blockGetProposalBlock, req)
		}
		root, SuList, count := structure.UpdateStateWithTxBlock(txlist, problock.ProposalBlocks[0].Body.LockedAccount, accList.AccountList, transaction.Height, state, shard)
		//log.Printf("locked Accounts:%v; valid txs: %v", problock.ProposalBlocks[0].Body.LockedAccount, count)

		ValidateSignTime := time.Since(validatesignstart)
		Uploadstart := time.Now()
		res := request.UploadRoot(shard, transaction.Height, count, root, SuList, HTTPURL, blockUploadRoot)
		UploadTime := time.Since(Uploadstart)
		Totaltime := time.Since(validatestart)
		fmt.Println(res.Message)
		str := fmt.Sprintf("验证总共用时: %v, 下载账户:%v,下载交易:%v,验证交易签名:%v, GSwrite:%v", Totaltime, DownloadAccTime, DownloadTxTime, ValidateSignTime, UploadTime)
		dstFile, err := os.OpenFile("/Users/xiading/Desktop/PorygonExeClient/validate.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer dstFile.Close()
		dstFile.WriteString(str + "\n")
		//str1 := fmt.Sprintln(accList)
		//dstFile1, err := os.OpenFile("/Users/xiading/Desktop/PorygonExeClient/transaction.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		//if err != nil {
		//	fmt.Println(err.Error())
		//	return
		//}
		//defer dstFile1.Close()
		//dstFile1.WriteString(str1 + "\n")
	}
}
