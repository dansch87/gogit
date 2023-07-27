package main

import (
	"fmt"
	"os"
	"path/filepath"
	"errors"
	"log"
)


// TODOs

/* 1) init Function to create directory structure

	check if file already exists

 - .gogit directory
	- HEAD
	- objects/
		- info/
		- pack/
	- refs/
		- heads/
		- tags/

*/

func initCmd() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	gogitDirPath := filepath.Join(path, ".gogit")
	if _, err := os.Stat(gogitDirPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(gogitDirPath, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	fmt.Printf("The subdirectory %s is created\n", gogitDirPath)
}





func main() {

	cmd := os.Args[1:]
	if len(cmd) < 1 {
		log.Println("You must pass a valid sub-command")
		os.Exit(1)
	}
	if cmd[0] == "init" {
		initCmd()
	} else {
		fmt.Printf("Unknown subcommand %s\n", cmd)
		os.Exit(1) 
	}

}