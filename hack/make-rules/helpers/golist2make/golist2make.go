/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"encoding/json"
	goflag "flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

var flDebug = pflag.BoolP("debug", "d", false, "enable debugging output")
var flHelp = pflag.BoolP("help", "h", false, "print help and exit")
var flPrefix = pflag.StringP("prefix", "p", "", "prefix make rule names with this string")
var flVendor = pflag.BoolP("prefix-nonroot-with-vendor", "v", false, "prefix the package portion of make rule names with ${rootModule}/vendor/ for packages outside the root module")
var flOnly = pflag.StringSlice("only", nil, "limit generated rules to packages beginning with these prefixes (may be specified multiple times)")

func main() {
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	pflag.Usage = func() { help(os.Stderr) }
	pflag.Parse()

	if *flHelp {
		help(os.Stdout)
		os.Exit(0)
	}
	if len(pflag.Args()) == 0 {
		help(os.Stderr)
		os.Exit(1)
	}
	sawNonFlagArg := ""
	for _, arg := range pflag.Args() {
		if strings.HasPrefix(arg, "-") {
			if len(sawNonFlagArg) > 0 {
				panic(fmt.Sprintf("flag args must come before all non-flag args, saw %q after %q", arg, sawNonFlagArg))
			}
			continue
		} else {
			sawNonFlagArg = arg
		}
		// guard against things that don't work consistently with go list across gopath and gomodule modes
		if strings.HasPrefix(arg, "./") {
			panic(fmt.Sprintf("packages must not be relative, got %q", arg))
		}
		if strings.HasSuffix(arg, "/...") {
			panic(fmt.Sprintf("packages must not be recursive, got %q", arg))
		}
	}

	// GO111MODULE=on go list -json -deps -e ...
	cmd := exec.Command("go", append([]string{"list", "-json", "-deps", "-e"}, pflag.Args()...)...)
	cmd.Env = append(os.Environ(), "GO111MODULE=on")
	cmd.Stderr = os.Stderr
	data, err := cmd.Output()
	if err != nil {
		panic(fmt.Sprintf("error running go list: %v", err))
	}

	// map of import path to package info
	packages := map[string]Package{}
	// list of import paths seen
	imports := []string{}

	rootModule := ""
	d := json.NewDecoder(bytes.NewBuffer(data))
	for {
		pkg := Package{}
		err := d.Decode(&pkg)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		if len(*flOnly) > 0 {
			found := false
			for _, prefix := range *flOnly {
				if strings.HasPrefix(pkg.ImportPath, prefix) {
					found = true
					break
				}
			}
			if !found {
				debug("skipping", pkg.ImportPath)
				continue
			}
		}

		if _, seen := packages[pkg.ImportPath]; !seen {
			imports = append(imports, pkg.ImportPath)
		}
		packages[pkg.ImportPath] = pkg
		if pkg.Module != nil && pkg.Module.Main {
			rootModule = pkg.Module.Path
		}
	}

	if len(imports) == 0 {
		panic("no packages found")
	}

	sort.Strings(imports)

	cachedDeps := map[string][]string{}

	for _, pkgPath := range imports {
		// ignore stdlib packages
		if packages[pkgPath].Standard {
			continue
		}

		if packages[pkgPath].Error != nil {
			debug("error loading package", pkgPath, packages[pkgPath].Error.Err)
			continue
		}

		deps := depsForPkg(pkgPath, cachedDeps, packages)
		makeRulePrefix := *flPrefix
		if len(rootModule) > 0 && *flVendor && !strings.HasPrefix(pkgPath, rootModule) {
			makeRulePrefix += rootModule + "/vendor/"
		}
		fmt.Printf("%s%s := \\\n  %s\n", makeRulePrefix, pkgPath, strings.Join(deps, " \\\n  "))
	}
}

func depsForPkg(pkgPath string, cachedDeps map[string][]string, packages map[string]Package) []string {
	if deps, ok := cachedDeps[pkgPath]; ok {
		return deps
	}

	pkg, ok := packages[pkgPath]
	if !ok || pkg.Standard {
		return nil
	}

	// Packages depend on their own directories, their own files, and a transitive list of all deps' directories and files.
	deps := make(map[string]struct{}, 1+len(pkg.Imports)+len(pkg.GoFiles))
	for _, i := range pkg.Imports {
		for _, dep := range depsForPkg(i, cachedDeps, packages) {
			deps[dep] = struct{}{}
		}
	}
	// enumerate specific go files instead of *.go to exclude autogenerated files
	for _, f := range pkg.GoFiles {
		deps[pkg.Dir+"/"+f] = struct{}{}
	}
	deps[pkg.Dir+"/"] = struct{}{}

	orderedDeps := make([]string, 0, len(deps))
	for dep := range deps {
		orderedDeps = append(orderedDeps, dep)
	}
	sort.Strings(orderedDeps)
	cachedDeps[pkgPath] = orderedDeps
	return orderedDeps
}

