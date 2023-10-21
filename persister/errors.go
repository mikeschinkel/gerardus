package persister

import (
	"gerardus/serr"
)

var (
	errFailedToInsertSpec    = serr.New("failed to insert Spec")
	errFailedWhilePersisting = serr.New("failed while persisting")
)
