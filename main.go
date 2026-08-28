package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/willGabriell/sqlplan/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		if errors.Is(err, cmd.ErrUsage) {
			fmt.Fprintln(os.Stderr, err)
			flag.Usage()
			os.Exit(2)
		}
		fmt.Fprintln(os.Stderr, "sqlplan:", err)
		os.Exit(1)
	}
}
