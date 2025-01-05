package blockchain

import "sync"

type Blockchain struct {
	Blocks []*Block
	mu     sync.Mutex
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		Blocks: []*Block{CreateGenesisBlock()},
	}
}

func (bc *Blockchain) AddBlock(data string) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := GenerateNewBlock(prevBlock, data)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func (bc *Blockchain) GetBlocks() []*Block {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	return bc.Blocks
}
