package pg_backend_killer

import (
	"github.com/prometheus/procfs"
	"strings"
	"syscall"
)

const (
	postgresComm = "postgres"
)

type PidStat struct {
	Pid      int
	CmdLine  string
	Comm     string
	RssUsage int
}

type PgBackendKiller struct {
	pids   []PidStat
	topPid int
	rssSum uint64
}

func ParseBackends(pids []int) *PgBackendKiller {
	fs, _ := procfs.NewFS("/proc")
	pidStats := make([]PidStat, len(pids))
	maxRss := -1
	topPid := -1
	rssSum := uint64(0)

	for i, pid := range pids {
		proc, err := fs.Proc(pid)
		if err != nil {
			continue
		}
		if isPgBackend(proc) {
			ps := getPidStat(proc)
			pidStats[i] = ps

			if maxRss < ps.RssUsage {
				maxRss = ps.RssUsage
				topPid = ps.Pid
			}

			rssSum += uint64(ps.RssUsage)
		}
	}

	return &PgBackendKiller{
		pids:   pidStats,
		topPid: topPid,
		rssSum: rssSum,
	}
}

func (pk *PgBackendKiller) IsExceedThreshold(threshold uint64) bool {
	return pk.rssSum > threshold
}

func (pk *PgBackendKiller) KillIntTopPid() error {
	return syscall.Kill(pk.topPid, syscall.SIGINT)
}

func (pk *PgBackendKiller) KillTermTopPid() error {
	return syscall.Kill(pk.topPid, syscall.SIGTERM)
}

func (pk *PgBackendKiller) GetStats() []PidStat {
	ret := make([]PidStat, len(pk.pids))
	for i, p := range pk.pids {
		retP := PidStat{
			Pid:      p.Pid,
			RssUsage: p.RssUsage,
			Comm:     p.Comm,
			CmdLine:  p.CmdLine,
		}
		ret[i] = retP
	}
	return ret
}

func isPgBackend(p procfs.Proc) bool {
	if comm, _ := p.Comm(); comm != postgresComm {
		return false
	}
	return true
}

func getPidStat(p procfs.Proc) PidStat {
	stat, _ := p.Stat()
	cmdLines, _ := p.CmdLine()
	comm, _ := p.Comm()

	return PidStat{
		Pid:      p.PID,
		Comm:     comm,
		CmdLine:  strings.Join(cmdLines, " "),
		RssUsage: stat.ResidentMemory(),
	}
}
