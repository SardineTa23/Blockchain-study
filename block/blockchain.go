package block

import (
	"blockchain-study/utils"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	// マイニングで設定した一致する行頭文字数
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
)

type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
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
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string
}

// 新しいブロックチェーンの作成
func NewBlockChain(blockchainAddress string) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
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

// 引数の情報を持つTransactionを新規作成、レシーバーのTransactionPoolに追加する
func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32,
	senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	t := NewTransaction(sender, recipient, value)

	// マイニングの際は、署名チェック不要
	if sender == MINING_SENDER {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		// 所持残高が送金量に満たない時はtransaction追加処理を中止する
		/*
			if bc.CalculateTotalAmount(sender) < value {
				log.Println("ERROR: Not enough balance in a wallet")
				return false
			}
		*/

		// transactionに追加
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	} else {
		log.Println("ERROR: Verify Transaction")
	}
	return false
}

// 鍵、署名情報などを使って正当なトランザクション作成かをチェック
func (bc *Blockchain) VerifyTransactionSignature(
	serderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(serderPublicKey, h[:], s.R, s.S)
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions,
			NewTransaction(
				t.senderBlockchainAddress,
				t.recipientBlockchainAddress,
				t.value))
	}
	return transactions
}

// 成功したらtrue, 失敗したらfalseを返す。
func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{0, nonce, previousHash, transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros
}

// レシーバーのnonceの適当な値が見つかるまでvalidProofを呼び続けるメソッド
func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}

// 新規BlockをChainするために必要となるMining処理全般を扱う。
func (bc *Blockchain) Mining() bool {
	// マイニングした人へ送金するためのTransaction作成
	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD, nil, nil)

	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	log.Println("action=mining, status=success")
	return true
}

// 呼び出し時点のチェーン内で、引数の人がどれだけのValueを持っているかを返す。
func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			// 受け取ったvalueは足し算
			if blockchainAddress == t.recipientBlockchainAddress {
				totalAmount += value
			}

			// 送信したvalueは引き算
			if blockchainAddress == t.senderBlockchainAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
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
