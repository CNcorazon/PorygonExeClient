package structure

// import "horizon/model"

type (
	PubKeySign struct {
		// Pub string
		// Sign []byte
		R string
		S string
	}

	InternalTransaction struct {
		Shard uint
		From  string
		To    string
		Value int
		// Signature []PubKeySign
	}

	CrossShardTransaction struct {
		Shard1 uint
		Shard2 uint
		From   string
		To     string
		Value  int
		// Signature []PubKeySign
	}

	SuperTransaction struct {
		Shard uint
		To    string
		Value int
		// Signature []PubKeySign
	}
)

func MakeInternalTransaction(s uint, from string, to string, value int) *InternalTransaction {
	trans := InternalTransaction{
		Shard: s,
		From:  from,
		To:    to,
		Value: value,
	}
	return &trans
}

func MakeCrossShardTransaction(s1 uint, s2 uint, from string, to string, value int) *CrossShardTransaction {
	trans := CrossShardTransaction{
		Shard1: s1,
		Shard2: s2,
		From:   from,
		To:     to,
		Value:  value,
	}
	return &trans
}
