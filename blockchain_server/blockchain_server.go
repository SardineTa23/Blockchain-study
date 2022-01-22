package main

import (
	"blockchain-study/block"
	"blockchain-study/utils"
	"blockchain-study/wallet"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

var cache map[string]*block.Blockchain = make(map[string]*block.Blockchain)

type BlockchainServer struct {
	port uint16
}

func NewBlockChainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{port}
}

func (bsc *BlockchainServer) Port() uint16 {
	return bsc.port
}

func (bsc *BlockchainServer) GetBlockchain() *block.Blockchain {
	bc, ok := cache["blockchain"]

	// キャッシュがない場合は新たにブロックチェーンを作成
	if !ok {
		minersWallet := wallet.NewWallet()
		bc = block.NewBlockChain(minersWallet.BlockchainAddress(), bsc.Port())
		cache["blockchain"] = bc
		log.Printf("private_key %v", minersWallet.PrivateKeyStr())
		log.Printf("public_key %v", minersWallet.PublicKeyStr())
		log.Printf("blockchain_address %v", minersWallet.BlockchainAddress())
	}
	return bc
}

// GetBlockchainで取得できたブロックチェーンをJsonで返すAPI
func (bcs *BlockchainServer) GetChain(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		// ヘッダーつけて
		w.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBlockchain()

		// 取得してきたブロックチェーンをJsonへ
		m, _ := bc.MarshalJSON()

		// レスポンスボディにJsonを追加
		io.WriteString(w, string(m[:]))
	}
}

func (bcs *BlockchainServer) Transactions(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBlockchain()
		transactions := bc.TransactionPool()
		m, _ := json.Marshal(struct {
			Transacions []*block.Transaction `json:"transactions"`
			Length      int                  `json:"length"`
		}{
			Transacions: transactions,
			Length:      len(transactions),
		})

		io.WriteString(w, string(m[:]))

	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var t block.TransactionRequest

		// リクエストのJsonがTransactionRequestにデコードできない場合はerrに入る
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		if !t.Validate() {
			log.Println("ERROR: missing field(s)")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		signature := utils.SignatureFromString(*t.Signature)

		bc := bcs.GetBlockchain()

		// wallet_serverから送られてきたJsonを元に、新しいTransactionを作成
		isCreated := bc.CreateTransaction(*t.SenderBlockchainAddress,
			*t.RecipientBlockchainAddress, *t.Value, publicKey, signature)

		w.Header().Add("Content-Type", "application/json")
		var m []byte
		if !isCreated {
			w.WriteHeader(http.StatusBadRequest)
			m = utils.JsonStatus("fail")
		} else {
			w.WriteHeader(http.StatusCreated)
			m = utils.JsonStatus("succsess")
		}

		io.WriteString(w, string(m))
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) Run() {
	http.HandleFunc("/", bcs.GetChain)
	http.HandleFunc("/transactions", bcs.Transactions)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcs.port)), nil))
}
