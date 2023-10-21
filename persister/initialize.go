package persister

import (
	"context"
	"log/slog"
)

func Initialize[STS []ST, ST symbolType](ctx context.Context, fp string, sts STS) (err error) {
	slog.Info("Initializing persister")
	dataStore, err = getDataStore(fp)
	if err != nil {
		goto end
	}
	err = dataStore.InitializeDataStore(ctx)
	if err != nil {
		goto end
	}
	for _, st := range sts {
		_, err = dataStore.UpsertSymbolType(ctx, UpsertSymbolTypeParams{
			ID:   int64(st.ID()),
			Name: st.Name(),
		})
		if err != nil {
			break
		}
	}
end:
	return err
}
