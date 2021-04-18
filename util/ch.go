package util

// 上层调用需加锁，否则并发调用时会有panic可能
func CloseChannel(ch chan struct{}) {
	select {
	case <-ch:
		return
	default:
	}
	close(ch)
}

// 上层调用需加锁
func ChannelIsClosed(ch chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
	}
	return false
}
