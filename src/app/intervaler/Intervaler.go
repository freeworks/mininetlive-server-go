package intervaler

import (
	easemob "app/easemob"
	logger "app/logger"
	"time"

	"github.com/coopernurse/gorp"
	cache "github.com/patrickmn/go-cache"
)

// Intervaler 是个接口用来让调用者自定义poller轮询时间间隔
type Intervaler interface {
	Interval() time.Duration
}

// IntervalerFunc 用来将 func() time.Duration 转化成 Intervaler
type IntervalerFunc func() time.Duration

func (intervalerFunc IntervalerFunc) Interval() time.Duration {
	return intervalerFunc()
}

type Poller struct {
	//要执行的方法
	do func() error
	//用于调用者传递停止信号
	cancle chan int
	//下次调用的时间间隔
	nextInterval Intervaler
}

// Poll 轮询
func (poller *Poller) Poll() {
	for {
		select {
		case <-poller.cancle:
			return
		case <-time.After(poller.nextInterval.Interval()):
			go func() {
				if err := poller.do(); err != nil {
					logger.Error("Poll poller.go: polling method returns a error: %v", err)
					// 或者结束整个轮询
					// poller.Cancel()
				}
			}()
		}
	}
}

// Cancel 向 cancel 发送信号
func (poller *Poller) Cancel() {
	println("Polling stopped")
	poller.cancle <- 1
}

// NewPoller 创建一个新的 Poller
func NewPoller(intervaler Intervaler, do func() error) *Poller {
	return &Poller{do: do, cancle: make(chan int), nextInterval: intervaler}
}

func PollGroupOnlineUser(c *cache.Cache, dbmap *gorp.DbMap) {
	base := time.Second * 0
	interval := IntervalerFunc(func() time.Duration {
		next := base
		base += time.Second * 30
		return next
	})
	poller := NewPoller(interval,
		func() error {
			easemob.GetGroupOnlineUserCount(c, dbmap)
			return nil
		})
	poller.Poll()
}
