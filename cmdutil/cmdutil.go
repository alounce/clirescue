package cmdutil

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

// Ask a question and retreives user's answer, or error if fails. If question is a sensetive one, then hide user's answer
func Ask(question string, isSensetive bool) (string, error) {
    if isSensetive {
        silence()
        defer unsilence() 
    }
    fmt.Printf("%s : ", question)
    bytes, _, err := bufio.NewReader(os.Stdin).ReadLine()
    if err != nil {
        return "", err
    }
    line := string(bytes)
    return strings.TrimSpace(line), nil
}

// silence - hides user input in the console, useful when user is entering his password
func silence() {
    runCommand(exec.Command("stty", "-echo"))
}

// unsilence - Restores user input in the console, useful when user is entering his password
func unsilence() {
    runCommand(exec.Command("stty", "echo"))
}

func runCommand(command *exec.Cmd) {
    command.Stdin = os.Stdin
    command.Stdout = os.Stdout
    command.Run()
}
