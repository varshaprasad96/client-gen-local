package custom

import (
	"fmt"
	"go/types"
	"path"
	"strings"
	"unicode"
	"unicode/utf8"

	"sigs.k8s.io/controller-tools/pkg/loader"
)

// ImportSpecs returns a string form of each import spec
// (i.e. `alias "path/to/import").  Aliases are only present
// when they don't match the package name.
func (l *importsList) ImportSpecs() []string {
	res := make([]string, 0, len(l.byPath))
	for importPath, alias := range l.byPath {
		pkg := l.pkg.Imports()[importPath]
		if pkg != nil && pkg.Name == alias {
			// don't print if alias is the same as package name
			// (we've already taken care of duplicates).
			res = append(res, fmt.Sprintf("%q", importPath))
		} else {
			res = append(res, fmt.Sprintf("%s %q", alias, importPath))
		}
	}
	return res
}

// NeedImport marks that the given package is needed in the list of imports,
// returning the ident (import alias) that should be used to reference the package.
func (l *importsList) NeedImport(importPath string) string {
	// we get an actual path from Package, which might include venddored
	// packages if running on a package in vendor.
	if ind := strings.LastIndex(importPath, "/vendor/"); ind != -1 {
		importPath = importPath[ind+8: /* len("/vendor/") */]
	}

	// check to see if we've already assigned an alias, and just return that.
	alias, exists := l.byPath[importPath]
	if exists {
		return alias
	}

	// otherwise, calculate an import alias by joining path parts till we get something unique
	restPath, nextWord := path.Split(importPath)

	for otherPath, exists := "", true; exists && otherPath != importPath; otherPath, exists = l.byAlias[alias] {
		if restPath == "" {
			// do something else to disambiguate if we're run out of parts and
			// still have duplicates, somehow
			alias += "x"
		}

		// can't have a first digit, per Go identifier rules, so just skip them
		for firstRune, runeLen := utf8.DecodeRuneInString(nextWord); unicode.IsDigit(firstRune); firstRune, runeLen = utf8.DecodeRuneInString(nextWord) {
			nextWord = nextWord[runeLen:]
		}

		// make a valid identifier by replacing "bad" characters with underscores
		nextWord = strings.Map(func(r rune) rune {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				return r
			}
			return '_'
		}, nextWord)

		alias = nextWord + alias
		if len(restPath) > 0 {
			restPath, nextWord = path.Split(restPath[:len(restPath)-1] /* chop off final slash */)
		}
	}

	l.byPath[importPath] = alias
	l.byAlias[alias] = importPath
	return alias
}

// namingInfo holds package and syntax for referencing a field, type,
// etc.  It's used to allow lazily marking import usage.
// You should generally retrieve the syntax using Syntax.
type namingInfo struct {
	// typeInfo is the type being named.
	typeInfo     types.Type
	nameOverride string
}

// Syntax calculates the code representation of the given type or name,
// and marks that is used (potentially marking an import as used).
func (n *namingInfo) Syntax(basePkg *loader.Package, imports *importsList) string {
	if n.nameOverride != "" {
		return n.nameOverride
	}

	// NB(directxman12): typeInfo.String gets us most of the way there,
	// but fails (for us) on named imports, since it uses the full package path.
	switch typeInfo := n.typeInfo.(type) {
	case *types.Named:
		// register that we need an import for this type,
		// so we can get the appropriate alias to use.
		typeName := typeInfo.Obj()
		otherPkg := typeName.Pkg()
		if otherPkg == basePkg.Types {
			// local import
			return typeName.Name()
		}
		alias := imports.NeedImport(loader.NonVendorPath(otherPkg.Path()))
		return alias + "." + typeName.Name()
	case *types.Basic:
		return typeInfo.String()
	case *types.Pointer:
		return "*" + (&namingInfo{typeInfo: typeInfo.Elem()}).Syntax(basePkg, imports)
	case *types.Slice:
		return "[]" + (&namingInfo{typeInfo: typeInfo.Elem()}).Syntax(basePkg, imports)
	case *types.Map:
		return fmt.Sprintf(
			"map[%s]%s",
			(&namingInfo{typeInfo: typeInfo.Key()}).Syntax(basePkg, imports),
			(&namingInfo{typeInfo: typeInfo.Elem()}).Syntax(basePkg, imports))
	default:
		basePkg.AddError(fmt.Errorf("name requested for invalid type: %s", typeInfo))
		return typeInfo.String()
	}
}
