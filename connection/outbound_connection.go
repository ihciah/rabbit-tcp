package connection

import (
	"context"
	"github.com/ihciah/rabbit-tcp/block"
	"io"
	"net"
)

const (
	OutboundRecvBuffer = 2048
)

type OutboundConnection struct {
	BaseConnection
	net.Conn
	cancel context.CancelFunc
}

func NewOutboundConnectionWithID(conn net.Conn, connectionID uint32, sendQueue chan<- block.Block) Connection {
	c := OutboundConnection{
		BaseConnection: BaseConnection{
			connectionID:     connectionID,
			ok:               true,
			sendQueue:        sendQueue,
			recvQueue:        make(chan block.Block, RecvQueueSize),
			orderedRecvQueue: make(chan block.Block, OrderedRecvQueueSize),
		},
		Conn: conn,
	}
	c.blockProcessor = newBlockProcessor(&c)
	// Will stop relay when Close
	var ctx context.Context
	ctx, c.cancel = context.WithCancel(context.Background())

	go c.RecvRelay(ctx)
	go c.SendRelay(ctx)
	return &c
}

// real connection -> ConnectionPool's SendQueue -> TunnelPool
func (oc *OutboundConnection) RecvRelay(ctx context.Context) {
	recvBuffer := make([]byte, OutboundRecvBuffer)
	for {
		n, err := oc.Conn.Read(recvBuffer)
		if err == nil {
			oc.SendData(recvBuffer[:n])
		} else if err == io.EOF {
			oc.ok = false
			oc.Conn.Close()
			oc.SendDisconnect()
			return
		} else {
			// TODO: error handle
		}
	}
}

// orderedRecvQueue -> real connection
func (oc *OutboundConnection) SendRelay(ctx context.Context) {
	for {
		select {
		case blk := <-oc.orderedRecvQueue:
			var err error
			switch blk.Type {
			case block.BLOCK_TYPE_DATA:
				if oc.ok {
					_, err = oc.Conn.Write(blk.BlockData)
				}
			case block.BLOCK_TYPE_DISCONNECT:
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
		case <-ctx.Done():
			// TODO: thread safe
			if oc.ok {
				oc.ok = false
				oc.Conn.Close()
			}
			return
		}
	}
}

func (oc *OutboundConnection) CancelDaemon() {
	oc.BaseConnection.CancelDaemon()
	if oc.cancel != nil {
		oc.cancel()
	}
}
