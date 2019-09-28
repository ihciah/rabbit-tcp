package connection

import "github.com/ihciah/rabbit-tcp/block"

type blockProcessor struct {
	blockBuffer chan block.Block
	connection  *Connection
	cache       map[uint32]block.Block
	sendBlockID uint32
	recvBlockID uint32
}

func (x *blockProcessor) SendBytes(data []byte) {

}

func (x *blockProcessor) RecvBlock(blk block.Block) {
	x.blockBuffer <- blk
}

func (x *blockProcessor) Connect() {

}

func (x *blockProcessor) Disconnect() {

}

// Join blocks and send buffer to connection
func (x *blockProcessor) Daemon() {
	for {

		blk := <-x.blockBuffer
		if x.recvBlockID == blk.BlockID {
			x.connection.dataQueue <- blk
			x.recvBlockID++
			for {
				blk, ok := x.cache[x.recvBlockID]
				if !ok {
					break
				}
				x.connection.dataQueue <- blk
				delete(x.cache, x.recvBlockID)
				x.recvBlockID++
			}
		} else {
			x.cache[blk.BlockID] = blk
		}
	}
}
