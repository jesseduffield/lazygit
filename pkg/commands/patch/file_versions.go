package patch

import "strings"

// FileVersions names the two versions of a file that a diff of it is between, as the
// diff's header states them: the id of the blob each version is stored as, abbreviated
// the way git wrote it, and the mode of each. git writes an id of zeros for a version
// that doesn't exist, the old one of an added file or the new one of a deleted one, and
// that version has no mode.
//
// With the same diff options, two diffs of a file between the same versions are the
// same diff, line for line. So the versions tell whether a diff read now is still the
// one read earlier, however the file, the index or HEAD changed in between.
type FileVersions struct {
	OldBlob string
	NewBlob string
	OldMode string
	NewMode string
}

// FileVersions returns the versions of the file the patch is a diff of.
func (self *Patch) FileVersions() FileVersions {
	return ParseFileVersions(self.header)
}

// ParseFileVersions reads the versions of a file from the header of its diff, given as
// the lines of the diff from its start; it stops at the first hunk. A header without
// an index line names no blobs. git leaves the index line out where the content
// doesn't change: a rename or copy of the file as it is, or a change of its mode alone.
func ParseFileVersions(lines []string) FileVersions {
	versions := FileVersions{}
	for _, line := range lines {
		if strings.HasPrefix(line, "@@") {
			break
		}
		switch {
		case strings.HasPrefix(line, "index "):
			// The mode is on the index line only if both versions have the same one.
			blobs, mode, _ := strings.Cut(strings.TrimPrefix(line, "index "), " ")
			versions.OldBlob, versions.NewBlob, _ = strings.Cut(blobs, "..")
			if mode != "" {
				versions.OldMode = mode
				versions.NewMode = mode
			}
		case strings.HasPrefix(line, "old mode "):
			versions.OldMode = strings.TrimPrefix(line, "old mode ")
		case strings.HasPrefix(line, "new mode "):
			versions.NewMode = strings.TrimPrefix(line, "new mode ")
		case strings.HasPrefix(line, "deleted file mode "):
			versions.OldMode = strings.TrimPrefix(line, "deleted file mode ")
		case strings.HasPrefix(line, "new file mode "):
			versions.NewMode = strings.TrimPrefix(line, "new file mode ")
		}
	}
	return versions
}

// Same reports whether two diffs are between the same versions of their file.
func (self FileVersions) Same(other FileVersions) bool {
	return sameBlob(self.OldBlob, other.OldBlob) &&
		sameBlob(self.NewBlob, other.NewBlob) &&
		self.OldMode == other.OldMode &&
		self.NewMode == other.NewMode
}

// sameBlob reports whether two abbreviated ids name the same blob. How many digits git
// abbreviates an id to depends on how many objects the repo has, so two diffs of the
// same blob can spell its id with different lengths; the shorter one is then a prefix
// of the longer one. An id that is missing matches only another missing one.
func sameBlob(a string, b string) bool {
	if a == "" || b == "" {
		return a == b
	}
	return strings.HasPrefix(a, b) || strings.HasPrefix(b, a)
}
