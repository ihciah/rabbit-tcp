package connection

import (
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
}

func NewOutboundConnectionWithID(conn net.Conn, connectionID uint32, sendQueue chan<- block.Block) Connection {
	c := OutboundConnection{
		BaseConnection: BaseConnection{
			connectionID:     connectionID,
			sendQueue:        sendQueue,
			recvQueue:        make(chan block.Block, RecvQueueSize),
			orderedRecvQueue: make(chan block.Block, OrderedRecvQueueSize),
		},
		Conn: conn,
	}
	c.blockProcessor = newBlockProcessor(&c)
	go c.RecvRelay()
	go c.SendRelay()
	return &c
}

func (oc *OutboundConnection) RecvRelay() {
	recvBuffer := make([]byte, OutboundRecvBuffer)
	for {
		n, err := oc.Conn.Read(recvBuffer)
		if err == nil {
			oc.SendData(recvBuffer[:n])
		} else if err == io.EOF {
			oc.SendDisconnect()
			return
		} else {
			// TODO: error handle
		}
	}
}

func (oc *OutboundConnection) SendRelay() {
	for {
		blk := <-oc.orderedRecvQueue
		var err error
		switch blk.Type {
		case block.BLOCK_TYPE_DATA:
			_, err = oc.Conn.Write(blk.BlockData)
		case block.BLOCK_TYPE_DISCONNECT:
			err = oc.Conn.Close()
		}
		if err != nil {
			// TODO: error handle
		}
	}
}
