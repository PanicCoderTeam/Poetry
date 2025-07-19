package game

import (
	"time"

	"github.com/lonng/nano/component"
	"github.com/lonng/nano/pipeline"
	"github.com/lonng/nano/scheduler"
	"github.com/lonng/nano/session"
)

// 流量统计
type Stats struct {
	// 继承 nano 组件，拥有完整的生命周期
	component.Base
	// 组件初始化完成后，做一些定时任务
	timer *scheduler.Timer
	// 出口流量统计
	outboundBytes int
	// 入口流量统计
	inboundBytes int
}

// 统计出口流量，会定义到 nano 的 pipeline
func (stats *Stats) Outbound(s *session.Session, msg *pipeline.Message) error {
	stats.outboundBytes += len(msg.Data)
	return nil
}

// 统计入口流量，会定义到 nano 的 pipeline
func (stats *Stats) Inbound(s *session.Session, msg *pipeline.Message) error {
	stats.inboundBytes += len(msg.Data)
	return nil
}

// 组件初始化完成后，会调用
// 每分钟会打印下出口与入口的流量
func (stats *Stats) AfterInit() {
	stats.timer = scheduler.NewTimer(time.Minute, func() {
		println("OutboundBytes", stats.outboundBytes)
		println("InboundBytes", stats.outboundBytes)
	})
}

func (st *Stats) Nil(s *session.Session, msg []byte) error {
	return nil
}
