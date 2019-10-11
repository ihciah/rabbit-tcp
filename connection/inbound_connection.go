package connection

import (
	"context"
	"fmt"
	"github.com/ihciah/rabbit-tcp/block"
	"github.com/ihciah/rabbit-tcp/logger"
	"go.uber.org/atomic"
	"io"
	"math/rand"
	"net"
	"time"
)

type InboundConnection struct {
	baseConnection
	dataBuffer LoopByteBuffer
}

func NewInboundConnection(sendQueue chan<- block.Block, ctx context.Context, removeFromPool context.CancelFunc) Connection {
	connectionID := rand.Uint32()
	c := InboundConnection{
		baseConnection: baseConnection{
			blockProcessor:   newBlockProcessor(ctx, removeFromPool),
			connectionID:     connectionID,
			closed:           atomic.NewBool(false),
			sendQueue:        sendQueue,
			recvQueue:        make(chan block.Block, RecvQueueSize),
			orderedRecvQueue: make(chan block.Block, OrderedRecvQueueSize),
			logger:           logger.NewLogger(fmt.Sprintf("[InboundConnection-%d]", connectionID)),
		},
		dataBuffer: NewLoopBuffer(block.MaxSize),
	}
	c.logger.Infof("InboundConnection %d created.\n", connectionID)
	return &c
}

func (c *InboundConnection) Read(b []byte) (n int, err error) {
	readN := 0
	if !c.dataBuffer.Empty() {
		// Read data from buffer
		readN += c.dataBuffer.Read(b)
		if readN == len(b) {
			// if dst is full, return
			return readN, nil
		}
	}

	if readN == 0 {
		// if no data in b, we must wait until a block comes and read something
		if c.closed.Load() {
			return 0, io.EOF
		}
		blk := <-c.orderedRecvQueue
		switch blk.Type {
		case block.TypeDisconnect:
			c.closed.Store(true)
			return 0, io.EOF
		case block.TypeData:
			dst := b[readN:]
			if len(dst) < len(blk.BlockData) {
				// if dst can't put a block, put part of it and return
				c.dataBuffer.OverWrite(blk.BlockData)
				readN += c.dataBuffer.Read(dst)
				return readN, nil
			}
			// if dst can put a block, put it
			readN += copy(dst, blk.BlockData)
		}
	}

	// readN should be more than 0 here
	for {
		select {
		case blk := <-c.orderedRecvQueue:
			switch blk.Type {
			case block.TypeDisconnect:
				c.closed.Store(true)
				return readN, nil
			case block.TypeData:
				dst := b[readN:]
				if len(dst) < len(blk.BlockData) {
					// if dst can't put a block, put part of it and return
					c.dataBuffer.OverWrite(blk.BlockData)
					readN += c.dataBuffer.Read(dst)
					return readN, nil
				}
				// if dst can put a block, put it
				readN += copy(dst, blk.BlockData)
			}
		default:
			return readN, nil
		}
	}
}

func (c *InboundConnection) Write(b []byte) (n int, err error) {
	// TODO: tag all blocks from b using WaitGroup
	// TODO: and wait all blocks sent
	c.sendData(b)
	return len(b), nil
}

func (c *InboundConnection) Close() error {
	// TODO
	c.SendDisconnect()
	return nil
}

func (c *InboundConnection) LocalAddr() net.Addr {
	// TODO
	return nil
}

func (c *InboundConnection) RemoteAddr() net.Addr {
	// TODO
	return nil
}

func (c *InboundConnection) SetDeadline(t time.Time) error {
	// TODO
	return nil
}

func (c *InboundConnection) SetReadDeadline(t time.Time) error {
	// TODO
	return nil
}

func (c *InboundConnection) SetWriteDeadline(t time.Time) error {
	// TODO
	return nil
}
