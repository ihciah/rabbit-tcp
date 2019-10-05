package connection

import (
	"context"
	"github.com/ihciah/rabbit-tcp/block"
	"go.uber.org/atomic"
	"log"
	"os"
	"sync"
)

// 1. Join blocks from chan to connection orderedRecvQueue
// 2. Send bytes or control block
type blockProcessor struct {
	cache  map[uint32]block.Block
	cancel context.CancelFunc
	logger *log.Logger

	sendBlockIDLock sync.Mutex
	sendBlockID     atomic.Uint32
	recvBlockID     uint32
}

func newBlockProcessor() blockProcessor {
	return blockProcessor{
		cache:  make(map[uint32]block.Block),
		logger: log.New(os.Stdout, "[BlockProcessor]", log.LstdFlags),
	}
}

// Join blocks and send buffer to connection
func (x *blockProcessor) Daemon(connection Connection) {
	var ctx context.Context
	ctx, x.cancel = context.WithCancel(context.Background())
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
		case <-ctx.Done():
			return
		}
	}
}

func (x *blockProcessor) StopDaemon() {
	if x.cancel != nil {
		x.cancel()
	}
}

func (x *blockProcessor) packData(data []byte, connectionID uint32) []block.Block {
	x.sendBlockIDLock.Lock()
	blocks := block.NewDataBlocks(connectionID, &x.sendBlockID, data)
	x.sendBlockIDLock.Unlock()
	return blocks
}

func (x *blockProcessor) packConnect(address string, connectionID uint32) block.Block {
	x.sendBlockIDLock.Lock()
	blkID := x.sendBlockID.Inc()
	x.sendBlockIDLock.Unlock()

	return block.NewConnectBlock(connectionID, blkID-1, address)
}

func (x *blockProcessor) packDisconnect(connectionID uint32) block.Block {
	x.sendBlockIDLock.Lock()
	blkID := x.sendBlockID.Inc()
	x.sendBlockIDLock.Unlock()

	return block.NewDisconnectBlock(connectionID, blkID-1)
}
