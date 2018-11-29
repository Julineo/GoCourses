package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	args := []string{"."}
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	offFile := true

	for _, f := range args {
		if f == "-f" {
			offFile = false
		}
	}

	for _, arg := range args {
		if arg == "-f" { continue }
		err := tree(arg, "", offFile)
		if err != nil {
			log.Printf("tree %s: %v\n", arg, err)
		}
	}
}

func tree(root, indent string, offFile bool) error {
	fi, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("could not stat %s: %v", root, err)
	}
	fmt.Println(fi.Name())
	if !fi.IsDir() {
		return nil
	}

	fis, err := ioutil.ReadDir(root)
	if err != nil {
		return fmt.Errorf("could not read dir %s: %v", root, err)
	}

	var names []string
	for _, fi := range fis {
		if fi.Name()[0] != '.' {

			// check if file -f continue
			fi, err := os.Stat(filepath.Join(root, fi.Name()))
				if err != nil {
					fmt.Errorf("could not stat %s: %v", filepath.Join(root, fi.Name()), err)
				}
			if !fi.IsDir() && offFile {
				continue
			}

			names = append(names, fi.Name())
		}
	}

	for i, name := range names {
		add := "│  "
		if i == len(names)-1 {
			fmt.Printf(indent + "└───")
			add = "   "
		} else {
			fmt.Printf(indent + "├───")
		}

		if err := tree(filepath.Join(root, name), indent+add, offFile); err != nil {
			return err
		}
	}

	return nil
}
