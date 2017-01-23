package models

import "strconv"

type Host struct {
    Hostname string
    Port     int
    Username string
    Password string
    Keypath  string
}

// api

func (self Host) GetHostname() string {
    hostname := self.Hostname

    if self.Port != 0 {
        hostname += ":" + strconv.Itoa(self.Port)
    }

    return hostname
}

func (self Host) GetUserAtHost() string {
    return self.Username + "@" + self.Hostname
}
