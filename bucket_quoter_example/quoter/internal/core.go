package internal

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"syscall"

	"golang.org/x/sys/unix"
)

type Core struct {
	g *CmdGlobal

	httpapi *Api
}

type CoreOverrides struct {
	Detach bool `json:"detach"`
}

func CreateCore(g *CmdGlobal) (*Core, error) {
	var err error

	var core Core
	core.g = g

	return &core, err
}

func (c *Core) Daemon() error {
	if os.Getppid() != 1 {
		// I am the parent, spawn child to run as daemon
		binary, err := exec.LookPath(os.Args[0])
		if err != nil {
			c.g.Log.Error(fmt.Sprintf("Failed to lookup binary, err:'%s'", err))
			return err
		}
		_, err = os.StartProcess(binary, os.Args, &os.ProcAttr{Dir: "", Env: nil,
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}, Sys: nil})
		if err != nil {
			c.g.Log.Error(fmt.Sprintf("Failed to start process, err:'%s'", err))
			return err
		}
		os.Exit(0)
	} else {
		// I am the child, i.e. the daemon, start new session and detach from terminal
		_, err := syscall.Setsid()
		if err != nil {
			c.g.Log.Error(fmt.Sprintf("Failed to create new session: err:'%s'", err))
			return err
		}
		file, err := os.OpenFile("/dev/null", os.O_RDWR, 0)
		if err != nil {
			c.g.Log.Error(fmt.Sprintf("Failed to open /dev/null, err:'%s'", err))
			return err
		}

		unix.Dup2(int(file.Fd()), int(os.Stdin.Fd()))
		unix.Dup2(int(file.Fd()), int(os.Stdout.Fd()))
		unix.Dup2(int(file.Fd()), int(os.Stderr.Fd()))
		file.Close()

		// writing pidfile
		pid := fmt.Sprintf("%d", os.Getpid())
		pidfile := c.g.Opts.RateLimiter.PidFile
		if err = ioutil.WriteFile(pidfile, []byte(pid), 0644); err != nil {
			c.g.Log.Error(fmt.Sprintf("error writing pidfile '%s', err:'%s'", pidfile, err))
			return err
		}
	}

	return nil
}

func (c *Core) IfProcessRun() (bool, int, *os.Process, error) {
	var pid int64

	pidfile := c.g.Opts.RateLimiter.PidFile
	var content []byte
	var err error
	if content, err = ioutil.ReadFile(pidfile); err != nil {
		c.g.Log.Error(fmt.Sprintf("error reading pidfile '%s', err:'%s'", pidfile, err))
		return false, int(pid), nil, err
	}

	if len(content) == 0 {
		err = errors.New("empty pidfile")
		c.g.Log.Error(fmt.Sprintf("error reading pidfile '%s', err:'%s'", pidfile, err))
		return false, int(pid), nil, err
	}

	spid := fmt.Sprintf("%s", content)
	if pid, err = strconv.ParseInt(spid, 10, 64); err != nil {
		c.g.Log.Error(fmt.Sprintf("error converting pidfile '%s' content into int, err:'%s'", pidfile, err))
		return false, int(pid), nil, err
	}

	c.g.Log.Debug(fmt.Sprintf("detecting pidfile:'%s' pid:'%d'", pidfile, pid))

	var p *os.Process
	if p, err = os.FindProcess(int(pid)); err != nil {
		c.g.Log.Error(fmt.Sprintf("error detecting process, pid:'%d', err:'%s'", pid, err))
		return false, int(pid), nil, err
	}

	if err = p.Signal(syscall.SIGCONT); err != nil {
		err = errors.New("process not running")
		c.g.Log.Debug(fmt.Sprintf("detected pid:'%d' not running", pid))
		return false, int(pid), p, err
	}

	return true, int(pid), p, nil
}

func (c *Core) StopCore(Overrides *CoreOverrides) error {
	if Overrides.Detach {
		var err error
		var run bool
		var pid int
		var p *os.Process

		if run, pid, p, err = c.IfProcessRun(); err != nil {
			c.g.Log.Error(fmt.Sprintf("%s error detecting run process, pid:'%d', err:'%s'", pid, err))
			return err
		}

		if !run {
			err = errors.New("no running process detected")
			c.g.Log.Error(fmt.Sprintf("%s no detecting run process, pid:'%d', err:'%s'", pid, err))
			return err
		}

		if err = p.Signal(syscall.SIGTERM); err != nil {
			c.g.Log.Error(fmt.Sprintf("error signalling pid:'%d', err:'%s'", pid, err))
			return err
		}
	}

	return nil
}

func (c *CoreOverrides) AsString() string {
	return fmt.Sprintf("detach:'%t'", c.Detach)
}

func (c *Core) StartCore(Overrides *CoreOverrides) error {
	if Overrides.Detach {
		var err error
		if os.Getppid() != 1 {
			// Detecting if process running already
			var run bool
			var pid int

			if run, pid, _, err = c.IfProcessRun(); err == nil && run {
				c.g.Log.Error(fmt.Sprintf("process already run on pid:'%d'", pid))
				return err
			}
		}

		if err = c.Daemon(); err != nil {
			c.g.Log.Error(fmt.Sprintf("error detaching process, err:'%s'", err))
			return err
		}
	}

	c.g.Log.Debug(fmt.Sprintf("starting core, overrides:'%s'", Overrides.AsString()))

	var waitGroup sync.WaitGroup
	count := 1
	waitGroup.Add(count)

	api, _ := CreateApi(c.g)
	api.core = c
	c.httpapi = api

	// API methods: unix socket and https
	// (for remote calls)
	//go api.Apiloop(&waitGroup, API_UNIXSOCKET)
	go api.Apiloop(&waitGroup, API_HTTPS)

	waitGroup.Wait()

	return nil
}
