package common

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type TestExecCmd struct {

	// binary to invoke
	// Example: "func", "kn", "kubectl", "/usr/bin/sh"
	Binary string

	// Binary args to append before actual args. Examples:
	//  when 'kn' binary binaryArgs should be ["func"]
	BinaryArgs []string

	// Run commands from Dir
	SourceDir string

	// Indicates shell should dump command line args during execution
	ShouldDumpCmdLine bool

	// Indicates shell should dump
	ShouldDumpOnSuccess bool

	// Fail Test on Error
	ShouldFailOnError bool

	// Environment variable to be used with the command
	Env []string

	// Optional function to be used to dump stdout command results
	DumpLogger func(out string)

	// Boolean
	T *testing.T
}

// TestExecCmdResult stored command result
type TestExecCmdResult struct {
	Out   string
	Error error
}

func (f *TestExecCmd) WithEnv(envKey string, envValue string) *TestExecCmd {
	env := envKey + "=" + envValue
	f.Env = append(f.Env, env)
	return f
}

func (f *TestExecCmd) FromDir(dir string) *TestExecCmd {
	f.SourceDir = dir
	return f
}

func (f *TestExecCmd) Run(oneArgs string) TestExecCmdResult {
	args := strings.Split(oneArgs, " ")
	return f.Exec(args...)
}

// Exec invokes go exec library and runs a shell command combining the binary args with args from method signature
func (f *TestExecCmd) Exec(args ...string) TestExecCmdResult {
	finalArgs := f.BinaryArgs
	if finalArgs == nil {
		finalArgs = args
	} else if args != nil {
		finalArgs = append(finalArgs, args...)
	}

	if f.ShouldDumpCmdLine {
		f.T.Log(f.Binary, strings.Join(finalArgs, " "))
	}

	var out bytes.Buffer

	cmd := exec.Command(f.Binary, finalArgs...)
	cmd.Stderr = &out
	cmd.Stdout = &out
	if f.SourceDir != "" {
		cmd.Dir = f.SourceDir
	}
	cmd.Env = append(os.Environ(), f.Env...)
	err := cmd.Run()

	result := TestExecCmdResult{
		Out:   out.String(),
		Error: err,
	}

	if err == nil && f.ShouldDumpOnSuccess {
		if result.Out != "" {
			if f.DumpLogger != nil {
				f.DumpLogger(result.Out)
			} else {
				f.T.Logf("%v", result.Out)
			}
		}
	}
	if err != nil {
		f.T.Log(result.Out)
		f.T.Log(err.Error())
		if f.ShouldFailOnError {
			f.T.Fail()
		}
	}

	return result
}
