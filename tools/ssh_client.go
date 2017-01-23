package tools

import "io/ioutil"
import "encoding/pem"

import "golang.org/x/crypto/ssh"

import "github.com/ivanbiker/scv/models"

type SSHClient struct {
    client  *ssh.Client
}

// api

func NewSSHClient() *SSHClient {
    return &SSHClient{}
}

func (self *SSHClient) Connect(host models.Host) error {
    config := createClientConfig(host)

    client, err := ssh.Dial("tcp", host.GetHostname(), config)

    if err != nil {
        return err
    }

    self.client = client

    return nil
}

func (self *SSHClient) Run(command string) (output string, err error) {
    session, err := self.client.NewSession()

    if err != nil {
        self.client.Close()
        return "", err
    }

    defer session.Close()

    outputBytes, err := session.CombinedOutput(command)
    output = string(outputBytes)

    if err != nil {
        return
    }

    return
}

func (self *SSHClient) Disconnect() {
    self.client.Close()
}

// helpers

func createClientConfig(host models.Host) *ssh.ClientConfig {
    authMethods := make([]ssh.AuthMethod, 0, 3)

    authMethods = addPublicKeyMethod(authMethods, host)
    authMethods = addPasswordMethod(authMethods, host)
    authMethods = addKeyboardInteractiveDummyMethod(authMethods, host)

    return &ssh.ClientConfig{
        User: host.Username,
        Auth: authMethods,
    }
}

func addPublicKeyMethod(methods []ssh.AuthMethod, host models.Host) []ssh.AuthMethod {
    if host.Keypath == "" {
        return methods
    }

    pemBytes, err := ioutil.ReadFile(host.Keypath)

    if err != nil {
        panic(err)
    }

    block, _ := pem.Decode(pemBytes)

    if block == nil {
        panic("no key found in keypath")
    }

    signer, err := ssh.ParsePrivateKey(pemBytes)

    if err != nil {
        panic(err)
    }

    return append(methods, ssh.PublicKeys(signer))
}

func addPasswordMethod(methods []ssh.AuthMethod, host models.Host) []ssh.AuthMethod {
    if host.Password == "" {
        return methods
    } else {
        return append(methods, ssh.Password(host.Password))
    }
}

func addKeyboardInteractiveDummyMethod(methods []ssh.AuthMethod, host models.Host) []ssh.AuthMethod {
    if host.Password == "" {
        return methods
    }

    keyboardInteractiveCallback := func(
        user,
        instruction string,
        questions []string,
        echos []bool,
    ) (answers []string, err error) {
        answers = make([]string, len(questions))

        for i, _ := range questions {
            answers[i] = host.Password
        }

        return
    }

    return append(methods, ssh.KeyboardInteractive(keyboardInteractiveCallback))
}

//func requestPassword(prompt string) (password string, err error) {
//    terminalState, err := terminal.GetState(0)
//
//    if err != nil {
//        return
//    }
//
//    incomeSignals := make(chan os.Signal, 1)
//    signal.Notify(
//        incomeSignals,
//        os.Interrupt,
//        syscall.SIGTERM,
//        syscall.SIGQUIT,
//    )
//
//    go func() {
//        <-incomeSignals
//        terminal.Restore(0, terminalState)
//        fmt.Println()
//        os.Exit(2)
//    }()
//
//    defer func() {
//        signal.Stop(incomeSignals)
//        close(incomeSignals)
//    }()
//
//    writer := bufio.NewWriter(os.Stdout)
//    writer.Write([]byte(prompt))
//    writer.Flush()
//
//    passwordBytes, err := terminal.ReadPassword(0)
//
//    if err != nil {
//        return
//    }
//
//    password = string(passwordBytes)
//
//    writer.Write([]byte("\n"))
//    writer.Flush()
//
//    return
//}
