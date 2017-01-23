package tools

import "strings"

func PackCommandLines(commands []string, args map[string]string) string {
    command := strings.Join(commands, "&&")

    for k, v := range args {
        command = strings.Replace(command, "$" + k, v, -1)
    }

    return command
}
