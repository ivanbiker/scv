package main

import "os"
import "fmt"
import "log"
import "strings"

import "github.com/ivanbiker/scv/types"
import "github.com/ivanbiker/scv/tools"
import "github.com/ivanbiker/scv/models"

func main() {
    fmt.Println("SCV good to go, sir!")
    fmt.Println("~~~~~~~~~~~~~~~~~~~~")

    configFilename, err := tools.FindConfig()

    if err != nil {
        log.Fatal(err)
    }

    config := tools.NewConfig()
    err = config.Parse(configFilename)

    if err != nil {
        log.Fatal(err)
    }

    task, err := config.FindTask(os.Args)

    if err != nil {
        log.Fatal(err)
    }

    localRunner := tools.NewLocalRunner()
    args := task.GetArgumentsMap(os.Args)

    for _, command := range task.Commands {
        if command.HostKey == "localhost" {
            err := executeCommand(command, args, localRunner)

            if err != nil {
                log.Fatal(err)
            }
        } else {
            host, err := config.FindHost(command.HostKey)

            if err != nil {
                log.Fatal(err)
            }

            fmt.Println("-> Connecting to " + host.GetUserAtHost() + "...")

            client := tools.NewSSHClient()
            err = client.Connect(*host)

            if err != nil {
                log.Fatal(err)
            }

            fmt.Println("<- Success!")

            err = executeCommand(command, args, client)

            if err != nil {
                log.Fatal(err)
            }

            fmt.Println("-> Disconnecting from " + host.GetUserAtHost() + "...")

            client.Disconnect()

            fmt.Println("<- Success!")
        } 
    }

    fmt.Println("~~~~~~~~~~~~~~~")
    fmt.Println("Job's finished.")
}

func executeCommand(command models.Command, args map[string]string, runner types.Runner) error {
    fmt.Println("-> Executing @ " + command.HostKey + ":")
    fmt.Println(formatLinesForPrint(command.Lines, args))

    line := packLines(command.Lines, args)
    output, err := runner.Run(line)

    if err != nil {
        return err
    }

    if strings.Trim(output, " \n\r") != "" {
        fmt.Println(output)
    }

    return nil
}

func packLines(lines []string, args map[string]string) string {
    line := strings.Join(lines, "&&")

    return applyArgs(line, args)
}

func formatLinesForPrint(lines []string, args map[string]string) string {
    newLines := make([]string, len(lines))

    for i, line := range lines {
        newLines[i] = "$ " + line
    }

    output := strings.Join(newLines, "\n")

    return applyArgs(output, args)
}

func applyArgs(s string, args map[string]string) string {
    for k, v := range args {
        s = strings.Replace(s, "$" + k, v, -1)
    }

    return s
}
