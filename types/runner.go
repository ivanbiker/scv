package types

type Runner interface {
    Run(string) (string, error)
}
