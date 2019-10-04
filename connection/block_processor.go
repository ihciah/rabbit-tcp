package connection

import (
	"context"
	"fmt"
	"github.com/ihciah/rabbit-tcp/block"
	"log"
	"os"
	"sync"
)

// 1. Join blocks from chan to connection orderedRecvQueue
// 2. Send bytes or control block
type blockProcessor struct {
	conn            Connection
	cache           map[uint32]block.Block
	cancel          context.CancelFunc
	sendBlockIDLock sync.Mutex
	recvBlockIDLock sync.Mutex
	sendBlockID     uint32
	recvBlockID     uint32
	logger          *log.Logger
}

func newBlockProcessor(conn Connection) blockProcessor {
	return blockProcessor{
		conn:        conn,
		cache:       make(map[uint32]block.Block),
		sendBlockID: 1,
		recvBlockID: 1,
		logger:      log.New(os.Stdout, fmt.Sprintf("[BlockProcessor%d]", conn.GetConnectionID()), log.LstdFlags),
	}
}

func (x *blockProcessor) SendData(data []byte) {
	x.logger.Printf("Send data block.\n")
	x.sendBlockIDLock.Lock()
	blocks := block.NewDataBlocks(x.conn.GetConnectionID(), x.sendBlockID, data)
	x.sendBlockID += uint32(len(blocks))
	x.sendBlockIDLock.Unlock()
	for _, blk := range blocks {
		x.conn.SendBlock(blk)
	}
}

func (x *blockProcessor) SendConnect(address string) {
	x.logger.Printf("Send connect to %s block.\n", address)
	x.sendBlockIDLock.Lock()
	blkID := x.sendBlockID
	x.sendBlockID += 1
	x.sendBlockIDLock.Unlock()

	blk := block.NewConnectBlock(x.conn.GetConnectionID(), blkID, address)
	x.conn.SendBlock(blk)
}

func (x *blockProcessor) SendDisconnect() {
	x.logger.Printf("Send disconnect block.\n")
	x.sendBlockIDLock.Lock()
	blkID := x.sendBlockID
	x.sendBlockID += 1
	x.sendBlockIDLock.Unlock()

	blk := block.NewDisconnectBlock(x.conn.GetConnectionID(), blkID)
	x.conn.SendBlock(blk)
}

// Join blocks and send buffer to connection
func (x *blockProcessor) Daemon(connection Connection) {
	var ctx context.Context
	ctx, x.cancel = context.WithCancel(context.Background())
	for {
		select {
		case blk := <-x.conn.GetRecvQueue():
			if x.recvBlockID == blk.BlockID {
				connection.GetOrderedRecvQueue() <- blk
				x.recvBlockID++
				for {
					blk, ok := x.cache[x.recvBlockID]
					if !ok {
						break
					}
					connection.GetOrderedRecvQueue() <- blk
					delete(x.cache, x.recvBlockID)
					x.recvBlockID++
				}
			} else {
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
