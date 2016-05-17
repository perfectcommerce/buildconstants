package main

import (
	"bufio"
	"bytes"
	"flag"
	"go/format"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
)

const (
	defaultfilename = "buildversion_generated.go"

	templatePrelude = `// Generated by running
//      buildversion
// DO NOT EDIT

package {{.GOPKG}}

const (
	MustRunBuildVersion = 0
`

	templatePostlude = ")"
)

type cmd struct {
	Var  string
	Line string
	echo bool // get value from env
}

// read commands from a text file.  Each command has its own line
//   and the lines are of the form `VAR = echo $FOO`
//   where $VAR is the env variable and everything after '=' will be evaluated for the value
func cmdRead(cmdfile string) ([]cmd, error) {

	f, err := os.Open(cmdfile)
	if err != nil {
		return []cmd{}, err
	}
	defer f.Close()
	var cmds []cmd
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		splits := strings.Split(scanner.Text(), "=")
		if len(splits) < 2 {
			continue
		}
		line := strings.TrimSpace(strings.Join(splits[1:], " "))
		c := cmd{
			Var:  strings.TrimSpace(splits[0]),
			Line: line,
			echo: strings.HasPrefix(line, "$"),
		}
		cmds = append(cmds, c)
	}
	return cmds, nil
}

// operates on inputs since it needs to be in the same order
func envTemplate(ins []cmd) string {
	var s string
	for _, in := range ins {
		s = s + " " + in.Var + " = \"{{." + in.Var + "}}\"\n"
	}
	return s
}

func do(command cmd) (string, error) {
	expanded := os.Expand(command.Line, os.Getenv)
	var out string
	if command.echo {
		out = expanded
	} else {
		split := strings.Split(expanded, " ")
		cmd := exec.Command(split[0], split[1:]...)
		bout, err := cmd.Output()
		if err != nil {
			return "", err
		}
		out = string(bout)
	}

	val := strings.TrimSpace(out)
	err := os.Setenv(command.Var, val)
	if err != nil {
		return "", err
	}
	return val, nil
}

func pkg() string {
	pkgCmd := exec.Command("go", "list", ".")
	out, err := pkgCmd.Output()
	if err != nil {
		return "main"
	}
	_, packageName := path.Split(strings.TrimSpace(string(out)))
	return packageName
}

func main() {

	cmdfile := flag.String("i", "commands.txt", "output file")
	fname := flag.String("o", defaultfilename, "output file")
	packageName := flag.String("package", pkg(), "package the generated file will be in.")
	flag.Parse()

	// incmds := []cmd{
	// 	cmd{Var: "GITVERSION", Line: "git rev-list --tags --max-count=1"},
	// 	cmd{Var: "GITTAG", Line: "git describe --always --tags ${GITVERSION}"},
	// 	cmd{Var: "GOVERSION", Line: "go version"},
	// 	cmd{Var: "BUILD_NUMBER", echo: true},
	// 	cmd{Var: "BRANCH_NAME", echo: true},
	// }
	incmds, err := cmdRead(*cmdfile)
	if err != nil {
		os.Exit(1)
	}

	outcmds := make(map[string]string, len(incmds)+1) //since we add GOPKG
	outcmds["GOPKG"] = *packageName
	if err != nil {
		os.Exit(1)
	}

	for _, cmd := range incmds {
		val, err := do(cmd)
		if err != nil {
			val = ""
			// os.Exit(1)
		}
		outcmds[cmd.Var] = val
	}

	templ := templatePrelude + envTemplate(incmds) + templatePostlude

	t := template.Must(template.New("templ").Parse(templ))

	var buf bytes.Buffer
	err = t.Execute(&buf, outcmds)
	if err != nil {
		os.Exit(1)
	}
	fmted, err := format.Source(buf.Bytes())
	if err != nil {
		os.Exit(1)
	}

	f, err := os.Create(*fname)
	if err != nil {
		os.Exit(1)
	}
	defer f.Close()
	f.Write(fmted)
}
