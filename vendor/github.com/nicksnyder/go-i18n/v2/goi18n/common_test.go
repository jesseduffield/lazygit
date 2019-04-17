package main

import "io/ioutil"

func mustTempDir(prefix string) string {
	outdir, err := ioutil.TempDir("", prefix)
	if err != nil {
		panic(err)
	}
	return outdir
}
