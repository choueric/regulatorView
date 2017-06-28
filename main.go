package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const REGULATOR_DIR = "/sys/class/regulator"

var regulators []*regulator

type ByIndex []*regulator

func (a ByIndex) Len() int           { return len(a) }
func (a ByIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByIndex) Less(i, j int) bool { return a[i].index < a[j].index }

func readRegulator(d string, r *regulator) error {
	files, err := ioutil.ReadDir(d)
	if err != nil {
		return err
	}

	for _, file := range files {
		fname := file.Name()
		fpath := filepath.Join(d, fname)
		switch fname {
		case "name":
			r.name, err = getLine(fpath)
			if err != nil {
				return err
			}
		case "num_users":
			r.userNum, err = getInt(fpath)
			if err != nil {
				return err
			}
			if r.userNum > 0 {
				getConsumers(r, files)
			}
		case "uevent":
			r.uevent, err = getString(fpath)
			if err != nil {
				return err
			}
		case "parent":
			fmt.Println("--------")
		}
	}

	return nil
}

func parsetRegulators(d string) bool {
	files, err := ioutil.ReadDir(d)
	if err != nil {
		fmt.Println(err)
		return false
	}

	var wait sync.WaitGroup
	for _, file := range files {
		fname := file.Name()
		idx, err := strconv.Atoi(strings.Split(fname, ".")[1])
		if err != nil {
			fmt.Println(fname, err)
			continue
		}
		r := &regulator{index: idx}
		regulators = append(regulators, r)
		sub := filepath.Join(d, fname)
		wait.Add(1)
		go func(d string, r *regulator) {
			defer wait.Done()
			err := readRegulator(d, r)
			if err != nil {
				fmt.Println(d, err)
			}
		}(sub, r)
	}

	wait.Wait()

	return true
}

func main() {
	if err := initUi(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !parsetRegulators(REGULATOR_DIR) {
		fmt.Println("read regulators failed")
		os.Exit(2)
	}
	sort.Sort(ByIndex(regulators))

	for {
		input, err := ui.Ask("\n>>>")
		if err != nil {
			fmt.Println(err)
		}

		if handleInput(input) {
			break
		}
	}
}
