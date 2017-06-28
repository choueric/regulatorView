package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

type regulator struct {
	name      string
	index     int
	userNum   int
	uevent    string
	consumers []string
}

func getLine(fname string) (string, error) {
	s, err := ioutil.ReadFile(fname)
	if err != nil {
		return "", err
	}
	return string(s[0 : len(s)-1]), nil
}

func getInt(fname string) (int, error) {
	ret, err := ioutil.ReadFile(fname)
	if err != nil {
		return 0, err
	}
	n, err := strconv.Atoi(string(ret[0 : len(ret)-1]))
	if err != nil {
		return 0, err
	}

	return n, nil
}

func getString(fname string) (string, error) {
	s, err := ioutil.ReadFile(fname)
	if err != nil {
		return "", err
	}
	return string(s), nil
}

func getConsumers(r *regulator, files []os.FileInfo) {
	for _, file := range files {
		if file.Mode()&os.ModeSymlink != 0 {
			name := file.Name()
			if name != "device" && name != "subsystem" {
				r.consumers = append(r.consumers, name)
			}
		}
	}
}

func printRegulator(w io.Writer, r *regulator, verbose bool) {
	fmt.Fprintf(w, "[%3d]: %s, %d\n", r.index, r.name, r.userNum)
	if !verbose {
		return
	}
	if r.userNum > 0 {
		for _, c := range r.consumers {
			fmt.Fprintln(w, "      ", c)
		}
	}

	if r.uevent != "" {
		fmt.Fprintln(w, r.uevent)
	}
}
