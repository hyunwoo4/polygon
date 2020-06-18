package clique

import (
	"context"

	"github.com/0xPolygon/minimal/consensus"
	"github.com/0xPolygon/minimal/types"
)

// Clique is a consensus algorithm for the clique protocol
type Clique struct {
}

func Factory(ctx context.Context, config *consensus.Config) (consensus.Consensus, error) {
	c := &Clique{}
	return c, nil
}

// VerifyHeader verifies the header is correct
func (c *Clique) VerifyHeader(chain consensus.ChainReader, header *types.Header, uncle, seal bool) error {
	return nil
}

// Seal seals the block
func (c *Clique) Seal(chain consensus.ChainReader, block *types.Block, ctx context.Context) (*types.Block, error) {
	return nil, nil
}

// Close closes the connection
func (c *Clique) Close() error {
	return nil
}
