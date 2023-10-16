package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gerardus/cli"
	"gerardus/paths"
)

func usage(msg string, args ...any) {
	cli.StdErr(msg+"\n\n", args...)
	cmd, _ := cli.CommandByName("help")
	am, err := cmd.ArgsMap()
	if err != nil {
		cli.StdErr(err.Error())
	}
	_ = cli.ExecHelp(am)
	os.Exit(1)
}

func makeAbs(path string) (string, error) {
	absDir, err := paths.Absolute(path)
	if err != nil {
		err = fmt.Errorf("error converting to absolute path: %s; %w",
			path, err)
	}
	return absDir, err
}

func defaultSourceDir() string {
	dir := os.Getenv("GOROOT")
	if len(dir) > 0 {
		dir = filepath.Join(dir, "src")
	}
	return dir
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
