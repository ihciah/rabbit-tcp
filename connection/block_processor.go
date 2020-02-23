package connection

import (
	"context"
	"time"

	"github.com/ihciah/rabbit-tcp/block"
	"github.com/ihciah/rabbit-tcp/logger"
	"go.uber.org/atomic"
)

// 1. Join blocks from chan to connection orderedRecvQueue
// 2. Send bytes or control block
type blockProcessor struct {
	cache          map[uint32]block.Block
	logger         *logger.Logger
	relayCtx       context.Context
	removeFromPool context.CancelFunc

	sendBlockID     atomic.Uint32
	recvBlockID     uint32
	lastRecvBlockID uint32
}

func newBlockProcessor(ctx context.Context, removeFromPool context.CancelFunc) blockProcessor {
	return blockProcessor{
		cache:          make(map[uint32]block.Block),
		relayCtx:       ctx,
		removeFromPool: removeFromPool,
		logger:         logger.NewLogger("[BlockProcessor]"),
	}
}

// Join blocks and send buffer to connection
// TODO: If waiting a packet for TIMEOUT, break the connection; otherwise re-countdown for next waiting packet.
func (x *blockProcessor) OrderedRelay(connection Connection) {
	x.logger.Infof("Ordered Relay of Connection %d started.\n", connection.GetConnectionID())
	for {
		select {
		case blk := <-connection.getRecvQueue():
			if blk.BlockID+1 > x.lastRecvBlockID {
				// Update lastRecvBlockID
				x.lastRecvBlockID = blk.BlockID + 1
			}
			if x.recvBlockID == blk.BlockID {
				// Can send directly
				x.logger.Debugf("Send Block %d directly\n", blk.BlockID)
				connection.getOrderedRecvQueue() <- blk
				x.recvBlockID++
				for {
					blk, ok := x.cache[x.recvBlockID]
					if !ok {
						break
					}
					x.logger.Debugf("Send Block %d from cache\n", blk.BlockID)
					connection.getOrderedRecvQueue() <- blk
					delete(x.cache, x.recvBlockID)
					x.recvBlockID++
				}
			} else {
				// Cannot send directly
				if blk.BlockID < x.recvBlockID {
					// We don't need this old block
					x.logger.Debugf("Block %d is too old to cache\n", blk.BlockID)
					continue
				}
				x.logger.Debugf("Put Block %d to cache\n", blk.BlockID)
				x.cache[blk.BlockID] = blk
			}
		case <-time.After(PacketWaitTimeoutSec * time.Second):
			x.logger.Debugf("Packet wait time exceed of Connection %d.\n", connection.GetConnectionID())
			if x.recvBlockID == x.lastRecvBlockID {
				x.logger.Debugf("recvBlockId == lastRecvBlockID(%d), but Connection %d is not in waiting status, continue.\n", x.recvBlockID, connection.GetConnectionID())
				continue
			}
			x.logger.Warnf("Connection %d is going to be killed due to timeout.\n", connection.GetConnectionID())
			x.removeFromPool()
		case <-x.relayCtx.Done():
			x.logger.Infof("Ordered Relay of Connection %d stopped.\n", connection.GetConnectionID())
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

func (x *blockProcessor) packDisconnect(connectionID uint32, shutdownType uint8) block.Block {
	return block.NewDisconnectBlock(connectionID, x.sendBlockID.Inc()-1, shutdownType)
}
