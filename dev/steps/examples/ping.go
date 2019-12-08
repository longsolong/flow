package examples

import (
	"fmt"
	"github.com/longsolong/flow/pkg/orchestration/request"
	"github.com/longsolong/flow/pkg/workflow/atom"
	"time"

	"github.com/longsolong/flow/pkg/workflow"
	"github.com/longsolong/flow/pkg/workflow/step"
	ping "github.com/sparrc/go-ping"
)

// Ping ...
type Ping struct {
	step.Step
	pinger *ping.Pinger
}

// NewPing ...
func NewPing(id, expansionDigest string) *Ping {
	p := &Ping{}
	p.SetID(atom.ID{
		ID:              id,
		ExpansionDigest: expansionDigest,
		Type:            atom.GenRunnableType(p, "dev/examples"),
	})
	return p
}

// Create ...
func (p *Ping) Create(ctx workflow.Context) error {
	req := ctx.Request().(*request.Request)

	pinger, err := ping.NewPinger(req.PrimaryRequestArgs["hostname"].(string))
	if err != nil {
		return err
	}
	pinger.OnRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %p: icmp_seq=%d time=%v ttl=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
	}
	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	pinger.Count = req.PrimaryRequestArgs["count"].(int)
	pinger.Interval = time.Second * time.Duration(req.PrimaryRequestArgs["interval"].(int))
	pinger.Timeout = time.Second * time.Duration(req.PrimaryRequestArgs["timeout"].(int))
	p.pinger = pinger
	return nil
}

// Run a ping
func (p *Ping) Run(ctx workflow.Context) (atom.Return, error) {
	ret := atom.Return{}
	p.pinger.Run()
	return ret, nil
}

// Stop run
func (p *Ping) Stop(ctx workflow.Context) error {
	p.pinger.Stop()
	return nil
}