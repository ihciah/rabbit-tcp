package connection

import (
	"context"
	"github.com/ihciah/rabbit-tcp/block"
	"go.uber.org/atomic"
	"log"
	"os"
)

// 1. Join blocks from chan to connection orderedRecvQueue
// 2. Send bytes or control block
type blockProcessor struct {
	cache    map[uint32]block.Block
	logger   *log.Logger
	relayCtx context.Context

	sendBlockID atomic.Uint32
	recvBlockID uint32
}

func newBlockProcessor(ctx context.Context) blockProcessor {
	return blockProcessor{
		cache:    make(map[uint32]block.Block),
		relayCtx: ctx,
		logger:   log.New(os.Stdout, "[BlockProcessor]", log.LstdFlags),
	}
}

// Join blocks and send buffer to connection
func (x *blockProcessor) OrderedRelay(connection Connection) {
	x.logger.Printf("Connection %d ordered relay started.\n", connection.GetConnectionID())
	for {
		select {
		case blk := <-connection.getRecvQueue():
			if x.recvBlockID == blk.BlockID {
				x.logger.Printf("Send %d directly\n", blk.BlockID)
				connection.getOrderedRecvQueue() <- blk
				x.recvBlockID++
				for {
					blk, ok := x.cache[x.recvBlockID]
					if !ok {
						break
					}
					x.logger.Printf("Send %d from cache\n", blk.BlockID)
					connection.getOrderedRecvQueue() <- blk
					delete(x.cache, x.recvBlockID)
					x.recvBlockID++
				}
			} else {
				x.logger.Printf("Put %d to cache\n", blk.BlockID)
				x.cache[blk.BlockID] = blk
			}
		case <-x.relayCtx.Done():
			x.logger.Printf("Connection %d ordered relay stoped.\n", connection.GetConnectionID())
			return
		}
	}
}

func (x *blockProcessor) packData(data []byte, connectionID uint32) []block.Block {
	return block.NewDataBlocks(connectionID, &x.sendBlockID, data)
}

func (x *blockProcessor) packConnect(address string, connectionID uint32) block.Block {
	return block.NewConnectBlock(connectionID, x.sendBlockID.Inc()-1, address)
}

func (x *blockProcessor) packDisconnect(connectionID uint32) block.Block {
	return block.NewDisconnectBlock(connectionID, x.sendBlockID.Inc()-1)
}
