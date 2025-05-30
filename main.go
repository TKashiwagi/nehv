package main

import (
	"fmt"
	"os"

	"configure/cmd"
)

func main() {
	// configureコマンドを直接実行
	cmdArgs := os.Args
	if len(cmdArgs) == 1 || (len(cmdArgs) > 1 && cmdArgs[1] == "configure") {
		// configureコマンドのRunを直接呼び出す
		cmd.ConfigureCmdRun(nil, nil)
	} else {
		if err := cmd.Execute(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}
