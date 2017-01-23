package tools

import "os/exec"
import "strings"

type LocalRunner struct {}

// api

func NewLocalRunner() *LocalRunner {
    return &LocalRunner{}
}

func (self LocalRunner) Run(line string) (string, error) {
    command := self.createCommand(line)

    output, err := command.Output()

    if err != nil {
        return "", err
    }

    return string(output), nil
}

// helpers

func (self LocalRunner) createCommand(line string) *exec.Cmd {
    var name string
    args := make([]string, 0, 10)

    for _, part := range strings.Split(line, " ") {
        if part == "" {
            continue
        }

        if name == "" {
            name = part
            continue
        }

        args = append(args, part)
    }

    return exec.Command(name, args...)
}
