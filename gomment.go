package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
)

var comment = flag.Bool("comments", true, "add explation comment add end of the line")

type errorType struct {
	VarName string
	ErrType string
}

func main() {
	flag.Parse()
	gofile := flag.Args()[0]
    tmpFile := path.Join(os.TempDir(), "gomment-tmp.go")

	errors := goBuild(gofile)
	ret, err := goFix(gofile, errors, []string{"import", "noNewVariable"})
	if err != nil {
		log.Fatal(err)
	}
	err = writeLines(ret, tmpFile)
	if err != nil {
		log.Fatal(err)
	}

	errors = goBuild(tmpFile)
	ret, err = goFix(tmpFile, errors, []string{"unusedVariable"})
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range ret {
		fmt.Print(line)
	}

    os.Remove(tmpFile)
    os.Remove(path.Join(os.TempDir(), "gomment-tmp"))

}

func goFix(filename, errors string, methods []string) (ret []string, err error) {

	//                               filename:linenr:  explanation : "package name"
	reImportError := regexp.MustCompile(`^([^:]+):(\d+): imported and not used: "([^"]+)"$`)
	reUnusedVar := regexp.MustCompile(`^([^:]+):(\d+): (\w+) declared and not used$`)
	reNoNewVariable := regexp.MustCompile(`^([^:]+):(\d+): no new variables on left side of :=$`)

	// map of linenumbers with with optional corrosponding offending variable
	var badlines = make(map[int]errorType)
	for _, method := range methods {
		for _, line := range strings.Split(errors, "\n") {
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			var ie []string
			var varname string
			switch method {
			case "import":
				ie = reImportError.FindStringSubmatch(strings.TrimRight(line, "\n"))
				varname = ""
			case "noNewVariable":
				ie = reNoNewVariable.FindStringSubmatch(strings.TrimRight(line, "\n"))
				varname = ""
			case "unusedVariable":
				ie = reUnusedVar.FindStringSubmatch(strings.TrimRight(line, "\n"))
				varname = ie[3]
			}

			if len(ie) == 0 {
				// There's an error that's not a comment or import error; we therefore can't fix it
				continue
			}
			lineNr, _ := strconv.Atoi(ie[2]) // make int of linenumber
			badlines[lineNr] = errorType{varname, method}
		}
	}
	inf := path.Clean(filename)
	ret, err = goCommentLines(inf, badlines)
	if err != nil {
		return
	}
	return
}

func goCommentLines(in string, lines map[int]errorType) (ret []string, err error) {
	var file *os.File
	if file, err = os.Open(in); err != nil {
		return nil, err
	}
	defer file.Close()

	binf := bufio.NewReader(file)
    alreadyCommented := regexp.MustCompile(`//(commented out unused .*)|(replaed := with =)$`)
	lineno := 0
	var line string

	for line, err = binf.ReadString('\n'); err == nil; line, err = binf.ReadString('\n') {
		lineno++
		if errTyp, ok := lines[lineno]; ok {
			line = strings.TrimRight(line, "\n")
            reSimpleDec1 := `(\s*var\s*` + errTyp.VarName + `\s*\w+)`
            reSimpleDec2 := `(\s*` + errTyp.VarName + `\s*:=.*)`
			reSimpleVarDec := regexp.MustCompile(reSimpleDec1 + "|" + reSimpleDec2)

            switch errTyp.ErrType{
            case "import":
				line = "// " + line
				if !alreadyCommented.MatchString(line) && *comment{
					line += " //commented out unused import"
				}

            case "noNewVariable":
                isRegex := regexp.MustCompile(":=")
                line = isRegex.ReplaceAllString(line, "=")
                if ! alreadyCommented.MatchString(line) && *comment {
                    line += " //replaced := with ="
                }
            case "unusedVariable":
                if reSimpleVarDec.MatchString(line) {
                    line = "// " + line
                    if !alreadyCommented.MatchString(line) && *comment {
                        line += " //commented unused variable"
                    }
                    break
                }

                varRegex := regexp.MustCompile(`\b` + errTyp.VarName + `\b`)
                line = varRegex.ReplaceAllString(line, "_")
                if !alreadyCommented.MatchString(line) && *comment {
                    line += " //commented unused variable " + errTyp.VarName
                }
            }
			line += "\n"
        }
        ret = append(ret, line)
	}
	if err == io.EOF {
		return ret, nil
	}
	return
}

func goBuild(gofile string) string {
    os.Chdir(path.Dir(gofile))
	args := append([]string{"build"}, gofile)
	build := exec.Command("go", args...)
	stdout, err := build.Output()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return string(stdout)
		}
		log.Fatal("goBuild", err)
	}
	return ""
}

func printFile(filename string) (err error) {
	fd, err := os.Open(filename)
	if err != nil {
		return
	}
	defer fd.Close()
	fdReader := bufio.NewReader(fd)

	for line, err := fdReader.ReadString('\n'); err == nil; line, err = fdReader.ReadString('\n') {
		fmt.Print(line)
	}
	if err == io.EOF {
		return nil
	}
	return
}

func writeLines(lines []string, path string) (err error) {
	var file *os.File

	if file, err = os.Create(path); err != nil {
		return
	}
	defer file.Close()

	for _, item := range lines {
		_, err := file.WriteString(item)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	return
}
