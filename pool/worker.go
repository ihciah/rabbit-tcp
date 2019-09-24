package pool

import (
	"bytes"
	"github.com/ihciah/rabbit-tcp/block"
	"io"
	"net"
)

type Worker interface {
	net.Conn
	StartSend(<-chan block.Block)
	StartRecv(chan<- block.Block)
	Exit()
}

type BasicWorker struct {
	net.Conn
	sendExit chan bool
	recvExit chan bool
}

func NewBasicWorker(conn net.Conn) Worker {
	return &BasicWorker{
		Conn:     conn,
		sendExit: make(chan bool, 1),
		recvExit: make(chan bool, 1),
	}
}

func (worker *BasicWorker) StartSend(input <-chan block.Block) {
	for {
		select {
		case <-worker.sendExit:
			return
		case blk := <-input:
			reader := bytes.NewReader(blk.Pack())
			io.Copy(worker.Conn, reader)
			// TODO: error handle
		}
	}
}

func (worker *BasicWorker) StartRecv(output chan<- block.Block) {
	for {
		select {
		case <-worker.recvExit:
			return
		default:
			blk := block.NewBlockFromReader(worker.Conn)
			output <- *blk
		}
	}
}

func (worker *BasicWorker) Exit() {
	worker.sendExit <- true
	worker.recvExit <- true
}
