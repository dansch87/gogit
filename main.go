package main

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
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

func compressData(data []byte) []byte {
	var buf bytes.Buffer
	compr := zlib.NewWriter(&buf)
	compr.Write(data)
	compr.Close()
	output := buf.Bytes()
	return output
}

func decompressData(data []byte) ([]byte, error) {
	buf := bytes.NewReader(data[:])
	z, err := zlib.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer z.Close()

	p, err := ioutil.ReadAll(z)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func writeObject(data []byte, objectType string, writeMode bool) (string, error) {

	// Extend data by header data
	byteSize := len(data)
	nullByte := []byte("\x00")
	header := fmt.Sprintf("%s %d%s", objectType, byteSize, string(nullByte))
	fullData := append([]byte(header), data...)

	// Hash full data
	hasher := sha1.New()
	hasher.Write(fullData)
	sha1Hash := hex.EncodeToString(hasher.Sum(nil))

	if writeMode {
		objDir := filepath.Join(".gogit", "objects", sha1Hash[:2])
		fileName := sha1Hash[2:]
		objFullPath := filepath.Join(objDir, fileName)

		if _, err := os.Stat(objDir); os.IsNotExist(err) {
			err := os.Mkdir(objDir, os.ModeDir)
			if err != nil {
				return "", errors.New("object directory couldn't be created")
			}
		}

		f, err := os.Create(objFullPath)
		if err != nil {
			return "", err
		}
		defer f.Close()

		zlibComprData := compressData(fullData)
		cnt, err := f.Write(zlibComprData)
		if err != nil {
			return "", err
		}
		fmt.Printf("%v bytes written to object database\n", cnt)

	}

	return sha1Hash, nil
}

func getNullCharacterIndex(data []byte) (int, error) {
	for i, v := range data {
		if v == 0 {
			return i, nil
		}
	}
	return 0, errors.New("no null character found")
}

func readObject(sha1Hash string) ([]byte, error) {
	fp := filepath.Join(".gogit", "objects", sha1Hash[:2], sha1Hash[2:])
	fmt.Println(fp)
	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	//fileInfo, _ := file.Stat()
	//var size int64 = fileInfo.Size()

	reader := bufio.NewReader(file)
	rawData, err := ioutil.ReadAll(reader)
	file.Close()
	if err != nil {
		return nil, err
	}

	data, err := decompressData(rawData)
	if err != nil {
		return nil, err
	}

	return data, nil
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
		sha1Hash, err := writeObject(testData, "blob", false)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(sha1Hash)
	} else if cmd[0] == "write" {
		sha1Hash, err := writeObject(testData, "blob", true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(sha1Hash)
	} else if cmd[0] == "read" {
		sha1Hash, err := writeObject(testData, "blob", false)
		if err != nil {
			fmt.Println(err)
		}
		data, err := readObject(sha1Hash)
		if err != nil {
			fmt.Println(err)
		}
		nullIndex, err := getNullCharacterIndex(data)
		if err != nil {
			fmt.Println(err)
		}
		header := strings.Fields(string(data[:nullIndex]))
		fmt.Printf("Object Type: %s\n", header[0])
		fmt.Printf("Byte Length: %s\n", header[1])
		fmt.Print("\nData:\n")
		fmt.Println(string(data[nullIndex:]))

	} else {
		fmt.Printf("unknown subcommand: %v\n", string(cmd[0]))
	}
}
