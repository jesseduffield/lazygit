// +build solaris

package ps

import (
	"encoding/binary"
	"fmt"
	"os"
)

type ushort_t uint16

type id_t int32
type pid_t int32
type uid_t int32
type gid_t int32

type dev_t uint64
type size_t uint64
type uintptr_t uint64

type timestruc_t [16]byte

// This is copy from /usr/include/sys/procfs.h
type psinfo_t struct {
	Pr_flag   int32     /* process flags (DEPRECATED; do not use) */
	Pr_nlwp   int32     /* number of active lwps in the process */
	Pr_pid    pid_t     /* unique process id */
	Pr_ppid   pid_t     /* process id of parent */
	Pr_pgid   pid_t     /* pid of process group leader */
	Pr_sid    pid_t     /* session id */
	Pr_uid    uid_t     /* real user id */
	Pr_euid   uid_t     /* effective user id */
	Pr_gid    gid_t     /* real group id */
	Pr_egid   gid_t     /* effective group id */
	Pr_addr   uintptr_t /* address of process */
	Pr_size   size_t    /* size of process image in Kbytes */
	Pr_rssize size_t    /* resident set size in Kbytes */
	Pr_pad1   size_t
	Pr_ttydev dev_t /* controlling tty device (or PRNODEV) */

	// Guess this following 2 ushort_t values require a padding to properly
	// align to the 64bit mark.
	Pr_pctcpu   ushort_t /* % of recent cpu time used by all lwps */
	Pr_pctmem   ushort_t /* % of system memory used by process */
	Pr_pad64bit [4]byte

	Pr_start    timestruc_t /* process start time, from the epoch */
	Pr_time     timestruc_t /* usr+sys cpu time for this process */
	Pr_ctime    timestruc_t /* usr+sys cpu time for reaped children */
	Pr_fname    [16]byte    /* name of execed file */
	Pr_psargs   [80]byte    /* initial characters of arg list */
	Pr_wstat    int32       /* if zombie, the wait() status */
	Pr_argc     int32       /* initial argument count */
	Pr_argv     uintptr_t   /* address of initial argument vector */
	Pr_envp     uintptr_t   /* address of initial environment vector */
	Pr_dmodel   [1]byte     /* data model of the process */
	Pr_pad2     [3]byte
	Pr_taskid   id_t      /* task id */
	Pr_projid   id_t      /* project id */
	Pr_nzomb    int32     /* number of zombie lwps in the process */
	Pr_poolid   id_t      /* pool id */
	Pr_zoneid   id_t      /* zone id */
	Pr_contract id_t      /* process contract */
	Pr_filler   int32     /* reserved for future use */
	Pr_lwp      [128]byte /* information for representative lwp */
}

func (p *UnixProcess) Refresh() error {
	var psinfo psinfo_t

	path := fmt.Sprintf("/proc/%d/psinfo", p.pid)
	fh, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fh.Close()

	err = binary.Read(fh, binary.LittleEndian, &psinfo)
	if err != nil {
		return err
	}

	p.ppid = int(psinfo.Pr_ppid)
	p.binary = toString(psinfo.Pr_fname[:], 16)
	return nil
}

func toString(array []byte, len int) string {
	for i := 0; i < len; i++ {
		if array[i] == 0 {
			return string(array[:i])
		}
	}
	return string(array[:])
}
