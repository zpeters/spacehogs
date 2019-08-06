package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

var html bool
var files Files
var path string
var topNumber int

const Version = "0.2"

type File struct {
	Path string
	Size int
}

type Files []*File

func (s Files) Len() int      { return len(s) }
func (s Files) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type BySize struct{ Files }

func (s BySize) Less(i, j int) bool { return s.Files[i].Size < s.Files[j].Size }

type Reverse struct{ Files }

func (s Files) Reverse() {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func WalkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		//log.Println(err)
	} else {
		f := &File{path, int(info.Size())}
		files = append(files, f)
	}
	return nil
}

func humanizeBytes(bytes int) string {
	var displayString string

	if bytes > 1073741824 {
		displayString = fmt.Sprintf("%0.2f Gigs", float64(bytes)/1024/1024/1024)
	} else if bytes > 1048576 {
		displayString = fmt.Sprintf("%0.2f Megs", float64(bytes)/1024/1024)
	} else {
		displayString = fmt.Sprintf("%d bytes")
	}

	return displayString
}

func usage() {
	fmt.Printf("Usage: %s -n [Number of resutls] -p [Path]\n", os.Args[0])
	fmt.Printf("Usage: %s -n [Number of resutls] -p [Path] -html \t\t HTML rendering\n", os.Args[0])
	fmt.Printf("Usage: %s -h\t - Display Help\n", os.Args[0])
	return
}

func main() {
	t0 := time.Now()
	files = []*File{}

	// get arguments
	var topNumber = flag.Int("n", 0, "Number of largest files you want to see")
	var path = flag.String("p", "YOUR PATH", "Path you want to crawl")
	var html = flag.Bool("html", false, "Render HTML")
	var help = flag.Bool("h", false, "Show Help")
	var version = flag.Bool("v", false, "Display Version")

	flag.Parse()

	if *version != false {
		fmt.Printf("Spacehogs - Version: %s\n", Version)
		return
	}

	if *help != false {
		usage()
		return
	}

	if *topNumber == 0 && *path == "YOUR PATH" {
		usage()
		return
	}

	err := filepath.Walk(*path, WalkFunc)
	if err != nil {
		log.Fatal(err)
	}
	t1 := time.Now()

	sort.Sort(BySize{files})
	files.Reverse()

	for i := 0; i < *topNumber; i++ {
		f := files[i]
		if *html {
			fmt.Printf("%s - %s<br>", f.Path, humanizeBytes(f.Size))
		} else {
			fmt.Printf("%s - %s\n", f.Path, humanizeBytes(f.Size))
		}
	}

	if *html {
		fmt.Printf("-----<br>")
		fmt.Printf("Took %v to find files<br>", t1.Sub(t0))
	} else {
		fmt.Printf("-----\n")
		fmt.Printf("Took %v to find files\n", t1.Sub(t0))
	}
}
