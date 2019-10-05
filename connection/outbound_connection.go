package connection

import (
	"context"
	"fmt"
	"github.com/ihciah/rabbit-tcp/block"
	"io"
	"log"
	"net"
	"os"
)

const (
	OutboundRecvBuffer = 2048
)

type OutboundConnection struct {
	baseConnection
	net.Conn
	ctx    context.Context
	cancel context.CancelFunc
}

func NewOutboundConnection(connectionID uint32, sendQueue chan<- block.Block, ctx context.Context, cancel context.CancelFunc) Connection {
	c := OutboundConnection{
		baseConnection: baseConnection{
			blockProcessor:   newBlockProcessor(ctx),
			connectionID:     connectionID,
			ok:               false,
			sendQueue:        sendQueue,
			recvQueue:        make(chan block.Block, RecvQueueSize),
			orderedRecvQueue: make(chan block.Block, OrderedRecvQueueSize),
			logger:           log.New(os.Stdout, fmt.Sprintf("[OutboundConnection-%d]", connectionID), log.LstdFlags),
		},
		ctx:    ctx,
		cancel: cancel,
	}
	return &c
}

// real connection -> ConnectionPool's SendQueue -> TunnelPool
func (oc *OutboundConnection) RecvRelay() {
	recvBuffer := make([]byte, OutboundRecvBuffer)
	for {
		n, err := oc.Conn.Read(recvBuffer)
		if err == nil {
			oc.sendData(recvBuffer[:n])
		} else if err == io.EOF {
			oc.ok = false
			oc.Conn.Close()
			oc.SendDisconnect()
			return
		} else {
			oc.logger.Printf("Error when relay outbound connection: %v\n.", err)
			// TODO: error handle
		}
		select {
		case <-oc.ctx.Done():
			return
		default:
			continue
		}
	}
}

// orderedRecvQueue -> real connection
func (oc *OutboundConnection) SendRelay() {
	for {
		select {
		case blk := <-oc.orderedRecvQueue:
			var err error
			switch blk.Type {
			case block.BLOCK_TYPE_CONNECT:
				// Will do nothing!
				continue
			case block.BLOCK_TYPE_DATA:
				oc.logger.Println("Received DATA block.")
				if oc.ok {
					_, err = oc.Conn.Write(blk.BlockData)
				}
			case block.BLOCK_TYPE_DISCONNECT:
				oc.logger.Println("Received DISCONNECT block.")
				if oc.ok {
					oc.ok = false
					err = oc.Conn.Close()
				}
			}
			if err != nil {
				// TODO: error handle
				// TODO: thread safe
				if oc.ok {
					oc.ok = false
					err = oc.Conn.Close()
				}
			}
		case <-oc.ctx.Done():
			// TODO: thread safe
			if oc.ok {
				oc.ok = false
				oc.Conn.Close()
			}
			return
		}
	}
}

func (oc *OutboundConnection) RecvBlock(blk block.Block) {
	if blk.Type == block.BLOCK_TYPE_CONNECT {
		address := string(blk.BlockData)
		go oc.connect(address)
	}
	oc.recvQueue <- blk
}

func (oc *OutboundConnection) connect(address string) {
	oc.logger.Println("Received CONNECTION block.")
	if oc.ok || oc.Conn != nil {
		return
	}
	rawConn, err := net.Dial("tcp", address)
	oc.logger.Printf("Dail to %s in error: %v.\n", address, err)
	if err == nil {
		oc.Conn = rawConn
		oc.ok = true
		go oc.RecvRelay()
		go oc.SendRelay()
	}
}
