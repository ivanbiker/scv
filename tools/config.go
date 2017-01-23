package tools

import "errors"

import "io/ioutil"
import "encoding/json"

import "github.com/ivanbiker/scv/models"

type Config struct {
    Hosts map[string]models.Host
    Tasks map[string]models.Task
}

// api

func NewConfig() *Config {
    return &Config{}
}

func (self *Config) Parse(filename string) error {
    bytes, err := ioutil.ReadFile(filename)

    if err != nil {
        return err
    }

    if err = json.Unmarshal(bytes, self); err != nil {
        return err
    }

    return nil
}

func (self Config) FindTask(cmdLineArgs []string) (*models.Task, error) {
    args := cmdLineArgs[1:]

    if len(args) == 0 {
        return nil, errors.New("list of available tasks: case 1")
    }

    key := args[0]

    task, ok := self.Tasks[key]

    if !ok {
        return nil, errors.New("list of available tasks: case 2")
    }

    if len(task.Arguments) != len(args[1:]) {
        return nil, errors.New("wrong task arguments")
    }

    return &task, nil
}

func (self Config) FindHost(key string) (*models.Host, error) {
    host, ok := self.Hosts[key]

    if !ok {
        return nil, errors.New("unknown host")
    }

    return &host, nil
}
