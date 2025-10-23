package main

import (
	"fmt"
	"os"
	"slices"

	_ "net/http/pprof"

	"github.com/DmitryNaumov/goplspd/gopackages"
	"github.com/goccy/go-json"
	"golang.org/x/tools/go/packages"
)

func main() {
	var (
		err error
		req packages.DriverRequest
		dr  DriverResponse
	)

	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		panic(err)
	}

	env, err := gopackages.ParseEnv(req.Env)
	if err != nil {
		panic(err)
	}

	fmt.Println("env:", env)
	fmt.Println("args:", os.Args)

	if !slices.Contains(os.Args, "./...") {
		dr.NotHandled = true
		writeResponse(&dr)
	}

	w := gopackages.NewWalker(env, env.Targets)

	dr.GoVersion = env.MinorVersion()
	dr.Arch = env.GOARCH
	dr.Compiler = "gc"

	dr.Packages, err = w.Packages()
	if err != nil {
		panic(err)
	}

	for _, p := range dr.Packages {
		if !p.DepOnly {
			dr.Roots = append(dr.Roots, p.ID)
		}
	}

	dr.Roots = append(dr.Roots, "builtin")

	writeResponse(&dr)
}

func writeResponse(dr *DriverResponse) {
	if err := json.NewEncoder(os.Stdout).Encode(dr); err != nil {
		panic(err)
	}

	os.Exit(0)
}

type DriverResponse struct {
	NotHandled bool

	Compiler string
	Arch     string

	Roots []string `json:",omitempty"`

	Packages []*gopackages.Package

	GoVersion int
}
