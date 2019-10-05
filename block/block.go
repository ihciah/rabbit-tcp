package block

import (
	"encoding/binary"
	"go.uber.org/atomic"
	"io"
)

const (
	BLOCK_TYPE_CONNECT = iota
	BLOCK_TYPE_DISCONNECT
	BLOCK_TYPE_DATA

	BLOCK_HEADER_SIZE = 1 + 4 + 4 + 4
	BLOCK_DATA_SIZE   = 1000
)

type Block struct {
	Type         uint8  // 1 byte
	ConnectionID uint32 // 4 bytes
	BlockID      uint32 // 4 bytes
	BlockLength  uint32 // 4 bytes
	BlockData    []byte
}

func (block *Block) Pack() []byte {
	packedData := make([]byte, BLOCK_HEADER_SIZE+len(block.BlockData))
	packedData[0] = block.Type
	binary.LittleEndian.PutUint32(packedData[1:], block.ConnectionID)
	binary.LittleEndian.PutUint32(packedData[5:], block.BlockID)
	binary.LittleEndian.PutUint32(packedData[9:], block.BlockLength)
	copy(packedData[BLOCK_HEADER_SIZE:], block.BlockData)
	return packedData
}

func NewBlockFromReader(reader io.Reader) (*Block, error) {
	headerBuf := make([]byte, BLOCK_HEADER_SIZE)
	block := Block{}
	_, err := io.ReadFull(reader, headerBuf)
	if err != nil {
		return nil, err
	}
	// TODO: error handle
	block.Type = headerBuf[0]
	block.ConnectionID = binary.LittleEndian.Uint32(headerBuf[1:])
	block.BlockID = binary.LittleEndian.Uint32(headerBuf[5:])
	block.BlockLength = binary.LittleEndian.Uint32(headerBuf[9:])
	block.BlockData = make([]byte, block.BlockLength)
	if block.BlockLength > 0 {
		_, err = io.ReadFull(reader, block.BlockData)
		if err != nil {
			return nil, err
		}
	}
	return &block, nil
}

func NewConnectBlock(connectID uint32, blockID uint32, address string) Block {
	data := []byte(address)
	return Block{
		Type:         BLOCK_TYPE_CONNECT,
		ConnectionID: connectID,
		BlockID:      blockID,
		BlockLength:  uint32(len(data)),
		BlockData:    data,
	}
}

func newDataBlock(connectID uint32, blockID uint32, data []byte) Block {
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)
	return Block{
		Type:         BLOCK_TYPE_DATA,
		ConnectionID: connectID,
		BlockID:      blockID,
		BlockLength:  uint32(len(dataCopy)),
		BlockData:    dataCopy,
	}
}

func NewDataBlocks(connectID uint32, blockID *atomic.Uint32, data []byte) []Block {
	blocks := make([]Block, 0)
	for cursor := 0; cursor < len(data); {
		end := cursor + BLOCK_DATA_SIZE
		if len(data) < end {
			end = len(data)
		}
		blocks = append(blocks, newDataBlock(connectID, blockID.Inc()-1, data[cursor:end]))
		cursor = end
	}
	return blocks
}

func NewDisconnectBlock(connectID uint32, blockID uint32) Block {
	return Block{
		Type:         BLOCK_TYPE_DISCONNECT,
		ConnectionID: connectID,
		BlockID:      blockID,
		BlockLength:  0,
		BlockData:    make([]byte, 0),
	}
}
