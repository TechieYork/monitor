package queue

import (
	"errors"
	"time"

	"github.com/DarkMetrix/monitor/agent/src/protocol"
)

type TransferQueue struct {
	queueChannel chan *protocol.Proto
}

func NewTransferQueue(bufferSize int) *TransferQueue {
	return &TransferQueue{
		queueChannel: make(chan *protocol.Proto, bufferSize),
	}
}

func (queue *TransferQueue) Push(item* protocol.Proto) error {
	select {
	case queue.queueChannel <- item:
		return nil
	default:
		return errors.New("Channel full")
	}
}

func (queue * TransferQueue) Pop(ms time.Duration) (*protocol.Proto, error) {
	select {
	case item, ok := <- queue.queueChannel:
		if !ok {
			return nil, errors.New("Channel closed!")
		}

		return item, nil
	case <- time.After(ms):
		return nil, errors.New("Channel pop timeout!")
	}
}


