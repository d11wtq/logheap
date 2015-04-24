package main

import (
	"github.com/samalba/dockerclient"
	"io"
	"os"
)

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
			info, _ := client.InspectContainer(c.Id)
			if !info.Config.Tty {
				ret = append(ret, c)
			}
		}
	}

	return ret
}

func processLogs(client dockerclient.Client, id string, done chan string) {
	opts := dockerclient.LogOptions{
		Stdout: true,
		Follow: true,
	}

	if s, err := client.ContainerLogs(id, &opts); err == nil {
		io.Copy(os.Stdout, Demuxer(s))
	}

	done <- id
}

func main() {
	jobs := make(map[string]bool, 10)
	done := make(chan string)

	client, _ := dockerclient.NewDockerClient(endpoint(), nil)

	for _, c := range filteredContainers(client) {
		jobs[c.Id] = true
		go processLogs(client, c.Id, done)
	}

	for id := range done {
		delete(jobs, id)
	}
}
