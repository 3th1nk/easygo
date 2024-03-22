package rabbitMQ

import "fmt"

type ChannelPoolBusyError struct {
	size int
}

func (this *ChannelPoolBusyError) Error() string {
	return fmt.Sprintf("channel-pool busy(%d)", this.size)
}

type ChannelOpenError struct {
	cause error
}

func (this *ChannelOpenError) Error() string {
	return fmt.Sprintf("open channel error: %v", this.cause)
}
