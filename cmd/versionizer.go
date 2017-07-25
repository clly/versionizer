package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/clly/versionizer"
	"github.com/clly/versionizer/cmd/version"
)

func main() {
	repo := flag.String("repo", "", "Git repository location")
	commits := flag.Int("commits", 10, "Number of commits to go back in the git log for scmlog")
	pkg := flag.String("pkg", "version", "Package name for version function to come out at")
	output := flag.String("output", "versionizer.go", "File name for output go file")
	outputdir := flag.String("output-dir", ".", "Output directory for versionizer.go")
	showversion := flag.Bool("version", false, "Show version string")
	//_ = flag.String("version", "", "SemVer formatted string (not checked to ensure it matches SemVer formatting constraints")
	flag.Parse()

	if *showversion {
		fmt.Fprintf(os.Stderr, version.VersionString())
		os.Exit(0)
	}

	if *repo == "" {
		flag.Usage()
		oopsf("repo argument is empty: %s", *repo)
	}

	g, err := versionizer.GetGit(*repo, *commits)
	if err != nil {
		oops(err)
	}

	path, err := getAndMakeFinalPath(*output, *outputdir, *pkg)
	if err != nil && !os.IsExist(err) {
		oopsf("Failed to create path %s: %s", path, err)
	}

	f, err := os.Create(path)
	if err != nil {
		oopsf("Failed to open %s : %s", path, err)
	}
	defer f.Close()
	vc := fmt.Sprintf(versionizerContent, *pkg, time.Now().String(), g.Hash)

	_, err = f.WriteString(vc)
	if err != nil {
		oopsf("Failed to write versionizer context to %s : %s", path, err)
	}
	err = f.Sync()
	if err != nil {
		oopsf("Failed to sync %s : %s", path, err)
	}
	os.Exit(0)

}

func getAndMakeFinalPath(output, outputdir, pkg string) (string, error) {
	b := path.Base(output)
	p := path.Join(path.Clean(outputdir), pkg)
	return path.Join(p, b), os.Mkdir(p, 0766)
}

func oopsf(formatter string, args ...interface{}) {
	log.Printf(formatter, args)
	os.Exit(1)
}

func oops(msg error) {
	log.Print(msg)
	os.Exit(1)
}

var versionizerContent = "// author versionizer\n" +
	"package %s\n" +
	"import \"runtime\"\n" +
	"import \"fmt\"\n" +
	"\n" +
	"func VersionString() string {\n" +
	"\to := fmt.Sprintf(\"Built with %%s at %s at git hash %s\\n\", runtime.Version())\n" +
	"\treturn o\n" +
	"}\n"
