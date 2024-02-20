package main

import (
	"credo/config"
	"credo/logger"
	"credo/modules"
	"credo/storage"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	logger := logger.Get()

	if len(os.Args) <= 1 {
		logger.Println("Usage: \n./credo moduleName [args...]")
		return
	}

	_, moduleName, args := os.Args[0], os.Args[1], os.Args[2:]

	module := modules.Modules[moduleName]()

	params := modules.Parameters{
		Env: map[string]string{},
	}

	for i := 0; i < len(args)/2; i += 2 {
		params.Env[args[i]] = args[i+1]
	}

	store := storage.FileStorage{
		Filename: "credospell.yaml",
	}

	prevFile := store.Read()

	config, err := config.FromFile(prevFile)

	// TODO: Spell from Params
	result := module.BareRun(&config, &params)

	module.Commit(&config, result)

	marshal, err := yaml.Marshal(config)
	if err != nil {
		logger.Fatal("Can't marshal")
	}

	store.Write(marshal)
}
