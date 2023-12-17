package persister

import (
	"context"
	"log/slog"

	"github.com/mikeschinkel/gerardus/collector"
	"github.com/mikeschinkel/gerardus/options"
	"github.com/mikeschinkel/gerardus/parser"
)

func Initialize(ctx context.Context, fp string, types ...any) (ds DataStore, err error) {

	var q DataStoreQueries

	slog.Info("Initializing persister")

	ds = NewSqliteDataStore(fp)

	err = ds.Initialize(ctx)
	if err != nil {
		goto end
	}
	q = ds.Queries()
	for _, typ := range types {
		switch t := typ.(type) {
		case []collector.SymbolType:
			for _, st := range t {
				_, err = q.UpsertSymbolType(ctx, UpsertSymbolTypeParams{
					ID:   int64(st.ID()),
					Name: st.Name(),
				})
				if err != nil {
					goto end
				}
			}
		case []parser.PackageType:
			for _, pt := range t {
				_, err = q.UpsertPackageType(ctx, UpsertPackageTypeParams{
					ID:   int64(pt.ID()),
					Name: pt.Name(),
				})
				if err != nil {
					goto end
				}
			}
		default:
			panicf("Unexpected invalid EnumTypes: %#v.", typ)
		}
	}
end:
	if err != nil {
		err = ErrFailedToInitDataStore.Err(err, "data_file", options.DataFile())
	}
	return ds, err
}
