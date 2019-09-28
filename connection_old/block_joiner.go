package connection_old

import (
	"github.com/ihciah/rabbit-tcp/block"
	"log"
)

type BlockJoiner struct {
	blockID uint32
	m       map[uint32][]byte
	conn    *RabbitTCPConn
}

func NewBlockJoiner(conn *RabbitTCPConn) *BlockJoiner {
	return &BlockJoiner{
		blockID: 1,
		m:       make(map[uint32][]byte),
		conn:    conn,
	}
}

func (j *BlockJoiner) AddBlock(block *block.Block) {
	j.m[block.BlockID] = block.BlockData
	log.Println("Block", block.BlockID, "added to BlockJoiner.")
}

func (j *BlockJoiner) GetBlock() []byte {
	if data, ok := j.m[j.blockID]; ok {
		log.Println("Block", j.blockID, "got and return connection.")
		j.blockID += 1
		delete(j.m, j.blockID)
		return data
	} else {
		log.Println("Cannot get block")
		return nil
	}
}
