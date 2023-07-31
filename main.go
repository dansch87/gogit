package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

func compressData(data []byte) []byte {
	var buf bytes.Buffer
	compr := zlib.NewWriter(&buf)
	compr.Write(data)
	compr.Close()
	output := buf.Bytes()
	return output
}

func writeDataToDb(data []byte, objectType string) error {

	sha1Hash := hashObject(data, objectType)

	objDir := filepath.Join(".gogit", "objects", sha1Hash[:2])
	fileName := sha1Hash[2:]
	objFullPath := filepath.Join(objDir, fileName)

	// Create object directory
	if _, err := os.Stat(objDir); os.IsNotExist(err) {
		err := os.Mkdir(objDir, os.ModeDir)
		if err != nil {
			return errors.New("object directory couldn't be created")
		}
	}

	f, err := os.Create(objFullPath)
	if err != nil {
		return err
	}
	defer f.Close()

	zlibComprData := compressData(data)
	cnt, err := f.Write(zlibComprData)
	if err != nil {
		return err
	}
	fmt.Printf("%v bytes written to object database\n", cnt)

	return nil
}

func main() {

	testData := []byte("what is up, doc?")

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
		sha1Hash := hashObject(testData, "blob")
		fmt.Println(sha1Hash)
	} else if cmd[0] == "write" {
		err := writeDataToDb(testData, "blob")
		if err != nil {
			fmt.Println("Error of test")
		}
	} else {
		fmt.Printf("unknown subcommand: %v\n", string(cmd[0]))
	}
}
