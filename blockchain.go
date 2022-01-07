package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type Block struct {
	nonce        int
	previousHash string
	timestamp    int64
	transactions []string
}

// Blockの作成
func NewBlock(nonce int, previousHash string) *Block {
	// newはポインタが返る
	b := new(Block)

	// timestampのint64を返す
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash

	return b
}

func (b *Block) Print() {
	fmt.Printf("timestamp     %d\n", b.timestamp)
	fmt.Printf("nonce     %d\n", b.nonce)
	fmt.Printf("previousHash     %s\n", b.previousHash)
	fmt.Printf("transactions     %s\n", b.transactions)
}

type Blockchain struct {
	transactionPool []string
	chain           []*Block
}

// 新しいブロックチェーンの作成
func NewBlockChain() *Blockchain {
	bc := new(Blockchain)
	bc.CreateBlock(0, "init hash")
	return bc
}

// chainするBlockを作成してチェーンに追加
func (bc *Blockchain) CreateBlock(nonce int, previousHash string) *Block {
	b := NewBlock(nonce, previousHash)
	bc.chain = append(bc.chain, b)
	return b
}

// チェーンの中身をループで回して、添字とともに表示
func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	blockchain := NewBlockChain()
	blockchain.Print()
	blockchain.CreateBlock(5, "hash 1")
	blockchain.Print()
	blockchain.CreateBlock(2, "hash 2")
	blockchain.Print()
}
