package structure

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/pochard/commons/randstr"
)

type (
	State struct {
		Shard     uint                    //表示该移动节点位于哪个分片中
		RootsVote map[uint]map[string]int //记录各个分片新状态的投票数
		// NewAccountMap map[uint]map[string]*Account
		AccountMap map[uint]map[string]*Account
	}

	Account struct {
		Shard   uint
		Address string
		Value   int
	}
)

//计算账户的状态
func (s *State) CalculateRoot(shard uint) string {
	jsonString, err := json.Marshal(s.AccountMap[shard])
	if err != nil {
		log.Fatalln("计算账户状态Root失败")
	}
	byte32 := sha256.Sum256(jsonString)
	return hex.EncodeToString(byte32[:])
}

//往全局状态中添加账户
func (s *State) AppendAccount(acc Account) {
	key := acc.Address
	s.AccountMap[acc.Shard][key] = &acc
	// fmt.Printf("1321+%p", &acc)
	// s.LogState(0)
	log.Printf("分片%v添加账户成功，账户地址为%v\n", acc.Shard, key)
}

// 验证交易，返回账户的树根
func UpdateState(tran TransactionBlock, height uint, s *State, shard uint) (string, map[uint][]SuperTransaction) {
	//处理超级交易
	Super := tran.SuperList
	IntTraList := tran.InternalList
	CroShaList := tran.CrossShardList
	SuList := make(map[uint][]SuperTransaction)

	// for i := 1; i <= ShardNum; i++ {
	for _, tran := range Super[shard] {
		ExcuteRelay(tran, s, int(shard))
	}
	//处理内部交易
	for _, tran := range IntTraList[shard] {
		ExcuteInteral(tran, s, int(shard))
	}
	//处理跨分片交易
	for _, tran := range CroShaList[shard] {
		res := ExcuteCross(tran, height, s, int(shard))
		SuList[res.Shard] = append(SuList[res.Shard], *res)
	}
	// }
	return s.CalculateRoot(shard), SuList
}

func ExcuteInteral(i InternalTransaction, s *State, shardNum int) {
	if uint(shardNum) != i.Shard {
		log.Printf("节点分片%v, 交易分片%v", shardNum, i.Shard)
		log.Fatalln("该交易不由本分片进行处理")
		return
	}
	Payer := i.From
	Beneficiary := i.To
	Value := i.Value
	// fmt.Println(Payer)
	// fmt.Println(Beneficiary)
	// _, flag := s.AccountMap[Payer]
	// if !flag {
	// 	log.Fatalf("该交易的付款者不是本分片的账户")
	// 	return
	// }
	// _, flag = s.AccountMap[Beneficiary]
	// if !flag {
	// 	log.Fatalf("该交易的收款者不是本分片的账户")
	// 	return
	// }

	// s.AccountMap[Payer].Value = s.AccountMap[Payer].Value + i.Value
	// s.AccountMap[Beneficiary].Value = s.AccountMap[Beneficiary].Value + i.Value

	value1 := s.AccountMap[uint(shardNum)][Payer].Value - Value
	s.AccountMap[uint(shardNum)][Payer].Value = value1
	// log.Printf("%+v\n", *s.AccountMap[Payer])
	// log.Printf("%+v\n", (*s.AccountMap[Beneficiary]))
	value2 := s.AccountMap[uint(shardNum)][Beneficiary].Value + Value
	s.AccountMap[uint(shardNum)][Beneficiary].Value = value2
	// log.Printf("%+v\n", (*s.AccountMap[Beneficiary]))
}

func ExcuteCross(e CrossShardTransaction, height uint, s *State, shardNum int) *SuperTransaction {
	if uint(shardNum) != e.Shard1 {
		log.Fatalln("该交易的发起用户不是本分片账户")
		return nil
	}
	Payer := e.From
	_, flag := s.AccountMap[uint(shardNum)][Payer]
	if !flag {
		log.Fatalf("该交易的付款者不是本分片的账户")
		return nil
	}
	s.AccountMap[uint(shardNum)][Payer].Value = s.AccountMap[uint(shardNum)][Payer].Value - e.Value
	res := SuperTransaction{
		Shard: e.Shard2,
		To:    e.To,
		Value: e.Value,
	}
	return &res
}

func ExcuteRelay(r SuperTransaction, s *State, shardNum int) {
	if uint(shardNum) != r.Shard {
		log.Fatalf("该交易不是由本分片执行")
		return
	}
	Beneficiary := r.To
	_, flag := s.AccountMap[uint(shardNum)][Beneficiary]
	if !flag {
		log.Fatalf("该交易的收款者不是本分片的账户")
		return
	}
	s.AccountMap[uint(shardNum)][Beneficiary].Value = s.AccountMap[uint(shardNum)][Beneficiary].Value + r.Value
}

//获取某一个分片中的当前所有的账户的状态
func (s *State) GetAccountList() []Account {
	var acc []Account
	for _, v := range s.AccountMap[uint(ShardNum)] {
		acc = append(acc, *v)
	}
	return acc
}

//为执行分片初始化生成n*shardNum个AccountList
func InitAccountList(shardNum int, n int) []Account {
	var accList []Account
	for j := 1; j < shardNum; j++ {
		addressList := GenerateAddressList(n)
		for i := 0; i < n; i++ {
			acc := Account{
				Shard:   uint(j),
				Address: addressList[i],
				Value:   100000, //初始化的Value设置
			}
			accList = append(accList, acc)
		}
	}
	return accList
}

func GenerateKey() string {
	return randstr.RandomAlphanumeric(16)
}

func GenerateAddressList(n int) []string {
	set := make(map[string]struct{})
	for len(set) < n {
		key := GenerateKey()
		set[key] = struct{}{}
	}
	var res []string
	for key := range set {
		res = append(res, key)
	}
	return res
}

//初始化构建本分片的全局状态
//s表示生成的状态的分片序列号
//n表示需要初始化的账户数目
func InitState(s uint, n int, shardNum int) *State {
	state := State{
		Shard:      s,
		RootsVote:  make(map[uint]map[string]int, shardNum),
		AccountMap: make(map[uint]map[string]*Account, shardNum),
	}
	accountList := InitAccountList(shardNum, n)
	for _, x := range accountList {
		// fmt.Printf("123%+v\n", x)
		state.AppendAccount(x)
	}
	return &state
}

//前端根据传输来的账户的状态重新构造全局状态
func MakeStateWithAccount(s uint, acc []Account, gsroot GSRoot) *State {
	state := State{
		Shard:      s,
		RootsVote:  make(map[uint]map[string]int),
		AccountMap: make(map[uint]map[string]*Account),
	}
	for i := 1; i <= ShardNum; i++ {
		state.AccountMap[uint(i)] = make(map[string]*Account)
	}
	for _, account := range acc {
		state.AccountMap[account.Shard][account.Address] = &account
	}
	state.RootsVote = gsroot.Vote
	return &state
}

func (s *State) LogState(height uint) {
	fmt.Printf("当前的区块高度是%v,此时的账户状态是\n", height)
	for i := 0; i < ShardNum; i++ {
		for key, acc := range s.AccountMap[uint(i)] {
			fmt.Printf("账户{%v}的余额为{%v}\n", key, acc.Value)

		}
	}
}

//根据区块更新世界状态
func UpdateStateWithTxBlock(transaction TransactionBlock, height uint, s *State, shard uint) (string, map[uint][]SuperTransaction) {
	root, SuList := UpdateState(transaction, height, s, shard)
	return root, SuList
}
