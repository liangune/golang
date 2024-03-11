package kafka

import "sync"

const (
	ConsumerTypeNone      = 0
	ConsumerTypeGroup     = 1
	ConsumerTypePartition = 2
)

type ConsumerWorker struct {
	groupConsumer     *GroupConsumer
	partitionConsumer *PartitionConsumer
	listener          ConsumerListener
	errorListener     ConsumerErrorListener
	closer            chan struct{}
	closed            chan struct{}
	closerError       chan struct{}
	closedError       chan struct{}
	closeOnce         sync.Once
	consumerType      int
}

func NewConsumerWorker(cfg *ConsumerConfig) (*ConsumerWorker, error) {
	var err error
	var groupConsumer *GroupConsumer
	var partitionConsumer *PartitionConsumer
	consumerType := ConsumerTypeNone
	if cfg.AutoResetOffset == AutoResetOffsetTimestamp {
		consumerType = ConsumerTypePartition
		partitionConsumer, err = NewPartitionConsumer(cfg)
	} else {
		consumerType = ConsumerTypeGroup
		groupConsumer, err = NewGroupConsumer(cfg)
	}

	if err != nil {
		return nil, err
	}

	w := &ConsumerWorker{
		groupConsumer:     groupConsumer,
		partitionConsumer: partitionConsumer,
		listener:          cfg.Listener,
		errorListener:     cfg.ErrorListener,
		closer:            make(chan struct{}),
		closed:            make(chan struct{}),
		closerError:       make(chan struct{}),
		closedError:       make(chan struct{}),
		consumerType:      consumerType,
	}

	return w, nil
}

func (w *ConsumerWorker) Start() {
	if w.consumerType == ConsumerTypePartition {
		go w.RecvPartitionMessage()
	} else {
		go w.RecvGroupMessage()
	}

	go w.ReturnError()
}

func (w *ConsumerWorker) Close() (err error) {
	w.closeOnce.Do(func() {
		close(w.closerError)
		<-w.closedError
		close(w.closer)
		<-w.closed
		if w.consumerType == ConsumerTypePartition {
			err = w.partitionConsumer.Close()
		} else {
			err = w.groupConsumer.Close()
		}
	})
	return
}

func (w *ConsumerWorker) RecvGroupMessage() {
	defer close(w.closed)
	for {
		select {
		case msg := <-w.groupConsumer.Messages():
			ctx := &ConsumerMessageContext{
				groupId:   w.groupConsumer.GetGroupId(),
				topic:     w.groupConsumer.GetTopic(),
				partition: msg.Partition,
				offset:    msg.Offset,
				message:   msg.Value,
				timeStamp: msg.Timestamp,
				val:       msg,
			}
			w.listener(ctx)
			w.groupConsumer.Consumer.MarkOffset(msg, "")
		case <-w.closer:
			return
		}
	}
}

func (w *ConsumerWorker) RecvPartitionMessage() {
	defer close(w.closed)
	for {
		select {
		case msg := <-w.partitionConsumer.Messages():
			ctx := &ConsumerMessageContext{
				groupId:   w.partitionConsumer.GetGroupId(),
				topic:     w.partitionConsumer.GetTopic(),
				partition: msg.Partition,
				offset:    msg.Offset,
				message:   msg.Value,
				timeStamp: msg.Timestamp,
				val:       msg,
			}
			timestamp := msg.Timestamp.Unix()
			if timestamp >= w.partitionConsumer.OffsetBeginTimestamp && timestamp <= w.partitionConsumer.OffsetEndTimestamp {
				w.listener(ctx)
			}
		case <-w.closer:
			return
		}
	}
}

func (w *ConsumerWorker) ReturnError() {
	defer close(w.closedError)
	if w.consumerType == ConsumerTypePartition {
		for {
			select {
			case err := <-w.partitionConsumer.Errors():
				if w.errorListener != nil {
					w.errorListener(err.Err)
				}
			case <-w.closerError:
				return
			}
		}
	} else {
		for {
			select {
			case err := <-w.groupConsumer.Errors():
				if w.errorListener != nil {
					w.errorListener(err)
				}
			case <-w.closerError:
				return
			}
		}
	}

}
