package main

import (
	"fmt"
	"os"

	"github.com/zexot-com/commitlint/internal/cmd"
)

var errExitCode = 1

func main() {
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(errExitCode)
	}
}
