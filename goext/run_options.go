package goext

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

//////////////////////////////
// Public Interface
//////////////////////////////

// Option that allows changing the working directory during a run.
func RunOptionInDirectory(path string) RunOption {
	return &runInDirectoryOption{path: path}
}

// Option that allows setting/overriding environment variables during a run.
func RunOptionWithEnvs(envVariables map[string]string) RunOption {
	return &runWithEnvsOption{envs: envVariables}
}

func RunInDirectory(path string, f func() error) (err error) {
	return RunWithOptions(f, RunOptionInDirectory(path))
}
func RunInDirectory1P[P1 any](path string, f func() (P1, error)) (P1, error) {
	return RunWithOptions1P(f, RunOptionInDirectory(path))
}
func RunInDirectory2P[P1 any, P2 any](path string, f func() (P1, P2, error)) (P1, P2, error) {
	return RunWithOptions2P(f, RunOptionInDirectory(path))
}
func RunInDirectory3P[P1 any, P2 any, P3 any](path string, f func() (P1, P2, P3, error)) (P1, P2, P3, error) {
	return RunWithOptions3P(f, RunOptionInDirectory(path))
}

func RunWithEnvs(envVariables map[string]string, f func() error) error {
	return RunWithOptions(f, RunOptionWithEnvs(envVariables))
}
func RunWithEnvs1P[P1 any](envVariables map[string]string, f func() (P1, error)) (P1, error) {
	return RunWithOptions1P(f, RunOptionWithEnvs(envVariables))
}
func RunWithEnvs2P[P1 any, P2 any](envVariables map[string]string, f func() (P1, P2, error)) (P1, P2, error) {
	return RunWithOptions2P(f, RunOptionWithEnvs(envVariables))
}
func RunWithEnvs3P[P1 any, P2 any, P3 any](envVariables map[string]string, f func() (P1, P2, P3, error)) (P1, P2, P3, error) {
	return RunWithOptions3P(f, RunOptionWithEnvs(envVariables))
}

//////////////////////////////
// Run In Directory
//////////////////////////////

// Option that allows changing the working directory during a run.
type runInDirectoryOption struct {
	path     string
	origPath string
}

func (r *runInDirectoryOption) apply() error {
	// Get the current directory
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current directory: %v", err)
	}
	r.origPath = pwd
	// Change the path
	err = os.Chdir(r.path)
	if err != nil {
		return fmt.Errorf("cannot change to directory %s: %v", strconv.Quote(r.path), err)
	}
	return nil
}

func (r *runInDirectoryOption) revert() error {
	// Reset to the previous folder
	err := os.Chdir(r.origPath)
	if err != nil {
		return fmt.Errorf("cannot change back to directory %s: %v", strconv.Quote(r.origPath), err)
	}
	return nil
}

//////////////////////////////
// Run With Envs
//////////////////////////////

// Option that allows setting/overriding environment variables during a run.
type runWithEnvsOption struct {
	envs     map[string]string
	origEnvs map[string]string
}

func (r *runWithEnvsOption) apply() (err error) {
	r.origEnvs = map[string]string{}
	for name, value := range r.envs {
		// Read and store original value
		if originalValue, ok := os.LookupEnv(name); ok {
			r.origEnvs[name] = originalValue
		}
		// Set the new value
		err = errors.Join(err, os.Setenv(name, value))
	}
	return
}

func (r *runWithEnvsOption) revert() (err error) {
	for name := range r.envs {
		origValue, ok := r.origEnvs[name]
		if ok {
			err = errors.Join(err, os.Setenv(name, origValue))
		} else {
			err = errors.Join(err, os.Unsetenv(name))
		}
	}
	return
}
