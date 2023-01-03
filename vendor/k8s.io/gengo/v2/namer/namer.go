/*
Copyright 2015 The Kubernetes Authors.

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

package namer

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"k8s.io/gengo/v2/types"
)

const (
	// GoSeperator is used to split go import paths.
	// Forward slash is used instead of filepath.Seperator because it is the
	// only universally-accepted path delimiter and the only delimiter not
	// potentially forbidden by Go compilers. (In particular gc does not allow
	// the use of backslashes in import paths.)
	// See https://golang.org/ref/spec#Import_declarations.
	// See also https://github.com/kubernetes/gengo/issues/83#issuecomment-367040772.
	GoSeperator = "/"
)

// Returns whether a name is a private Go name.
func IsPrivateGoName(name string) bool {
	return len(name) == 0 || strings.ToLower(name[:1]) == name[:1]
}

// NewPublicNamer is a helper function that returns a namer that makes
// CamelCase names. See the NameStrategy struct for an explanation of the
// arguments to this constructor.
func NewPublicNamer(prependPackageNames int, ignoreWords ...string) *NameStrategy {
	n := &NameStrategy{
		Join:                Joiner(IC, IC),
		IgnoreWords:         map[string]bool{},
		PrependPackageNames: prependPackageNames,
	}
	for _, w := range ignoreWords {
		n.IgnoreWords[w] = true
	}
	return n
}

// NewPrivateNamer is a helper function that returns a namer that makes
// camelCase names. See the NameStrategy struct for an explanation of the
// arguments to this constructor.
func NewPrivateNamer(prependPackageNames int, ignoreWords ...string) *NameStrategy {
	n := &NameStrategy{
		Join:                Joiner(IL, IC),
		IgnoreWords:         map[string]bool{},
		PrependPackageNames: prependPackageNames,
	}
	for _, w := range ignoreWords {
		n.IgnoreWords[w] = true
	}
	return n
}

// NewRawNamer will return a Namer that makes a name by which you would
// directly refer to a type, optionally keeping track of the import paths
// necessary to reference the names it provides. Tracker may be nil.
// The 'pkg' is the full package name, in which the Namer is used - all
// types from that package will be referenced by just type name without
// referencing the package.
//
// For example, if the type is map[string]int, a raw namer will literally
// return "map[string]int".
//
// Or if the type, in package foo, is "type Bar struct { ... }", then the raw
// namer will return "foo.Bar" as the name of the type, and if 'tracker' was
// not nil, will record that package foo needs to be imported.
func NewRawNamer(pkg string, tracker ImportTracker) *rawNamer {
	return &rawNamer{pkg: pkg, tracker: tracker}
}

// Names is a map from Type to name, as defined by some Namer.
type Names map[*types.Type]string

// Namer takes a type, and assigns a name.
//
// The purpose of this complexity is so that you can assign coherent
// side-by-side systems of names for the types. For example, you might want a
// public interface, a private implementation struct, and also to reference
// literally the type name.
//
// Note that it is safe to call your own Name() function recursively to find
// the names of keys, elements, etc. This is because anonymous types can't have
// cycles in their names, and named types don't require the sort of recursion
// that would be problematic.
type Namer interface {
	Name(*types.Type) string
}

// NameStrategy is a general Namer. The easiest way to use it is to copy the
// Public/PrivateNamer variables, and modify the members you wish to change.
//
// The Name method produces a name for the given type, of the forms:
// Anonymous types: <Prefix><Type description><Suffix>
// Named types: <Prefix><Optional Prepended Package name(s)><Original name><Suffix>
//
// In all cases, every part of the name is run through the capitalization
// functions.
//
// The IgnoreWords map can be set if you have directory names that are
// semantically meaningless for naming purposes, e.g. "proto".
//
// Prefix and Suffix can be used to disambiguate parallel systems of type
// names. For example, if you want to generate an interface and an
// implementation, you might want to suffix one with "Interface" and the other
// with "Implementation". Another common use-- if you want to generate private
// types, and one of your source types could be "string", you can't use the
// default lowercase private namer. You'll have to add a suffix or prefix.
type NameStrategy struct {
	Prefix, Suffix string
	Join           func(parts ...string) string

	// Add non-meaningful package directory names here (e.g. "proto") and
	// they will be ignored.
	IgnoreWords map[string]bool

	// If > 0, prepend exactly that many package directory names (or as
	// many as there are).  Package names listed in "IgnoreWords" will be
	// ignored.
	//
	// For example, if Ignore words lists "proto" and type Foo is in
	// pkg/server/frobbing/proto, then a value of 1 will give a type name
	// of FrobbingFoo, 2 gives ServerFrobbingFoo, etc.
	PrependPackageNames int

	// A cache of names thus far assigned by this namer.
	Names
}

// IC ensures the first character is uppercase.
func IC(in string) string {
	if in == "" {
		return in
	}
	return strings.ToUpper(in[:1]) + in[1:]
}

// IL ensures the first character is lowercase.
func IL(in string) string {
	if in == "" {
		return in
	}
	return strings.ToLower(in[:1]) + in[1:]
}

// Joiner lets you specify functions that preprocess the various components of
// a name before joining them. You can construct e.g. camelCase or CamelCase or
// any other way of joining words. (See the IC and IL convenience functions.)
func Joiner(first, others func(string) string) func(in ...string) string {
	return func(in ...string) string {
		tmp := make([]string, 0, len(in))
		for i := range in {
			tmp = append(tmp, others(in[i]))
		}
		return first(strings.Join(tmp, ""))
	}
}

var (
	importPathNameSanitizer = strings.NewReplacer("-", "_", ".", "")
)

// filters out unwanted directory names and sanitizes remaining names.
func (ns *NameStrategy) filterDirs(path string) []string {
	allDirs := strings.Split(path, GoSeperator)
	dirs := make([]string, 0, len(allDirs))
	for _, p := range allDirs {
		if ns.IgnoreWords == nil || !ns.IgnoreWords[p] {
			dirs = append(dirs, importPathNameSanitizer.Replace(p))
		}
	}
	return dirs
}

// See the comment on NameStrategy.
func (ns *NameStrategy) Name(t *types.Type) string {
	if ns.Names == nil {
		ns.Names = Names{}
	}
	if name, ok := ns.Names[t]; ok {
		return name
	}

	// Some implementations of Join are not no-op wrt empty parts.
	parts := make([]string, 0, 3)
	if ns.Prefix != "" {
		parts = append(parts, ns.Prefix)
	}
	parts = append(parts, ns.name(t))
	if ns.Suffix != "" {
		parts = append(parts, ns.Suffix)
	}

	name := ns.Join(parts...)
	ns.Names[t] = name
	return name
}

// name is like Name but without the prefix and suffix.
func (ns *NameStrategy) name(t *types.Type) string {
	if t.Name.Package != "" {
		dirs := append(ns.filterDirs(t.Name.Package), t.Name.Name)
		i := ns.PrependPackageNames + 1
		dn := len(dirs)
		if i > dn {
			i = dn
		}
		name := ns.Join(dirs[dn-i:]...)
		return name
	}

	// Only anonymous types remain.
	var name string
	switch t.Kind {
	case types.Builtin:
		name = t.Name.Name
	case types.Map:
		name = ns.Join(
			"Map",
			ns.name(t.Key),
			"To",
			ns.name(t.Elem))
	case types.Slice:
		name = ns.Join(
			"Slice",
			ns.name(t.Elem))
	case types.Array:
		name = ns.Join(
			"Array",
			fmt.Sprintf("%d", t.Len),
			ns.name(t.Elem))
	case types.Pointer:
		name = ns.Join(
			"Pointer",
			ns.name(t.Elem))
	case types.Struct:
		parts := []string{"Struct"}
		for _, m := range t.Members {
			parts = append(parts, ns.name(m.Type))
		}
		name = ns.Join(parts...)
	case types.Chan:
		name = ns.Join(
			"Chan",
			ns.name(t.Elem))
	case types.Interface:
		// TODO: add to name test
		if len(t.Methods) == 0 {
			name = "Any"
		} else if t.Name.Name == "error" {
			name = "Error"
		} else {
			parts := []string{"Interface"}
			for n, m := range t.Methods {
				parts = append(parts, ns.nameFunc(n, m.Signature))
			}
			name = ns.Join(parts...)
		}
	case types.Func:
		// TODO: add to name test
		name = ns.nameFunc("Func", t.Signature)
	default:
		name = "unnameable_" + string(t.Kind)
	}
	return name
}

// nameFunc returns a name string for a function or method.  The name argument
// can be used to indicate the method name when naming interfaces.
func (ns *NameStrategy) nameFunc(name string, sig *types.Signature) string {
	parts := make([]string, 0, 2+len(sig.Parameters)+len(sig.Results))
	parts = append(parts, name)
	for _, pt := range sig.Parameters {
		parts = append(parts, ns.name(pt))
	}
	if len(sig.Results) > 0 {
		parts = append(parts, "Returns")
		for _, rt := range sig.Results {
			parts = append(parts, ns.name(rt))
		}
	}
	return ns.Join(parts...)
}

// ImportTracker allows a raw namer to keep track of the packages needed for
// import. You can implement yourself or use the one in the generation package.
type ImportTracker interface {
	AddType(*types.Type)
	AddImport(packagePath string)
	LocalNameOf(packagePath string) string
	PathOf(localName string) (string, bool)
	ImportLines() []string
}

type rawNamer struct {
	pkg     string
	tracker ImportTracker
	Names
}

// Name makes a name the way you'd write it to literally refer to type t,
// making ordinary assumptions about how you've imported t's package (or using
// r.tracker to specifically track the package imports).
func (r *rawNamer) Name(t *types.Type) string {
	if r.Names == nil {
		r.Names = Names{}
	}
	if name, ok := r.Names[t]; ok {
		return name
	}
	if t.Name.Package != "" {
		var name string
		if t.Name.Package == r.pkg {
			name = t.Name.Name
		} else {
			if r.tracker != nil {
				r.tracker.AddType(t)
				name = r.tracker.LocalNameOf(t.Name.Package) + "." + t.Name.Name
			} else {
				name = filepath.Base(t.Name.Package) + "." + t.Name.Name
			}
		}
		r.Names[t] = name
		return name
	}
	var name string
	switch t.Kind {
	case types.Builtin:
		name = t.Name.Name
	case types.Map:
		name = "map[" + r.Name(t.Key) + "]" + r.Name(t.Elem)
	case types.Slice:
		name = "[]" + r.Name(t.Elem)
	case types.Array:
		l := strconv.Itoa(int(t.Len))
		name = "[" + l + "]" + r.Name(t.Elem)
	case types.Pointer:
		name = "*" + r.Name(t.Elem)
	case types.Struct:
		elems := []string{}
		for _, m := range t.Members {
			elems = append(elems, m.Name+" "+r.Name(m.Type))
		}
		name = "struct{" + strings.Join(elems, "; ") + "}"
	case types.Chan:
		// TODO: include directionality
		name = "chan " + r.Name(t.Elem)
	case types.Interface:
		// TODO: add to name test
		elems := []string{}
		for _, m := range t.Methods {
			// TODO: include function signature
			elems = append(elems, m.Name.Name)
		}
		name = "interface{" + strings.Join(elems, "; ") + "}"
	case types.Func:
		// TODO: add to name test
		params := []string{}
		for _, pt := range t.Signature.Parameters {
			params = append(params, r.Name(pt))
		}
		results := []string{}
		for _, rt := range t.Signature.Results {
			results = append(results, r.Name(rt))
		}
		name = "func(" + strings.Join(params, ",") + ")"
		if len(results) == 1 {
			name += " " + results[0]
		} else if len(results) > 1 {
			name += " (" + strings.Join(results, ",") + ")"
		}
	default:
		name = "unnameable_" + string(t.Kind)
	}
	r.Names[t] = name
	return name
}
