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

	Vertices    map[atom.ID]*Node // All vertices in the graph (node id -> node)
	VerticesMux *sync.RWMutex     // for access to vertices maps
}

// Node represents a single vertex within a Graph.
// Each node consists of a Payload (i.e. the data that the
// user cares about), a list of next and prev Nodes, and other
// information about the node such as the number of times it
// should be retried on error. Next defines all the out edges
// from Node, and Prev defines all the in edges to Node.
type Node struct {
	Datum   atom.Atom         // Data stored at this Node
	Next    map[atom.ID]*Node // out edges ( node id -> Node )
	Prev    map[atom.ID]*Node // in edges ( node id -> Node )
	EdgeMux *sync.RWMutex     // for access to vertices maps

	Name          string        // the name of the node
	Retry         uint          // the number of times to retry a node
	RetryWait     time.Duration // the time, in seconds, to sleep between retries
	SequenceID    atom.ID       // ID for first node in sequence
	SequenceRetry uint          // Number of times to retry a sequence. Only set for first node in sequence.
}

// NewDAG ...
func NewDAG(name string, version int) *DAG {
	return &DAG{
		Name:        name,
		Version:     version,
		Vertices:    make(map[atom.ID]*Node),
		VerticesMux: &sync.RWMutex{},
	}
}

// NewNode ...
func NewNode(a atom.Atom, name string, retry uint, retryWait time.Duration) *Node {
	return &Node{
		Datum:   a,
		Next:    make(map[atom.ID]*Node),
		Prev:    make(map[atom.ID]*Node),
		EdgeMux: &sync.RWMutex{},

		Name:      name,
		Retry:     retry,
		RetryWait: retryWait,
	}
}

// AddNode ...
func (g *DAG) AddNode(node *Node) error {
	g.VerticesMux.Lock()
	defer g.VerticesMux.Unlock()
	if _, ok := g.Vertices[node.Datum.ID()]; ok {
		return workflow.ErrAlreadyRegisteredNode
	}
	g.Vertices[node.Datum.ID()] = node
	return nil
}

// GetNode ...
func (g *DAG) GetNode(atomID atom.ID) (*Node, error) {
	g.VerticesMux.RLock()
	defer g.VerticesMux.RUnlock()
	if node, ok := g.Vertices[atomID]; ok {
		return node, nil
	}
	return nil, workflow.ErrNotRegisteredNode
}

// SetUpstream ...
func (n *Node) SetUpstream(upstream *Node) error {
	n.EdgeMux.Lock()
	defer n.EdgeMux.Unlock()
	if _, ok := n.Prev[upstream.Datum.ID()]; ok {
		return workflow.ErrAlreadyRegisteredUpstream
	}
	if _, ok := upstream.Next[n.Datum.ID()]; ok {
		return workflow.ErrAlreadyRegisteredDownstream
	}
	n.Prev[upstream.Datum.ID()] = upstream
	upstream.Next[n.Datum.ID()] = n
	return nil
}
