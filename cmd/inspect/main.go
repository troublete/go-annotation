package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/troublete/go-annotation/annotation"
	"github.com/troublete/go-annotation/inspect"
)

func main() {
	root := flag.String("root", "./", "root path for inspection")
	flag.Parse()

	slog.Info("inspecting structure", "root", *root)

	if *root == "" {
		slog.Error("-root is required.")
		os.Exit(1)
	}

	types, err := inspect.FindAllTypes(*root)
	if err != nil {
		slog.Error("failed to find all types", "err", err)
		os.Exit(1)
	}

	funcs, err := inspect.FindAllFunctions(*root)
	if err != nil {
		slog.Error("failed to find all funcs", "err", err)
		os.Exit(1)
	}

	def, err := annotation.Read(types, funcs)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	c, err := json.MarshalIndent(def, "", "\t")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	fmt.Println(bytes.NewBuffer(c).String())
}
