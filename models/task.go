package models

type Task struct {
    Arguments []string
    Commands  []Command
}

// api

func (self Task) GetArgumentsMap(cmdLineArgs []string) map[string]string {
    values := cmdLineArgs[2:]

    args := make(map[string]string)

    for i, key := range self.Arguments {
        args[key] = values[i]
    }

    return args
}
