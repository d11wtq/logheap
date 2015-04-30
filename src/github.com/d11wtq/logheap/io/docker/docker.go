package docker

import (
	"bufio"
	"github.com/d11wtq/logheap/io"
	"github.com/samalba/dockerclient"
	"net/url"
	"time"
)

// Job class, for a given container ID
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

func (i *Input) tags(job Job, info *dockerclient.ContainerInfo) map[string]interface{} {
	var stream string
	switch {
	case job.Stdout:
		stream = "stdout"
	case job.Stderr:
		stream = "stderr"
	}

	return map[string]interface{}{
		"stream":         stream,
		"container_id":   info.Id,
		"container_name": info.Name,
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
			tags := i.tags(job, info)

			if init || info.State.Running {
				if s, err := i.client.ContainerLogs(job.Id, &opts); err == nil {
					scanner := bufio.NewScanner(Demuxer(s))
					for scanner.Scan() {
						if doc, err := io.Encode(scanner.Text(), tags); err == nil {
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

// Register this input handler for use.
func Register() {
	io.RegisterInput(
		"docker",
		func(u *url.URL) (io.Input, error) {
			return &Input{}, nil
		},
	)
}
