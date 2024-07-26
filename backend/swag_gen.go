//go:build ignore

package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = filepath.Join(os.Getenv("HOME"), "go")
	}

	swagPath := filepath.Join(gopath, "bin", "swag")

	cmd := exec.Command(swagPath, "init")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error running swag init: %v", err)
	}

	log.Println("Swagger documentation generated successfully")
}