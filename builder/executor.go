package builder

import (
	"bufio"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/sirupsen/logrus"
	"io"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

func runCommand(command string) (err error) {
	cmd := exec.Command("cmd", "/c", command)
	cmd.Env = []string{
		"LANG=en_US.UTF-8",
	}

	var enc mahonia.Decoder
	enc = mahonia.NewDecoder("gbk")

	if runtime.GOOS != "windows" {
		enc = mahonia.NewDecoder("utf-8")
		cmd = exec.Command("sh", "-c", command)
	}

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()
		reader := bufio.NewReader(stdout)
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}
			fmt.Println(enc.ConvertString(strings.ReplaceAll(readString, "\n", "")))
		}
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		reader := bufio.NewReader(stderr)
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}
			fmt.Println(enc.ConvertString(strings.ReplaceAll(readString, "\n", "")))
		}
	}()

	err = cmd.Start()
	if err != nil {
		logrus.Errorf("failed to start command: %s", command)
		return
	}
	logrus.Infof("PID: %d Command: %s", cmd.Process.Pid, command)
	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		logrus.Errorf("failed to complete command: %s", command)
		return
	}
	return
}
