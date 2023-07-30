package main

import (
	"fmt"
	"os"
	"path/filepath"
	"errors"
	"log"
	"crypto/sha1"
	"encoding/hex"
)


func createRepo() error {
	/*
	Repository dir/file structure:
	 - .gogit
	    - .gogit/objects/
		- .gogit/refs/
		   - gogit/refs/heads/
		   - gogit/refs/tags/
		- .gogit/HEAD
		- .gogit/config (not included)
		- .gogit/description (not included)
	*/

	projDir, err := os.Getwd()
	if err != nil {
		return err
	}

	repoDir := filepath.Join(projDir, ".gogit/")

	if _, err := os.Stat(repoDir); os.IsExist(err) {
		return errors.New("repository already exist")
	}

	// Create dir struct
	subDir := []string{"/objects/", "/refs/heads/", "/refs/tags/"}
	for _, subDir := range subDir {
		fp := filepath.Join(".gogit/", subDir)
		err := os.MkdirAll(fp, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// Create HEAD file
	file, err := os.Create(".gogit/HEAD")
    if err != nil {
        return err
    }
    defer file.Close()
    buffer := []byte("ref: refs/heads/main\n")
    file.Write(buffer)

	return nil
}


func hashObject(data []byte, objectType string) string {
	byteSize := len(data)
	nullByte := []byte("\x00")
	header := fmt.Sprintf("%s %d%s", objectType, byteSize, string(nullByte))

	fullData := append([]byte(header), data...)

	hasher := sha1.New()
	hasher.Write(fullData)
	sha1Hash := hex.EncodeToString(hasher.Sum(nil))

	return sha1Hash
}






func main() {
	cmd := os.Args[1:]
	if len(cmd) < 1 {
		log.Println("You must pass a valid sub-command")
	}
	if cmd[0] == "init" {
		err := createRepo()
		if err != nil {
			log.Fatal(err)
		}
	} else if cmd[0] == "hash" {
		testData := []byte("what is up, doc?")
		sha1Hash := hashObject(testData, "blob")
		fmt.Println(sha1Hash)
	} else {

		fmt.Printf("unknown subcommand: %v\n", string(cmd[0]))
	}
}