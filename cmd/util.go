package main

import (
	"os"

	"gerardus/options"
)

func usage(msg string, args ...any) {
	options.StdErr(msg+"\n\n", args...)
	options.StdErr("\tUsage: gerardus [-src=<source_dir>] [-out=<output_dir>] [-data=<sqlite_db>] run\n")
	os.Exit(1)
}

// printSymbolTypes is just an archetype for how to call SQLC generated funcs
//func printSymbolTypes() {
//	ds := gerardus.GetDataStore()
//	sts, err := ds.ListSymbolTypes(context.Background())
//	if err != nil {
//		ErrOut(err)
//	}
//	for i, st := range sts {
//		fmt.Printf("%d. id=%d, name=%s\n", i, st.ID, st.Name)
//	}
//	os.Exit(1)
//}
