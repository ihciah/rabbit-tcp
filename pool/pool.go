package pool

func NewPool(threadsWait []uint, timeoutSec uint) *Pool{
	// TODO: Start pool manager
	return &Pool{
		threadsWait: threadsWait,
		timeoutSec:  timeoutSec,
	}
}

type Pool struct {
	threadsWait []uint
	timeoutSec uint
}

func (p *Pool) SendBlock() {

}

func (p *Pool) RecvBlock() {

}