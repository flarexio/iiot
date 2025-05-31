package stdio

import (
	"bufio"
	"fmt"
	"os/exec"
	"testing"
	"time"
)

func TestManualExecution(t *testing.T) {
	cmd := exec.Command("modbus_tool")

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()

	cmd.Start()

	req := `{"method": "driver.schema"}`
	stdin.Write([]byte(req + "\n"))

	scanner := bufio.NewScanner(stdout)
	if scanner.Scan() {
		fmt.Println("Response 1: " + scanner.Text())
	}

	time.Sleep(10 * time.Second)

	stdin.Write([]byte(req + "\n"))

	scanner = bufio.NewScanner(stdout)
	if scanner.Scan() {
		fmt.Println("Response 2: " + scanner.Text())
	} else {
		fmt.Println("No response received")
	}
}
