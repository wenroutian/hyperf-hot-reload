package main

import (
	"flag"
	"fmt"
	"github.com/andreaskoch/go-fswatch"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var path = flag.String("path", "", "the hyperf path to check")

func main() {
	flag.Parse()
	p := *path

	if p == "" {
		log.Fatal("fail get the restart path")
	}

	_, err := os.Stat(p)

	if err == os.ErrNotExist {
		log.Fatal("path not exist")
	}

	fmt.Println(p)

	go func() {

		skipDotFilesAndFolders := func(path string) bool {
			return strings.HasPrefix(filepath.Base(path), ".")
		}

		checkIntervalInSeconds := 1

		folderWatcher := fswatch.NewFolderWatcher(
			p+"/app",
			true,
			skipDotFilesAndFolders,
			checkIntervalInSeconds,
		)

		folderWatcher.Start()

		for folderWatcher.IsRunning() {
			select {
			case <-folderWatcher.Modified():
				fmt.Println("file changed")
				killCurrent(p)
				start(p)
			}
		}

	}()

	select {}
}

func fileExist(f string) bool {
	_, err := os.Stat(f)
	return err != os.ErrNotExist
}

func killCurrent(p string) {
	rP := "/runtime/hyperf.pid"
	path := p + rP
	if fileExist(path) {
		f, _ := os.Open(path)
		pid, _ := ioutil.ReadAll(f)
		execute("kill", "-9", string(pid))
	}
}

func start(p string) {
	execute("php", p+"/bin/hyperf.php", "start")
}

func execute(c string, args ...string) {
	err := exec.Command(c, args...).Run()
	if err != nil {
		log.Println(err)
	}
}
