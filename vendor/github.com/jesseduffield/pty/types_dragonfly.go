// +build ignore

package pty

/*
#define _KERNEL
#include <sys/conf.h>
#include <sys/param.h>
#include <sys/filio.h>
*/
import "C"

const (
	_C_SPECNAMELEN = C.SPECNAMELEN /* max length of devicename */
)

type fiodgnameArg C.struct_fiodname_args
