package handlers

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"
)

type Process struct {
    cmd *exec.Cmd
    stdin   io.WriteCloser
    stdout  io.ReadCloser
    scanner *bufio.Scanner
}

var (
    processes = make(map[string]*Process)
    mu  sync.RWMutex
)

func NewProcess(name string) (*Process, error) {
    cmd := exec.Command(fmt.Sprintf("./handlers/%s", name))
    stdin, err := cmd.StdinPipe()
    if err != nil {
        return nil, fmt.Errorf("error creating stdin pipe: %v", err)
    }
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return nil, fmt.Errorf("error ceating stdout pipe: %v", err)
    }

    if err := cmd.Start(); err != nil {
        return nil, fmt.Errorf("error starting process: %v", err)
    }
    
    scanner := bufio.NewScanner(stdout)
    p := &Process{
        cmd:    cmd,
        stdin:  stdin,
        stdout: stdout,
        scanner: scanner,
    }

    mu.Lock()
    processes[name] = p
    mu.Unlock()

    log.Printf("Process %s started.\n", name)
    return p, nil
}

func (p *Process) Exec (command interface{}) (string, error) {
    _, err := fmt.Fprintln(p.stdin, command)
    if err != nil {
        return "", fmt.Errorf("error writing to stdin: %v", err)
    } 

    if p.scanner.Scan() {
        response := p.scanner.Text()
        return response, nil
    } else {
        if err := p.scanner.Err(); err != nil {
            return "", fmt.Errorf("error writing to stdin: %v", err)
        }
        return "", fmt.Errorf("no output received from process")
    }
}

func GetProcess(name string) (*Process, error) {
    mu.RLock()
    p, exists := processes[name]
    mu.RUnlock()
    if !exists {
        return nil, fmt.Errorf("process w ith name %s does not exist", name)
    }
    return p, nil
}

func DeleteProcess(name string) error {
    mu.Lock()
    p, exists := processes[name]
    if !exists {
        mu.Unlock()
        return fmt.Errorf("process with name %s does not exist", name)
    }
    delete(processes, name)
    mu.Unlock()

    return p.Close()
}

func (p *Process) Close() error {
    if err := p.stdin.Close(); err != nil {
        return fmt.Errorf("error closing stdin: %v", err)
    }

    if err := p.cmd.Wait(); err != nil {
        return fmt.Errorf("error wating for process to finish: %v", err)
    }

    return nil
}