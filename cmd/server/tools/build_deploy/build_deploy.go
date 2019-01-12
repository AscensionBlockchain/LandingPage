package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	bindata "github.com/go-bindata/go-bindata"
)

const assetsFileName = "assets_bindata.go"

func main() {
	err := os.Chdir(filepath.Join("..", ".."))
	if err != nil {
		log.Fatalln("Unable to chdir to ../.. because:", err)
	}
	config := new(bindata.Config)
	config.Package = "main"
	config.Input = []bindata.InputConfig{
		bindata.InputConfig{
			Path:      "assets",
			Recursive: true,
		},
		bindata.InputConfig{
			Path:      "templates",
			Recursive: true,
		},
		bindata.InputConfig{
			Path: "config.toml",
		},
	}
	config.Output = filepath.Join("cmd", "server", assetsFileName)
	err = bindata.Translate(config)
	if err != nil {
		log.Fatalln("unable to create", assetsFileName, "because:", err)
	}

	fmt.Println("Successfully wrote", assetsFileName)
}
