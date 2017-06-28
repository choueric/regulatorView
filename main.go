package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

const REGULATOR_DIR = "/sys/class/regulator"

var regulators []*regulator

type ByIndex []*regulator

func (a ByIndex) Len() int           { return len(a) }
func (a ByIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByIndex) Less(i, j int) bool { return a[i].index < a[j].index }

func readRegulator(d string) (*regulator, error) {
	os.Chdir(d)
	defer os.Chdir(REGULATOR_DIR)

	s, _ := os.Getwd()
	fmt.Println(d, ":", s)

	idx, err := strconv.Atoi(strings.Split(d, ".")[1])
	if err != nil {
		return nil, err
	}
	r := &regulator{
		index: idx,
	}
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fname := file.Name()
		switch fname {
		case "name":
			r.name, err = getName(fname)
			if err != nil {
				return nil, err
			}
		case "num_users":
			r.userNum, err = getUserNum(fname)
			if err != nil {
				return nil, err
			}
			if r.userNum > 0 {
				getConsumers(r, files)
			}
		case "parent":
			fmt.Println("--------")
		}
	}

	return r, nil
}

func parsetRegulators() bool {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println(err)
		return false
	}

	for _, file := range files {
		if !file.IsDir() {
			r, err := readRegulator(file.Name())
			if err != nil {
				fmt.Println(file.Name(), err)
				continue
			}
			regulators = append(regulators, r)
		}
	}

	return true
}

func main() {
	fmt.Println("== regulator tree ==")
	os.Chdir(REGULATOR_DIR)

	parsetRegulators()

	sort.Sort(ByIndex(regulators))

	for _, r := range regulators {
		fmt.Printf("[%3d]: %s, %d\n", r.index, r.name, r.userNum)
		if r.userNum > 0 {
			for _, c := range r.consumers {
				fmt.Println("      ", c)
			}
		}
	}
}
