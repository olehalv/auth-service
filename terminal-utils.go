package main

import (
	"os"
	"os/exec"
)

func clearConsole() {
	clears := []string{"clear", "cls"}
	for _, s := range clears {
		cmd := exec.Command(s)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}
}
