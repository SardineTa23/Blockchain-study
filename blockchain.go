package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type Block struct {
	nonce        int
	previousHash [32]byte
	timestamp    int64
	transactions []*Transaction
}

// Blockの作成
func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	// newはポインタが返る
	b := new(Block)

	// timestampのint64を返す
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = transactions

	return b
}

func (b *Block) Print() {
	fmt.Printf("timestamp     %d\n", b.timestamp)
	fmt.Printf("nonce     %d\n", b.nonce)
	fmt.Printf("previous_hash     %x\n", b.previousHash)
	for _, t := range b.transactions {
		t.Print()
	}
}

// レシーバーのブロックをハッシュ化したものを返す。
func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256([]byte(m))
}

// ただのjson.Marshalではプライベートなプロパティにアクセスできないため、Marshalを上書き
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash [32]byte       `json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}

type Blockchain struct {
	transactionPool []*Transaction
	chain           []*Block
}

// 新しいブロックチェーンの作成
func NewBlockChain() *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.CreateBlock(0, b.Hash())
	return bc
}

// chainするBlockを作成してチェーンに追加
// レシーバーのTransactionPoolの中身を、作成するBlockのTransactionsにいれて、レシーバーのPoolは空にする
func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

// レシーバーのチェーンで繋がっている最後のブロックを取得
func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

// チェーンの中身をループで回して、添字とともに表示
func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

// 引数情報を持つTransactionを新規作成、レシーバーのTransactionPoolに追加する
func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32) {
	t := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
}

type Transaction struct {
	// 例：送金した人
	senderBlockchainAddress string

	// 例：お金を受け取った人
	recipientBlockchainAddress string

	// 例：送金した内容（金額など）
	value float32
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, value}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_blockchain_address        %s\n", t.senderBlockchainAddress)
	fmt.Printf(" recipient_blockchain_address     %s\n", t.recipientBlockchainAddress)
	fmt.Printf(" value                            %.1f\n", t.value)

}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})
}

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {

	blockchain := NewBlockChain()
	blockchain.Print()

	blockchain.AddTransaction("A", "B", 1.0)
	previousHash := blockchain.LastBlock().Hash()
	blockchain.CreateBlock(5, previousHash)
	blockchain.Print()

	blockchain.AddTransaction("A", "B", 2.0)
	blockchain.AddTransaction("A", "B", 3.0)
	previousHash = blockchain.LastBlock().Hash()
	blockchain.CreateBlock(2, previousHash)
	blockchain.Print()
}
