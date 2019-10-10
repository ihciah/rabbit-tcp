package connection

import (
	"context"
	"fmt"
	"github.com/ihciah/rabbit-tcp/block"
	"github.com/ihciah/rabbit-tcp/logger"
	"io"
	"net"
)

const (
	OutboundRecvBuffer = 16 * 1024 // 16K
)

type OutboundConnection struct {
	baseConnection
	net.Conn
	ctx    context.Context
	cancel context.CancelFunc
}

func NewOutboundConnection(connectionID uint32, sendQueue chan<- block.Block, ctx context.Context, removeFromPool context.CancelFunc) Connection {
	c := OutboundConnection{
		baseConnection: baseConnection{
			blockProcessor:   newBlockProcessor(ctx, removeFromPool),
			connectionID:     connectionID,
			ok:               false,
			sendQueue:        sendQueue,
			recvQueue:        make(chan block.Block, RecvQueueSize),
			orderedRecvQueue: make(chan block.Block, OrderedRecvQueueSize),
			logger:           logger.NewLogger(fmt.Sprintf("[OutboundConnection-%d]", connectionID)),
		},
		ctx: ctx,
	}
	c.logger.Infof("OutboundConnection %d created.\n", connectionID)
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
			oc.logger.Debugln("EOF received from outbound connection.")
			oc.ok = false
			oc.SendDisconnect()
			oc.Conn.Close()
			// TODO: error handle
			return
		} else {
			oc.logger.Errorf("Error when relay outbound connection: %v\n.", err)
			oc.ok = false
			oc.SendDisconnect()
			oc.Conn.Close()
			// TODO: error handle
			return
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
			case block.TypeConnect:
				// Will do nothing!
				continue
			case block.TypeData:
				oc.logger.Debugln("Received DATA block.")
				if oc.ok {
					_, err = oc.Conn.Write(blk.BlockData)
				}
			case block.TypeDisconnect:
				oc.logger.Debugln("Received DISCONNECT block.")
				if oc.ok {
					oc.ok = false
					err = oc.Conn.Close()
				}
			}
			if err != nil {
				// TODO: error handle
				oc.ok = false
				err = oc.Conn.Close()
			}
		case <-oc.ctx.Done():
			// TODO: error handle
			oc.ok = false
			oc.Conn.Close()
			return
		}
	}
}

func (oc *OutboundConnection) RecvBlock(blk block.Block) {
	if blk.Type == block.TypeConnect {
		address := string(blk.BlockData)
		go oc.connect(address)
	}
	oc.recvQueue <- blk
}

func (oc *OutboundConnection) connect(address string) {
	oc.logger.Debugln("Received CONNECTION block.")
	if oc.ok || oc.Conn != nil {
		return
	}
	rawConn, err := net.Dial("tcp", address)
	if err == nil {
		oc.logger.Infof("Dail to %s successfully.\n", address)
		oc.Conn = rawConn
		oc.ok = true
		go oc.RecvRelay()
		go oc.SendRelay()
	} else {
		oc.logger.Warnf("Error when dail to %s: %v.\n", address, err)
		oc.SendDisconnect()
	}
}