func help(out io.Writer) {
	fmt.Fprintf(out, "Usage: %s [FLAG...] <PKG...>\n", os.Args[0])
	fmt.Fprintf(out, "\n")
	fmt.Fprintf(out, "go2make calculates all of the dependencies of a set of Go packages and prints\n")
	fmt.Fprintf(out, "them as variable definitions suitable for use as a Makefile.\n")
	fmt.Fprintf(out, "\n")
	fmt.Fprintf(out, "Package specifications may be simple (e.g. 'example.com/txt/color') or\n")
	fmt.Fprintf(out, "recursive (e.g. 'example.com/txt/...')\n")
	fmt.Fprintf(out, " Example:\n")
	fmt.Fprintf(out, "  $ %s example.com/pretty\n", os.Args[0])
	fmt.Fprintf(out, "  example.com/txt/split := \\\n")
	fmt.Fprintf(out, "    /go/src/example.com/txt/split/ \\\n")
	fmt.Fprintf(out, "    /go/src/example.com/txt/split/split.go \\\n")
	fmt.Fprintf(out, "  example.com/pretty := \\\n")
	fmt.Fprintf(out, "    /go/src/example.com/pretty/ \\\n")
	fmt.Fprintf(out, "    /go/src/example.com/pretty/print.go \\\n")
	fmt.Fprintf(out, "    /go/src/example.com/txt/split\n")
	fmt.Fprintf(out, "\n")
	fmt.Fprintf(out, " Flags:\n")

	pflag.PrintDefaults()
}

func debug(items ...interface{}) {
	if *flDebug {
		x := []interface{}{"DBG:"}
		x = append(x, items...)
		fmt.Fprintln(os.Stderr, x...)
	}
}

type Package struct {
	Dir           string   // directory containing package sources
	ImportPath    string   // import path of package in dir
	ImportComment string   // path in import comment on package statement
	Name          string   // package name
	Doc           string   // package documentation string
	Target        string   // install path
	Shlib         string   // the shared library that contains this package (only set when -linkshared)
	Goroot        bool     // is this package in the Go root?
	Standard      bool     // is this package part of the standard Go library?
	Stale         bool     // would 'go install' do anything for this package?
	StaleReason   string   // explanation for Stale==true
	Root          string   // Go root or Go path dir containing this package
	ConflictDir   string   // this directory shadows Dir in $GOPATH
	BinaryOnly    bool     // binary-only package (no longer supported)
	ForTest       string   // package is only for use in named test
	Export        string   // file containing export data (when using -export)
	BuildID       string   // build ID of the compiled package (when using -export)
	Module        *Module  // info about package's containing module, if any (can be nil)
	Match         []string // command-line patterns matching this package
	DepOnly       bool     // package is only a dependency, not explicitly listed

	// Source files
	GoFiles           []string // .go source files (excluding CgoFiles, TestGoFiles, XTestGoFiles)
	CgoFiles          []string // .go source files that import "C"
	CompiledGoFiles   []string // .go files presented to compiler (when using -compiled)
	IgnoredGoFiles    []string // .go source files ignored due to build constraints
	IgnoredOtherFiles []string // non-.go source files ignored due to build constraints
	CFiles            []string // .c source files
	CXXFiles          []string // .cc, .cxx and .cpp source files
	MFiles            []string // .m source files
	HFiles            []string // .h, .hh, .hpp and .hxx source files
	FFiles            []string // .f, .F, .for and .f90 Fortran source files
	SFiles            []string // .s source files
	SwigFiles         []string // .swig files
	SwigCXXFiles      []string // .swigcxx files
	SysoFiles         []string // .syso object files to add to archive
	TestGoFiles       []string // _test.go files in package
	XTestGoFiles      []string // _test.go files outside package

	// Embedded files
	EmbedPatterns      []string // //go:embed patterns
	EmbedFiles         []string // files matched by EmbedPatterns
	TestEmbedPatterns  []string // //go:embed patterns in TestGoFiles
	TestEmbedFiles     []string // files matched by TestEmbedPatterns
	XTestEmbedPatterns []string // //go:embed patterns in XTestGoFiles
	XTestEmbedFiles    []string // files matched by XTestEmbedPatterns

	// Cgo directives
	CgoCFLAGS    []string // cgo: flags for C compiler
	CgoCPPFLAGS  []string // cgo: flags for C preprocessor
	CgoCXXFLAGS  []string // cgo: flags for C++ compiler
	CgoFFLAGS    []string // cgo: flags for Fortran compiler
	CgoLDFLAGS   []string // cgo: flags for linker
	CgoPkgConfig []string // cgo: pkg-config names

	// Dependency information
	Imports      []string          // import paths used by this package
	ImportMap    map[string]string // map from source import to ImportPath (identity entries omitted)
	Deps         []string          // all (recursively) imported dependencies
	TestImports  []string          // imports from TestGoFiles
	XTestImports []string          // imports from XTestGoFiles

	// Error information
	Incomplete bool            // this package or a dependency has an error
	Error      *PackageError   // error loading package
	DepsErrors []*PackageError // errors loading dependencies
}

type PackageError struct {
	ImportStack []string // shortest path from package named on command line to this one
	Pos         string   // position of error (if present, file:line:col)
	Err         string   // the error itself
}

type Module struct {
	Path      string       // module path
	Version   string       // module version
	Versions  []string     // available module versions (with -versions)
	Replace   *Module      // replaced by this module
	Time      *time.Time   // time version was created
	Update    *Module      // available update, if any (with -u)
	Main      bool         // is this the main module?
	Indirect  bool         // is this module only an indirect dependency of main module?
	Dir       string       // directory holding files for this module, if any
	GoMod     string       // path to go.mod file used when loading this module, if any
	GoVersion string       // go version used in module
	Retracted string       // retraction information, if any (with -retracted or -u)
	Error     *ModuleError // error loading module
}

type ModuleError struct {
	Err string // the error itself
}
