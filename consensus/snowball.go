package consensus

import (
	"avalanche-consensus/chain"
	"context"
	"errors"
)

type Consensus struct {
	config     Config
	preference chain.DataType
	confidence int
	isRunning  bool
}

func NewConcensus(conf Config, preference chain.DataType) (*Consensus, error) {
	err := conf.Verify()
	if err != nil {
		return nil, err
	}

	consensus := &Consensus{
		config:     conf,
		preference: preference,
		confidence: 0,
		isRunning:  false,
	}

	return consensus, nil
}

func (c *Consensus) Run(ctx context.Context, index int, chain *chain.Blockchain) error {
	if c.isRunning {
		return errors.New("Consensus is running")
	}

	c.isRunning = true
	c.confidence = 1

	for c.confidence < c.config.Beta {
		preferenceFromKpeer, err := chain.GetDataFromKRandomBlock(ctx, index, c.config.K)
		if err != nil {
			return err
		}
		preference, frequent, err := c.getMostPreference(preferenceFromKpeer)
		if err != nil {
			return errors.New(err.Error() + ", unable to get the most frequent")
		}

		if frequent >= c.config.Alphal {
			oldPrefer := c.preference
			c.preference = preference
			err := chain.Set(index, preference)
			if err != nil {
				return errors.New(err.Error() + ", error when update the preference")
			}

			if preference == oldPrefer {
				c.confidence++
			} else {
				c.confidence = 1
			}

		} else {
			c.confidence = 0
		}
	}
	c.isRunning = false

	return nil
}

func (c *Consensus) getMostPreference(preferences []chain.DataType) (chain.DataType, int, error) {
	if len(preferences) == 0 {
		return 0, 0, errors.New("the preferences is empty")
	}
	countMap := map[chain.DataType]int{}
	maxCount := 0
	var prefer chain.DataType
	for _, p := range preferences {
		countMap[p]++
		if countMap[p] > maxCount {
			maxCount = countMap[p]
			prefer = p
		}
	}
	return prefer, maxCount, nil
}
