package docker

import (
	"bufio"
	"encoding/json"
	"github.com/d11wtq/logheap/io"
	"github.com/samalba/dockerclient"
	"os"
	"time"
)

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

type Job struct {
	Id     string
	Stdout bool
	Stderr bool
}

// Input handler for reading from docker.
type Input struct {
	client dockerclient.Client
	jobs   map[Job]bool
	done   chan Job
}

// Listen for incoming documents and process them.
func (i *Input) Listen(o io.Output) {
	i.client, _ = dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
	i.jobs = make(map[Job]bool)
	i.done = make(chan Job)

	go i.run(o)

	for job := range i.done {
		delete(i.jobs, job)
	}
}

func (i *Input) run(o io.Output) {
	for {
		for _, c := range i.filteredContainers() {
			todos := []Job{
				{Id: c.Id, Stdout: true},
				{Id: c.Id, Stderr: true},
			}

			for _, job := range todos {
				if !i.jobs[job] {
					i.jobs[job] = true
					go i.processLogs(job, o)
				}
			}
		}

		time.Sleep(time.Second * 3)
	}
}

func (i *Input) processLogs(job Job, o io.Output) {
	opts := dockerclient.LogOptions{
		Stdout: job.Stdout,
		Stderr: job.Stderr,
		Follow: true,
		Tail:   0,
	}
	init := true

	for {
		if info, err := i.client.InspectContainer(job.Id); err == nil {
			if init || info.State.Running {
				if s, err := i.client.ContainerLogs(job.Id, &opts); err == nil {
					scanner := bufio.NewScanner(Demuxer(s))
					for scanner.Scan() {
						if doc, err := encode(job, scanner.Text()); err == nil {
							o.Push(doc)
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

	i.done <- job
}

func (i *Input) filteredContainers() []dockerclient.Container {
	ret := make([]dockerclient.Container, 0)

	if list, err := i.client.ListContainers(true, false, ""); err == nil {
		for _, c := range list {
			if info, err := i.client.InspectContainer(c.Id); err == nil {
				if !info.Config.Tty {
					ret = append(ret, c)
				}
			}
		}
	}

	return ret
}
