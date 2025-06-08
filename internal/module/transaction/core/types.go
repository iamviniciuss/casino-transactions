package core

type TransactionType string

const (
	TransactionTypeBet TransactionType = "bet"
	TransactionTypeWin TransactionType = "win"
)

func (t TransactionType) IsValid() bool {
	switch t {
	case TransactionTypeBet, TransactionTypeWin:
		return true
	default:
		return false
	}
}
