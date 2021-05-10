package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)
const targetBits = 20
var (
	maxNonce = math.MaxInt64
)
//代表一次工作量证明
type ProofOfWork struct {
	block *Block //当前需要验证的区块
	target *big.Int //挖矿难度
}

//数据拼接，返回字节数组
func (w *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			w.block.PrevBlockHash,
			w.block.Data,
			IntToHex(w.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
		)
	return data
}

func (w ProofOfWork) Run() (int, []byte) {
	var hash [32]byte
	var hashInt big.Int
	nonce := 0
	fmt.Printf("Mining the block containing \"%s\"\n", w.block.Data)
	for nonce < maxNonce {
		data := w.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])
		// hashInt < w.target = -1
		// hashInt = w.target = 0
		// hashInt > w.target > -1
		if hashInt.Cmp(w.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce,hash[:]
}

func (w *ProofOfWork) Validate() bool  {
	var hashInt big.Int
	data := w.prepareData(w.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(w.target) == -1
	return isValid
}

func NewProofOfWork(block *Block) *ProofOfWork  {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	pow := &ProofOfWork{block, target}
	return pow
}
