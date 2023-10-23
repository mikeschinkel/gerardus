package channels

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Stage func() error

type Pipeline struct {
	group  *errgroup.Group
	cancel context.CancelFunc
	stages []Stage
}

func NewPipeline(ctx context.Context) *Pipeline {
	group, ctx := errgroup.WithContext(ctx)
	ctx, cancel := context.WithCancel(ctx)
	return &Pipeline{
		group:  group,
		cancel: cancel,
		stages: make([]Stage, 0),
	}
}

func (pl *Pipeline) AddStage(stage Stage) {
	pl.stages = append(pl.stages, CancelOnErr(pl.cancel, stage))
}

func (pl *Pipeline) Go() error {
	for i := len(pl.stages) - 1; i >= 0; i-- {
		// Call in reverse order do the downstream function will be ready before the
		// upstream function starts.
		pl.group.Go(pl.stages[i])
	}
	return pl.group.Wait()
}
