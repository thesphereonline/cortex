package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Block struct {
	Index        int       `json:"index"`
	Timestamp    time.Time `json:"timestamp"`
	PreviousHash string    `json:"previousHash"`
	Hash         string    `json:"hash"`
	Data         string    `json:"data"`
}

func CreateGenesisBlock() *Block {
	return &Block{
		Index:        0,
		Timestamp:    time.Now(),
		PreviousHash: "0",
		Hash:         "GENESIS_BLOCK",
		Data:         "Cortex Blockchain Genesis Bock",
	}
}

func GenerateHash(block Block) string {
	record := string(block.Index) + block.Timestamp.String() + block.PreviousHash + block.Data
	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}

func GenerateNewBlock(prevBlock *Block, data string) *Block {
	newBlock := &Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now(),
		PreviousHash: prevBlock.Hash,
		Data:         data,
	}
	newBlock.Hash = GenerateHash(*newBlock)
	return newBlock
}
