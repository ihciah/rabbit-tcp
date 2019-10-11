package connection

const (
	OrderedRecvQueueSize    = 24        // OrderedRecvQueue channel cap
	RecvQueueSize           = 24        // RecvQueue channel cap
	OutboundRecvBuffer      = 16 * 1024 // 16K receive buffer for Outbound Connection
	OutboundBlockTimeoutSec = 3         // Wait the period and check exit signal
	PacketWaitTimeoutSec    = 7         // If block processor is waiting for a "hole", and no packet comes within this limit, the Connection will be closed
)
