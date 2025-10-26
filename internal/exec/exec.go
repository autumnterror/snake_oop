package exec

import (
	"os"
	"os/exec"
)

func Clean() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
