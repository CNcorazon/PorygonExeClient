package structure

import (
	"os"
	"strconv"
)

const ShardNum1 = 10
const AccountNum = 1000
const ProposerNum1 = 20
const CLIENT_MAX = 10
const SIGN_VERIFY_TIME = 300 //microsecond
const TX_NUM1 = 20000        //per shard per catagory
const CORE = 1
const NodeNum1 = ShardNum1 * CLIENT_MAX

var ShardNum int
var ProposerNum int
var TX_NUM int
var NodeNum int

func init() {
	ShardNum = ShardNum1
	ProposerNum = ProposerNum1
	TX_NUM = TX_NUM1

	// 如果命令行参数存在，尝试将其转换为整数并修改 ModifiedValue
	if len(os.Args) > 1 {
		ModifiedShardNum, err := strconv.Atoi(os.Args[1])
		if err == nil {
			ShardNum = ModifiedShardNum
			// log.Println(ShardNum)
		}
		ModifiedProposerNum, err := strconv.Atoi(os.Args[2])
		if err == nil {
			ProposerNum = ModifiedProposerNum
			// log.Println(ProposerNum)
		}
		ModifiedTXNUM, err := strconv.Atoi(os.Args[3])
		if err == nil {
			TX_NUM = ModifiedTXNUM
			// log.Println(TX_NUM)
		}
	}

	NodeNum = ShardNum * CLIENT_MAX
}
