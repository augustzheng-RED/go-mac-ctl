package main

import (
	"go-mac-ctl/cmd"
	"go-mac-ctl/internal/executor"
)

func main() {
	_ = executor.LoadDotEnv(".env")
	cmd.Execute()
}
