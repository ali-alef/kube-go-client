package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Filter[E any](s []E, f func(E) bool) []E {
	res := make([]E, 0, len(s))
	for _, e := range s {
		if f(e) {
			res = append(res, e)
		}
	}
	return res
}

func ExecCommand(command string, args ...string) string {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + "-> " + stderr.String())
		return ""
	}

	return out.String()
}
