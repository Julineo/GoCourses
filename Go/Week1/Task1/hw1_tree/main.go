package main

import (
	"fmt"
	"io"
	"io/ioutil"
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

	flag := false

	var tree func(root, indent string) error

	tree = func(root, indent string) error {
		fi, err := os.Stat(root)
		if err != nil {
			return fmt.Errorf("could not stat %s: %v", root, err)
		}
		if fi.IsDir() {
			if flag {
				fmt.Println(fi.Name())
			}
		} else {
			fmt.Println(fi.Name())
			return nil
		}
		flag = true

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
				if !fi.IsDir() && !printFiles {
					continue
				}

				names = append(names, fi.Name())
			}
		}

		for i, name := range names {
			add := "│	"
			if i == len(names)-1 {
				fmt.Printf(indent + "└───")
				add = "	"
			} else {
				fmt.Printf(indent + "├───")
			}

			if err := tree(filepath.Join(root, name), indent+add); err != nil {
				return err
			}
		}

		return nil
	}

	if err := tree(path, ""); err != nil {
		return err
	}
	return nil
}
