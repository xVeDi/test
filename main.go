package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

var stack []string

func push(st string) {
	stack = append(stack, st)
	return
}

func pop() {
	ln := len(stack)
	if ln > 0 {
		stack = stack[:ln-1]
	} else {
		stack = []string{}
	}
	return
}

func printPathGraph() (out string) {
	for _, i := range stack {
		out += i
	}
	return
}

func printPrefix(isLast bool) string {
	out := printPathGraph()
	if isLast {
		out += "└───"
	} else {
		out += "├───"
	}
	return out
}

func getSizeStr(size int64) (strSize string) {
	strSize = " ("
	if size > 0 {
		strSize += strconv.Itoa(int(size)) + "b"
	} else {
		strSize += "empty"
	}
	strSize += ")"
	return
}

func dirTree(out io.Writer, path string, printFiles bool) error {

	defer pop()

	var err error
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		//log.Fatal(err)
	}

	countFiles := 0

	//
	if printFiles {
		countFiles = len(files)
	} else {
		for _, f := range files {
			fi, err := os.Stat(path + "/" + f.Name())
			if err != nil {
				//log.Fatal(err)
			}

			if fi.Mode().IsDir() {
				countFiles++
			}
		}
	}

	if countFiles == 0 {
		return err
	}

	for _, f := range files {
		countFiles--
		fi, err := os.Stat(path + "/" + f.Name())
		if err != nil {
			log.Fatal(err)
		}
		isLast := (countFiles == 0)

		switch mode := fi.Mode(); {
		case mode.IsRegular():
			if printFiles {
				fmt.Fprint(out, printPrefix(isLast), f.Name(), getSizeStr(f.Size()), "\n")
			} else {
				countFiles++
			}
		case mode.IsDir():
			fmt.Fprint(out, printPrefix(isLast)+f.Name()+"\n")
			if isLast {
				push("	")
			} else {
				push("│	")
			}
			dirTree(out, path+"/"+f.Name(), printFiles)
		}
	}
	return err
}

func main() {

	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
