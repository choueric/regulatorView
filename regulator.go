package main

import (
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

func getString(fname string) (string, error) {
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
