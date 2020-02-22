package examples

import (
	"context"
	"fmt"
	"github.com/longsolong/flow/pkg/infra"
	"github.com/longsolong/flow/pkg/orchestration/request"
	"github.com/longsolong/flow/pkg/workflow/atom"
	flowcontext "github.com/longsolong/flow/pkg/workflow/context"
	"time"

	"github.com/longsolong/flow/pkg/workflow/step"
	ping "github.com/sparrc/go-ping"
)

//go:generate genatom -type=Ping

// Ping ...
type Ping struct {
	step.Step
	pinger *ping.Pinger
}

// NewPing ...
func NewPing(id, expansionDigest string) *Ping {
	p := &Ping{}
	p.ID = id
	p.ExpansionDigest = expansionDigest
	return p
}

// Create ...
func (p *Ping) Create(ctx context.Context, req *request.Request) error {
	var logger *infra.Logger
	if l := ctx.Value(flowcontext.FlowContextKey("logger")); l != nil {
		logger = l.(*infra.Logger)
	}

	pinger, err := ping.NewPinger(req.RequestArgs["hostname"].(string))
	if err != nil {
		return err
	}
	pinger.OnRecv = func(pkt *ping.Packet) {
		logger.Log.Info(fmt.Sprintf("%d bytes from %p: icmp_seq=%d time=%v ttl=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl))
	}
	pinger.OnFinish = func(stats *ping.Statistics) {
		logger.Log.Info(fmt.Sprintf("ping statistics %v", stats.Addr))
		logger.Log.Info(fmt.Sprintf("%d packets transmitted, %d packets received, %v%% packet loss",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss))
		logger.Log.Info(fmt.Sprintf("round-trip min/avg/max/stddev = %v/%v/%v/%v",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt))
	}

	pinger.Count = int(req.RequestArgs["count"].(float64))
	pinger.Interval = time.Second * time.Duration(req.RequestArgs["interval"].(float64))
	pinger.Timeout = time.Second * time.Duration(req.RequestArgs["timeout"].(float64))
	p.pinger = pinger
	return nil
}

// Run a ping
func (p *Ping) Run(ctx context.Context) (atom.Return, error) {
	ret := atom.Return{}
	p.pinger.Run()
	return ret, nil
}

// Stop run
func (p *Ping) Stop(ctx context.Context) error {
	p.pinger.Stop()
	return nil
}