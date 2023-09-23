package structure

import (
	"os"
	"strconv"
)

const ShardNum1 = 2
const AccountNum = 500
const ProposerNum1 = 10
const CLIENT_MAX = 10
const SIGN_VERIFY_TIME = 300 //microsecond
const ValidateTxNum1 = 4000  //per shard per catagory
const CORE = 1
const NodeNum1 = ShardNum1 * CLIENT_MAX
const SelfShardNum1 = 1

var ShardNum int
var ProposerNum int
var ValidateTxNum int
var SelfShardNum int

func init() {
	SelfShardNum = SelfShardNum1
	ShardNum = ShardNum1
	ProposerNum = ProposerNum1
	ValidateTxNum = ValidateTxNum1

	// 如果命令行参数存在，尝试将其转换为整数并修改 ModifiedValue
	if len(os.Args) > 1 {
		ModifiedSelfShardNum, err := strconv.Atoi(os.Args[1])
		if err == nil {
			SelfShardNum = ModifiedSelfShardNum
		}

		ModifiedShardNum, err := strconv.Atoi(os.Args[2])
		if err == nil {
			ShardNum = ModifiedShardNum
			// log.Println(ShardNum)
		}
		ModifiedProposerNum, err := strconv.Atoi(os.Args[3])
		if err == nil {
			ProposerNum = ModifiedProposerNum
			// log.Println(ProposerNum)
		}
		ModifiedTXNUM, err := strconv.Atoi(os.Args[4])
		if err == nil {
			ValidateTxNum = ModifiedTXNUM
			// log.Println(TX_NUM)
		}
	}

	//NodeNum = ShardNum * CLIENT_MAX
}
