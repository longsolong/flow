package dag

import (
	"github.com/longsolong/flow/pkg/workflow"
	"github.com/longsolong/flow/pkg/workflow/atom"
	"sync"
	"time"
)

// DAG represents a directed acyclic graph.
type DAG struct {
	Name    string // Name of the Graph
	Version int    // Version of the Graph

	Vertices    map[atom.AtomID]*Node // All vertices in the graph (node id -> node)
	VerticesMux *sync.RWMutex         // for access to vertices maps
}

// Node represents a single vertex within a Graph.
// Each node consists of a Payload (i.e. the data that the
// user cares about), a list of next and prev Nodes, and other
// information about the node such as the number of times it
// should be retried on error. Next defines all the out edges
// from Node, and Prev defines all the in edges to Node.
type Node struct {
	Datum   atom.Atom             // Data stored at this Node
	Next    map[atom.AtomID]*Node // out edges ( node id -> Node )
	Prev    map[atom.AtomID]*Node // in edges ( node id -> Node )
	EdgeMux *sync.RWMutex         // for access to vertices maps

	Name          string        // the name of the node
	Retry         uint          // the number of times to retry a node
	RetryWait     time.Duration // the time, in seconds, to sleep between retries
	SequenceID    atom.AtomID   // AtomID for first node in sequence
	SequenceRetry uint          // Number of times to retry a sequence. Only set for first node in sequence.
}

// NewDAG ...
func NewDAG(name string, version int) *DAG {
	return &DAG{
		Name:        name,
		Version:     version,
		Vertices:    make(map[atom.AtomID]*Node),
		VerticesMux: &sync.RWMutex{},
	}
}

// NewNode ...
func NewNode(a atom.Atom, name string, retry uint, retryWait time.Duration) *Node {
	return &Node{
		Datum:   a,
		Next:    make(map[atom.AtomID]*Node),
		Prev:    make(map[atom.AtomID]*Node),
		EdgeMux: &sync.RWMutex{},

		Name:      name,
		Retry:     retry,
		RetryWait: retryWait,
	}
}

// MustAddNode ...
func (g *DAG) MustAddNode(node *Node) {
	g.VerticesMux.Lock()
	defer g.VerticesMux.Unlock()
	if _, ok := g.Vertices[node.Datum.AtomID()]; ok {
		panic(workflow.ErrAlreadyRegisteredNode)
	}
	g.Vertices[node.Datum.AtomID()] = node
}

// GetNode ...
func (g *DAG) GetNode(atomID atom.AtomID) (*Node, error) {
	g.VerticesMux.RLock()
	defer g.VerticesMux.RUnlock()
	if node, ok := g.Vertices[atomID]; ok {
		return node, nil
	}
	return nil, workflow.ErrNotRegisteredNode
}

// MustGetNode ...
func (g *DAG) MustGetNode(atomID atom.AtomID) *Node {
	g.VerticesMux.RLock()
	defer g.VerticesMux.RUnlock()
	if node, ok := g.Vertices[atomID]; ok {
		return node
	}
	panic(workflow.ErrNotRegisteredNode)
}

// SetUpstream ...
func (n *Node) SetUpstream(upstream *Node) error {
	n.EdgeMux.Lock()
	defer n.EdgeMux.Unlock()
	if _, ok := n.Prev[upstream.Datum.AtomID()]; ok {
		return workflow.ErrAlreadyRegisteredUpstream
	}
	if _, ok := upstream.Next[n.Datum.AtomID()]; ok {
		return workflow.ErrAlreadyRegisteredDownstream
	}
	n.Prev[upstream.Datum.AtomID()] = upstream
	upstream.Next[n.Datum.AtomID()] = n
	return nil
}

// Upstream ...
func (n *Node) Upstream() map[atom.AtomID]*Node {
	n.EdgeMux.RLock()
	defer n.EdgeMux.RUnlock()
	return n.Prev
}

// Downstream ...
func (n *Node) Downstream() map[atom.AtomID]*Node {
	n.EdgeMux.RLock()
	defer n.EdgeMux.RUnlock()
	return n.Next
}
