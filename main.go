package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	files := listFilesInPath(os.Args[1])
	dupes := findDuplicates(files)
	printDupes(dupes)
}

type File struct {
	path string
	size int64
	md5 string
}

func (file *File) Md5() string {
	if(file.md5 == "") {
		file.md5 = file.getMd5()
	}

	return file.md5
}

func (file File) Path() string {
	return file.path
}

func (file File) Size() int64 {
	return file.size
}

func (file File) getMd5() string {
	f, err := os.Open(file.path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return hex.EncodeToString(h.Sum(nil))
}

func printDupes(dupes  map[int64][]File) {
	matched := findDuplicateMD5(dupes)
	for _, matchFiles := range matched {
		if len(matchFiles) > 1 {
			fmt.Println("Found Matches")
			for _, file := range matchFiles {
				fmt.Println(file.Md5() + ": " + file.path)
			}
		}
	}
}

func findDuplicateMD5(dupes map[int64][]File) map[string][]File {
	var matched = make(map[string][]File)
	for _, dupe := range dupes {
		for _, file := range dupe {
			matched[file.Md5()] = append(matched[file.Md5()], file)
		}
	}

	return matched
}

func findDuplicates(files []File) map[int64][]File {
	var dupes = make(map[int64][]File)
	for _, file := range files {
		dupes[file.size] = append(dupes[file.size], file)
	}
	return dupes
}

func listFilesInPath(path string) []File {
	var files []File
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			if err != nil {
				return err
			}

			files = append(files, File{ path: path, size: info.Size() })
			return nil
		})

	if err != nil {
		log.Fatal(err)
	}

	return files
}

func getFileSize(filepath string) int64 {
	file, err := os.Stat(filepath)
	if err != nil {
		log.Fatal(err)
	}

	return file.Size()
}
