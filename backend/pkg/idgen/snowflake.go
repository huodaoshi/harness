package idgen

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
)

// IDGenerator wraps a snowflake node to produce unique string IDs.
type IDGenerator struct {
	node *snowflake.Node
}

// NewIDGenerator creates an IDGenerator using the given node ID (0–1023).
func NewIDGenerator(nodeID int64) (*IDGenerator, error) {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, fmt.Errorf("idgen: snowflake: failed to create node: %w", err)
	}
	return &IDGenerator{node: node}, nil
}

// Generate produces a new unique ID as a string.
func (g *IDGenerator) Generate() string {
	return g.node.Generate().String()
}
