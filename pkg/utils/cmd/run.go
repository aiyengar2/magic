package cmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

var (
	VerboseEnvVar = "VERBOSE"
	Verbose       = false
)

func init() {
	if os.Getenv(VerboseEnvVar) == "1" {
		Verbose = true
	}
}

func RunV(cmd string, args ...string) error {
	return RunWithV(nil, cmd, args...)
}

func RunWithV(env map[string]string, cmd string, args ...string) error {
	_, err := RunWith(env, cmd, args...)
	return err
}

func Run(cmd string, args ...string) (string, error) {
	return RunWith(nil, cmd, args...)
}

func RunWith(env map[string]string, cmd string, args ...string) (string, error) {
	buf := &bytes.Buffer{}
	var stdout, stderr io.Writer
	if Verbose {
		stdout = io.MultiWriter(buf, os.Stdout)
		stderr = os.Stderr
	} else {
		stdout = buf
		stderr = os.Stderr
	}
	err := Exec(env, os.Stdin, stdout, stderr, cmd, args...)
	if err != nil {
		_, _ = io.Copy(os.Stdout, buf)
	}
	return buf.String(), err
}

func Exec(env map[string]string, stdin io.Reader, stdout, stderr io.Writer, cmd string, args ...string) error {
	PrintWith(env, "internal:exec", cmd, args...)
	c := exec.Command(cmd, args...)
	c.Env = append(os.Environ(), MapToEnvSlice(env)...)
	c.Stderr = stderr
	c.Stdout = stdout
	return c.Run()
}
