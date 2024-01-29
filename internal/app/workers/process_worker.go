package workers

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
)

type Worker struct {
	ID     string
	Cancel context.CancelFunc
	Ctx    context.Context

	ProcessStdout chan string
	ProcessStderr chan string
	ProcessStdIn  chan string
	ProcessSignal chan os.Signal

	startCommand string
	stopCommand  string
	directory    string
	cmd          *exec.Cmd
}

type MinecraftServerWorker struct {
	Worker
	MemoryMB int
}

func CreateMinecraftServerWorker(id string, memoryMB int) (*MinecraftServerWorker, error) {
	msw := &MinecraftServerWorker{}

	if memoryMB < 128 {
		return nil, errors.New("MemoryMB must be greater than 128")
	}

	msw.Ctx, msw.Cancel = context.WithCancel(context.Background())

	msw.ID = id
	msw.MemoryMB = memoryMB
	msw.ProcessStdout = make(chan string)
	msw.ProcessStderr = make(chan string)
	msw.ProcessStdIn = make(chan string)
	msw.ProcessSignal = make(chan os.Signal)

	minMemoryMB := 128
	maxMemoryMB := msw.MemoryMB

	if maxMemoryMB < minMemoryMB {
		maxMemoryMB = minMemoryMB
	}

	msw.startCommand = "java -Xms" + fmt.Sprint(minMemoryMB) + "M -Xmx" + fmt.Sprint(maxMemoryMB) + "M -jar server.jar nogui"
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	msw.directory = workingDir + "/minecraft_servers/" + msw.ID
	// Make sure the directory exists
	if _, err := os.Stat(msw.directory); os.IsNotExist(err) {
		err = os.MkdirAll(msw.directory, 0755)
		if err != nil {
			return nil, err
		}
	}

	return msw, nil
}

func (worker *Worker) Start() {
	isWindows := runtime.GOOS == "windows"
	var cmd *exec.Cmd
	if isWindows {
		cmd = exec.CommandContext(worker.Ctx, "cmd", "/C", worker.startCommand)
	} else {
		cmd = exec.CommandContext(worker.Ctx, "sh", "-c", worker.startCommand)
	}
	cmd.Dir = worker.directory

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		worker.ProcessStderr <- err.Error()
	}
	defer stderr.Close()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		worker.ProcessStderr <- err.Error()
	}
	defer stdin.Close()

	if err := cmd.Start(); err != nil {
		worker.ProcessStderr <- err.Error()
	}
	worker.cmd = cmd

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			worker.ProcessStdout <- scanner.Text()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			worker.ProcessStderr <- scanner.Text()
		}
	}()
	wg.Add(1)
	go func() {
		for {
			select {
			case <-worker.Ctx.Done():
                wg.Done()
				return
			case message := <-worker.ProcessStdIn:
				if len(message) > 0 {
					stdin.Write([]byte(message))
				}
			}
		}
	}()
	wg.Wait()
}
