package main

import (
	"fmt"
	"github.com/chobostar/go-oom-guard/pkg/eventfd"
	pgkiller "github.com/chobostar/go-oom-guard/pkg/pg-backend-killer"
	"github.com/containerd/cgroups"
	"log"
)

const (
	defaultThresholdBytes = 100 * 1024 * 1024
	defaultKillSignal     = "TERM" //INT
)

func main() {

	printStat()

	threshold := uint64(defaultThresholdBytes)
	killSignal := defaultKillSignal

	efd := registerThreshold(threshold)
	log.Printf("registered event: %v, with threshold: %v, with killSignal: %v", efd, threshold, killSignal)

	log.Printf("starting to read events...")
	for {
		val, err := efd.ReadEvents()
		if err != nil {
			log.Printf("error while reading from eventfd: %v", err)
			break
		}
		log.Printf("got threshold event: %v", val)

		pgk := pgkiller.ParseBackends(getPids())
		log.Printf("%+v", pgk.GetStats())

		if pgk.IsExceedThreshold(threshold) {
			log.Printf("trying to kill top pid...")
			switch killSignal {
			case "TERM":
				err = pgk.KillTermTopPid()
			case "INT":
				err = pgk.KillIntTopPid()
			default:
				log.Fatalf("unknown killSignal %v", killSignal)
			}
			if err != nil {
				log.Fatalf("failed to kill, err: %v", err)
			}
			log.Printf("pid successfully killed")
		} else {
			log.Printf("total RSS usage not exceeding threshold, do nothing")
		}
	}
}

// https://github.com/containerd/cgroups#registering-for-memory-events
func registerThreshold(threshold uint64) *eventfd.EventFD {
	control := getCurrentCGroup()
	event := cgroups.MemoryThresholdEvent(threshold, false)
	efd, err := control.RegisterMemoryEvent(event)
	if err != nil {
		log.Fatalf("error while registering memory event %v", err)
	}
	return eventfd.FromFd(efd)
}

func printStat() {
	control := getCurrentCGroup()
	procs, err := control.Processes("memory", true)
	if err != nil {
		log.Fatalf("unable to load subsystem %v", err)
	}
	for _, proc := range procs {
		fmt.Printf("%+v\n", proc)
	}
}

func getPids() []int {
	control := getCurrentCGroup()
	procs, err := control.Processes("memory", true)
	if err != nil {
		log.Fatalf("unable to load subsystem %v", err)
	}
	pids := make([]int, len(procs))
	for i, proc := range procs {
		pids[i] = proc.Pid
	}
	return pids
}

func getCurrentCGroup() cgroups.Cgroup {
	control, err := cgroups.Load(cgroups.V1, cgroups.NestedPath(""))
	if err != nil {
		log.Fatalf("error while loading cgroup %v", err)
	}
	return control
}
