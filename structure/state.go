package structure

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"sync"

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
		Id      int
		Shard   uint
		Address string
		Value   int
	}
	LockedAccount struct {
		AccountsByID   map[int]Account // Use a map to store accounts for quick access
		AccountsByAddr map[string]int
		Locked         map[int]bool // Store locked account Ids
	}
)

func NewLockedAccount() *LockedAccount {
	return &LockedAccount{
		AccountsByID:   make(map[int]Account),
		AccountsByAddr: make(map[string]int),
		Locked:         make(map[int]bool),
	}
}

func (la *LockedAccount) IsAccountAccessible(accountAddress string) bool {
	accountId, exists := la.AccountsByAddr[accountAddress]
	if !exists {
		return false
	}
	return !la.Locked[accountId]
}

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
func findMissingNumbers(lockedAccount []int) []int {
	// 创建一个map用于记录出现过的数字
	appeared := make(map[int]bool)

	// 标记出现过的数字
	for _, num := range lockedAccount {
		appeared[num] = true
	}

	// 找到未出现的数字
	var missingNumbers []int
	for i := 0; i <= 1000; i++ {
		if !appeared[i] {
			missingNumbers = append(missingNumbers, i)
		}
	}

	return missingNumbers
}

// UpdateState 验证交易，返回账户的树根
func UpdateState(tran TransactionBlock, lockedAccountID []int, accounts []Account, height uint, s *State, shard uint) (string, map[uint][]SuperTransaction, int) {
	//处理超级交易
	Super := tran.SuperList
	IntTraList := tran.InternalList
	CroShaList := tran.CrossShardList
	SuList := make(map[uint][]SuperTransaction)

	lockedAccounts := NewLockedAccount()
	for _, acc := range accounts {
		lockedAccounts.AccountsByID[acc.Id] = acc
		lockedAccounts.AccountsByAddr[acc.Address] = acc.Id
	}
	for _, id := range lockedAccountID {
		lockedAccounts.Locked[id] = true
	}
	var wg1 sync.WaitGroup
	var wg2 sync.WaitGroup
	accessibleTransactions1 := make(chan InternalTransaction)
	accessibleTransactions2 := make(chan CrossShardTransaction)
	for _, tx := range IntTraList[shard] {
		wg1.Add(1)
		go func(trans InternalTransaction) {
			defer wg1.Done()
			if lockedAccounts.IsAccountAccessible(trans.From) && lockedAccounts.IsAccountAccessible(trans.To) {
				accessibleTransactions1 <- trans
			}
		}(tx)
	}

	go func() {
		wg1.Wait()
		close(accessibleTransactions1)
	}()

	for _, tx := range CroShaList[shard] {
		wg2.Add(1)
		go func(trans CrossShardTransaction) {
			defer wg2.Done()
			if lockedAccounts.IsAccountAccessible(trans.From) && lockedAccounts.IsAccountAccessible(trans.To) {
				accessibleTransactions2 <- trans
			}
		}(tx)
	}

	go func() {
		wg2.Wait()
		close(accessibleTransactions2)
	}()

	var count int

	for _, trans := range Super[shard] {
		ExcuteRelay(trans, s, int(shard))
		count++
	}
	//处理内部交易
	for trans := range accessibleTransactions1 {
		//if lockedAccounts.IsAccountAccessible(trans.From) && lockedAccounts.IsAccountAccessible(trans.To) {
		ExcuteInteral(trans, s, int(shard))
		count++
	}
	//}
	//处理跨分片交易
	for trans := range accessibleTransactions2 {
		//if lockedAccounts.IsAccountAccessible(trans.From) && lockedAccounts.IsAccountAccessible(trans.To) {
		res := ExcuteCross(trans, height, s, int(shard))
		SuList[res.Shard] = append(SuList[res.Shard], *res)
		count++
		//}
	}
	return s.CalculateRoot(shard), SuList, count
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

func MakeStateWithAccount(s uint, acc []Account) *State {
	state := State{
		Shard:      s,
		AccountMap: make(map[uint]map[string]*Account),
	}
	for i := 1; i <= ShardNum; i++ {
		state.AccountMap[uint(i)] = make(map[string]*Account)
	}
	for _, account := range acc {
		state.AccountMap[account.Shard][account.Address] = &account
	}
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

// UpdateStateWithTxBlock 根据区块更新世界状态
func UpdateStateWithTxBlock(transaction TransactionBlock, lockedAccounts []int, accounts []Account, height uint, s *State, shard uint) (string, map[uint][]SuperTransaction, int) {
	root, SuList, count := UpdateState(transaction, lockedAccounts, accounts, height, s, shard)
	return root, SuList, count
}
