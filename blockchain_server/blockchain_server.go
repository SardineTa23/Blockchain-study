package main

import (
	"blockchain-study/block"
	"blockchain-study/wallet"
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

func (bcs *BlockchainServer) Run() {
	http.HandleFunc("/", bcs.GetChain)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcs.port)), nil))
}
