package distributor

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Claim struct {
	Index   int            `json:"index"`
	Account common.Address `json:"account"`
	Amount  *big.Int       `json:"amount"`
	Proof   []common.Hash  `json:"proof"`
}

type MerkleDistributorInfo struct {
	MerkleRoot common.Hash `json:"merkleRoot"`
	TokenTotal *big.Int    `json:"tokenTotal"`
	Claims     []Claim     `json:"claims"`
}

func ParseBalanceMap(balances []Balance) (MerkleDistributorInfo, error) {
	info := MerkleDistributorInfo{
		Claims: make([]Claim, 0),
	}

	tree, err := NewBalanceTree(balances)
	if err != nil {
		return info, err
	}

	tokenTotal := big.NewInt(0)
	for idx, balance := range balances {
		proof, err := tree.GetProof(idx, balance.Account, balance.Amount)
		if err != nil {
			return info, err
		}
		tokenTotal = big.NewInt(0).Add(tokenTotal, balance.Amount)

		info.Claims = append(info.Claims, Claim{
			Index:   idx,
			Account: balance.Account,
			Amount:  balance.Amount,
			Proof:   proof,
		})
	}

	info.MerkleRoot = tree.GetRoot()
	info.TokenTotal = tokenTotal

	return info, nil
}
