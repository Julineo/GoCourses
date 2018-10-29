package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

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

func dirTree(out io.Writer, path string, printFiles bool) error {
	visit := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			fmt.Println("dir:  ", path)
		} /* else {
			fmt.Println("file: ", path)
		}*/
		return nil
	}

	err := filepath.Walk(path, visit)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

/*
func dirTree(out io.Writer, path string, printFiles bool) error {
	var txt string
	var level int

	var visit func(prev, dir string)
	visit = func(prev, dir string) {
		level++
		fmt.Println(strings.Repeat("\t", level-1) + dir)
		//txt = txt + dir + "\n"
		_ = txt

		files, err := ioutil.ReadDir(prev + "/" + dir)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			if file.IsDir() {
				visit(prev+"/"+dir, file.Name())
			}
		}
		level--
	}

	visit(".", path)
	// ├ ─ └ │
	//fmt.Println(path)
	fmt.Fprintln(out, txt)
	return nil
}
*/
