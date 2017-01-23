package tools

import "os"
import "path/filepath"

const jsonConfigName = "scv.json"

func FindConfig() (string, error) {
    currentDir, err := os.Getwd()

    if err != nil {
        return "", err
    }

    filename := filepath.Join(currentDir, jsonConfigName)

    _, err = os.Stat(filename)

    return filename, err
}
