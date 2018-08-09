// +build !norwfs

package dotgit

import (
	"os"

	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/utils/ioutil"
)

func (d *DotGit) setRef(fileName, content string, old *plumbing.Reference) (err error) {
	// If we are not checking an old ref, just truncate the file.
	mode := os.O_RDWR | os.O_CREATE
	if old == nil {
		mode |= os.O_TRUNC
	}

	f, err := d.fs.OpenFile(fileName, mode, 0666)
	if err != nil {
		return err
	}

	defer ioutil.CheckClose(f, &err)

	// Lock is unlocked by the deferred Close above. This is because Unlock
	// does not imply a fsync and thus there would be a race between
	// Unlock+Close and other concurrent writers. Adding Sync to go-billy
	// could work, but this is better (and avoids superfluous syncs).
	err = f.Lock()
	if err != nil {
		return err
	}

	// this is a no-op to call even when old is nil.
	err = d.checkReferenceAndTruncate(f, old)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte(content))
	return err
}
