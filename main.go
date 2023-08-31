package main

import (
	"fmt"
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
	blockTransaction  = "/block/transaction"
	blockAccount      = "/block/account"
	blockUpload       = "/block/upload"
	blockWitness      = "/block/witness"
	blockWitness_2    = "/block/witness_2"
	blockTxValidation = "/block/validate"
	blockUploadRoot   = "/block/uploadroot"

	shardNum       = "/shard/shardNum"
	consensusflag  = "/shard/flag"
	register       = "/shard/register"
	muliticastconn = "/shard/multicast"
	multicastblock = "/shard/block"
	sendtvote      = "/shard/vote"
	heightNum      = "/shard/height"
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
		shard := uint(rand.Intn(structure.ShardNum) + 1)
		transaction := request.RequestBlock(shard, HTTPURL, blockTxValidation)
		time.Sleep(4 * time.Second)
		log.Printf("交易下载完成,有%v条", transaction.Num)
		accList := request.RequestAccount(shard, HTTPURL, blockAccount)
		time.Sleep(560 * time.Millisecond)
		log.Println("账户信息下载完成")
		log.Printf("gsroot:%v", accList.GSRoot)
		// 带宽限制只能限制上传带宽，不能限制下载带宽，先sleep过去
		DownloadTxTime := time.Since(validatestart)
		state := structure.MakeStateWithAccount(shard, accList.AccountList, accList.GSRoot)
		txlist := structure.TransactionBlock{
			InternalList:   transaction.InternalList,
			CrossShardList: transaction.CrossShardList,
			SuperList:      transaction.RelayList,
		}
		validatesignstart := time.Now()
		time.Sleep(time.Duration(40000*structure.SIGN_VERIFY_TIME) * time.Microsecond)
		root, SuList := structure.UpdateStateWithTxBlock(txlist, transaction.Height, state, shard)
		ValidateSignTime := time.Since(validatesignstart)
		Uploadstart := time.Now()
		res := request.UploadRoot(shard, transaction.Height, transaction.Num, root, SuList, HTTPURL, blockUploadRoot)
		UploadTime := time.Since(Uploadstart)
		Totaltime := time.Since(validatestart)
		fmt.Println(res.Message)
		str := fmt.Sprintf("验证总共用时: %v, GSread:%v,验证交易签名:%v, GSwrite:%v", Totaltime, DownloadTxTime, ValidateSignTime, UploadTime)
		dstFile, err := os.OpenFile("/Users/xiading/Library/Mobile Documents/com~apple~CloudDocs/学习/中山大学/论文代码/validator/validate.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer dstFile.Close()
		dstFile.WriteString(str + "\n")
	}
}
