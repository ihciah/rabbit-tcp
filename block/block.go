package block

import (
	"encoding/binary"
	"io"
)

const (
	BLOCK_TYPE_REG = iota
	BLOCK_TYPE_CONNECT
	BLOCK_TYPE_DISCONNECT
	BLOCK_TYPE_DATA

	BLOCK_HEADER_SIZE = 1 + 4 + 4 + 2
	BLOCK_DATA_SIZE   = 1000
)

type Block struct {
	Type         uint8  // 1 byte
	ConnectionID uint32 // 4 bytes
	BlockID      uint32 // 4 bytes
	BlockLength  uint16 // 2 bytes
	BlockData    []byte
}

func (block *Block) Pack() []byte {
	packedData := make([]byte, BLOCK_HEADER_SIZE+len(block.BlockData))
	packedData[0] = block.Type
	binary.LittleEndian.PutUint32(packedData[1:], block.ConnectionID)
	binary.LittleEndian.PutUint32(packedData[5:], block.BlockID)
	binary.LittleEndian.PutUint16(packedData[9:], block.BlockLength)
	copy(packedData[BLOCK_HEADER_SIZE:], block.BlockData)
	return packedData
}

func NewBlockFromReader(reader io.Reader) *Block {
	headerBuf := make([]byte, BLOCK_HEADER_SIZE)
	block := Block{}
	io.ReadFull(reader, headerBuf)
	// TODO: error handle
	block.Type = headerBuf[0]
	block.ConnectionID = binary.LittleEndian.Uint32(headerBuf[1:])
	block.BlockID = binary.LittleEndian.Uint32(headerBuf[5:])
	block.BlockLength = binary.LittleEndian.Uint16(headerBuf[9:])
	block.BlockData = make([]byte, block.BlockLength)
	io.ReadFull(reader, block.BlockData)
	// TODO: error handle
	return &block
}

func NewConnectBlock(connectID uint32, address string) Block {
	data := []byte(address)
	return Block{
		Type:         BLOCK_TYPE_CONNECT,
		ConnectionID: connectID,
		BlockID:      0,
		BlockLength:  uint16(len(data)),
		BlockData:    data,
	}
}

func NewRegBlock(clientID uint16) Block {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, clientID)
	return Block{
		Type:         BLOCK_TYPE_REG,
		ConnectionID: 0,
		BlockID:      0,
		BlockLength:  uint16(len(data)),
		BlockData:    data,
	}
}

func newDataBlock(connectID uint32, blockID uint32, data []byte) Block {
	return Block{
		Type:         BLOCK_TYPE_DATA,
		ConnectionID: connectID,
		BlockID:      blockID,
		BlockLength:  uint16(len(data)),
		BlockData:    data,
	}
}

func NewDataBlocks(connectID uint32, blockID uint32, data []byte) []Block {
	blocks := make([]Block, 0)
	for cursor := 0; cursor < len(data); {
		end := cursor + BLOCK_DATA_SIZE
		if len(data) < end {
			end = len(data)
		}
		blocks = append(blocks, newDataBlock(connectID, blockID, data[cursor:end]))
		blockID += 1
		cursor = end
	}
	return blocks
}

func NewDisconnectBlock(connectID uint32) Block {
	return Block{
		Type:         BLOCK_TYPE_DISCONNECT,
		ConnectionID: connectID,
		BlockID:      0,
		BlockLength:  0,
		BlockData:    make([]byte, 0),
	}
}
