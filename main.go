package main

import (
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"strings"
)

type Pod struct {
	name     string
	ready    string
	status   string
	restarts string
	age      string
	line     string
}

func GetPods(nameSpace string, nameFilter ...string) []Pod {
	var output string

	if nameSpace != "" {
		output = ExecCommand("kubectl", "get", "po", "-n", nameSpace)
	} else {
		output = ExecCommand("kubectl", "get", "po")
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")[1:]
	pods := make([]Pod, len(lines))

	for i, line := range lines {
		res := strings.Split(line, "  ")
		res = Filter(res, func(s string) bool {
			return s != "" && s != " "
		})

		name := strings.TrimSpace(res[0])
		ready := strings.TrimSpace(res[1])
		status := strings.TrimSpace(res[2])
		restarts := strings.TrimSpace(res[3])
		age := strings.TrimSpace(res[4])

		pods[i] = Pod{
			name:     name,
			ready:    ready,
			status:   status,
			restarts: restarts,
			age:      age,
			line:     line,
		}
	}

	pods = Filter(pods, func(pod Pod) bool {
		status := true
		for _, name := range nameFilter {
			if !strings.Contains(pod.name, name) {
				status = false
				break
			}
		}

		return status
	})

	return pods
}

func ExecPod(nameSpace string, nameFilters ...string) {
	pods := GetPods(nameSpace, nameFilters...)

	var pod Pod
	if len(pods) == 0 {
		fmt.Println("No Pod Found :(")
		return
	}

	if len(pods) > 1 {
		fmt.Println("Which one?")
		for i, p := range pods {
			fmt.Printf("\t%d) %s %s\n", i, p.name, p.age)
		}

		var inputId int

		fmt.Println("Enter id")
		_, _ = fmt.Scan()
		_, err := fmt.Scanf("%d", &inputId)
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		pod = pods[inputId]
	} else {
		pod = pods[0]
	}

	baseExecCommand := "kubectl exec -it %s %s -- bash"

	var command string
	if nameSpace != "" {
		command = fmt.Sprintf(baseExecCommand, "-n "+nameSpace, pod.name)
	} else {
		command = fmt.Sprintf(baseExecCommand, "", pod.name)
	}

	err := clipboard.WriteAll(command)
	if err != nil {
		fmt.Println("Error With copying in clipboard:", err)
		return
	}
	fmt.Println("copied to clipboard :D")
}

type NameFilter []string

func (n *NameFilter) String() string {
	return fmt.Sprintf("%v", *n)
}

func (n *NameFilter) Set(value string) error {
	parts := strings.Split(value, ",")
	for _, part := range parts {
		*n = append(*n, part)
	}
	return nil
}

func main() {
	nameSpace := flag.String("n", "native", "name space")
	action := flag.Int("act", 1, "available actions:\n\t1) get pods\n\t2) exec pod")
	var nameFilters NameFilter
	flag.Var(&nameFilters, "filter", "filter pod names using comma")
	flag.Parse()

	switch *action {
	case 1:
		pods := GetPods(*nameSpace, nameFilters...)
		for _, pod := range pods {
			fmt.Println(pod.line)
		}
	case 2:
		ExecPod(*nameSpace, nameFilters...)
	}
}
