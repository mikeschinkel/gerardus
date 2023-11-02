package persister

import (
	"context"
	"log/slog"

	"gerardus/collector"
	"gerardus/parser"
)

func Initialize(ctx context.Context, fp string, etss ...any) (err error) {
	slog.Info("Initializing persister")
	dataStore, err = getDataStore(fp)
	if err != nil {
		goto end
	}
	err = dataStore.InitializeDataStore(ctx)
	if err != nil {
		goto end
	}
	for _, ets := range etss {
		switch t := ets.(type) {
		case []collector.SymbolType:
			for _, st := range t {
				_, err = dataStore.UpsertSymbolType(ctx, UpsertSymbolTypeParams{
					ID:   int64(st.ID()),
					Name: st.Name(),
				})
				if err != nil {
					goto end
				}
			}
		case []parser.PackageType:
			for _, pt := range t {
				_, err = dataStore.UpsertPackageType(ctx, UpsertPackageTypeParams{
					ID:   int64(pt.ID()),
					Name: pt.Name(),
				})
				if err != nil {
					goto end
				}
			}
		default:
			panicf("Unexpected invalid EnumTypes: %#v.", ets)
		}
	}
end:
	return err
}
