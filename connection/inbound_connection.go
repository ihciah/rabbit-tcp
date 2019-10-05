package connection

import (
	"fmt"
	"github.com/ihciah/rabbit-tcp/block"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

type InboundConnection struct {
	baseConnection
	dataBuffer LoopByteBuffer
}

func NewInboundConnection(sendQueue chan<- block.Block) Connection {
	connectionID := rand.Uint32()
	return NewInboundConnectionWithID(connectionID, sendQueue)
}

func NewInboundConnectionWithID(connectionID uint32, sendQueue chan<- block.Block) Connection {
	c := InboundConnection{
		baseConnection: baseConnection{
			blockProcessor:   newBlockProcessor(),
			connectionID:     connectionID,
			ok:               true,
			sendQueue:        sendQueue,
			recvQueue:        make(chan block.Block, RecvQueueSize),
			orderedRecvQueue: make(chan block.Block, OrderedRecvQueueSize),
			logger:           log.New(os.Stdout, fmt.Sprintf("[InboundConnection%d]", connectionID), log.LstdFlags),
		},
		dataBuffer: NewLoopBuffer(BlockMaxSize),
	}
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
		if !c.ok {
			return 0, io.EOF
		}
		blk := <-c.orderedRecvQueue
		switch blk.Type {
		case block.BLOCK_TYPE_DISCONNECT:
			c.ok = false
			return 0, io.EOF
		case block.BLOCK_TYPE_DATA:
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
			case block.BLOCK_TYPE_DISCONNECT:
				c.ok = false
				return readN, nil
			case block.BLOCK_TYPE_DATA:
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
	return nil
}

func (c *InboundConnection) LocalAddr() net.Addr {
	return nil
}

func (c *InboundConnection) RemoteAddr() net.Addr {
	return nil
}

func (c *InboundConnection) SetDeadline(t time.Time) error {
	return nil
}

func (c *InboundConnection) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *InboundConnection) SetWriteDeadline(t time.Time) error {
	return nil
}
