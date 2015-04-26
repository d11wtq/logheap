package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/samalba/dockerclient"
	"os"
	"time"
)

type Job struct {
	Id     string
	Stdout bool
	Stderr bool
}

func endpoint() string {
	if host := os.Getenv("DOCKER_HOST"); host != "" {
		return host
	} else {
		return "unix:///var/run/docker.sock"
	}
}

func filteredContainers(client dockerclient.Client) []dockerclient.Container {
	ret := make([]dockerclient.Container, 0)

	if list, err := client.ListContainers(true, false, ""); err == nil {
		for _, c := range list {
			if info, err := client.InspectContainer(c.Id); err == nil {
				if !info.Config.Tty {
					ret = append(ret, c)
				}
			}
		}
	}

	return ret
}

func encode(j Job, s string) (string, error) {
	var stream string

	switch {
	case j.Stderr:
		stream = "stderr"
	case j.Stdout:
		stream = "stdout"
	}

	document := map[string]interface{}{
		"message": s,
		"host":    os.Getenv("HOSTNAME"),
		"stream":  stream,
	}

	bytes, err := json.Marshal(document)
	return string(bytes), err
}

func processLogs(client dockerclient.Client, j Job, done chan Job) {
	opts := dockerclient.LogOptions{
		Stdout: j.Stdout,
		Stderr: j.Stderr,
		Follow: true,
		Tail:   0,
	}
	init := true

	for {
		if info, err := client.InspectContainer(j.Id); err == nil {
			if init || info.State.Running {
				if s, err := client.ContainerLogs(j.Id, &opts); err == nil {
					scanner := bufio.NewScanner(Demuxer(s))
					for scanner.Scan() {
						if doc, err := encode(j, scanner.Text()); err == nil {
							fmt.Println(doc)
						}
					}
					init = false
					opts.Tail = 1 // FIXME: Not accurate!
				}
			}
		} else {
			break
		}

		time.Sleep(time.Second * 3)
	}

	done <- j
}

func queueJobs(client dockerclient.Client, jobs map[Job]bool, done chan Job) {
	for {
		for _, c := range filteredContainers(client) {
			todos := []Job{
				{Id: c.Id, Stdout: true},
				{Id: c.Id, Stderr: true},
			}

			for _, item := range todos {
				if !jobs[item] {
					jobs[item] = true
					go processLogs(client, item, done)
				}
			}
		}
		time.Sleep(time.Second * 3)
	}
}

func main() {
	jobs := make(map[Job]bool)
	done := make(chan Job)

	client, _ := dockerclient.NewDockerClient(endpoint(), nil)

	go queueJobs(client, jobs, done)

	for k := range done {
		delete(jobs, k)
	}
}
