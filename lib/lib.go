package lib

import (
	"cmp"
	"os"
	"slices"

	"github.com/dcaiafa/bag3l/internal/export"
	"github.com/dcaiafa/bag3l/internal/vm"
	libbool "github.com/dcaiafa/bag3l/lib/bool"
	"github.com/dcaiafa/bag3l/lib/crypto"
	"github.com/dcaiafa/bag3l/lib/encoding/base64"
	libcsv "github.com/dcaiafa/bag3l/lib/encoding/csv"
	libhex "github.com/dcaiafa/bag3l/lib/encoding/hex"
	"github.com/dcaiafa/bag3l/lib/encoding/json"
	"github.com/dcaiafa/bag3l/lib/file"
	"github.com/dcaiafa/bag3l/lib/float"
	"github.com/dcaiafa/bag3l/lib/global"
	libint "github.com/dcaiafa/bag3l/lib/int"
	"github.com/dcaiafa/bag3l/lib/io"
	"github.com/dcaiafa/bag3l/lib/list"
	"github.com/dcaiafa/bag3l/lib/maps"
	ospkg "github.com/dcaiafa/bag3l/lib/os"
	"github.com/dcaiafa/bag3l/lib/path/filepath"
	"github.com/dcaiafa/bag3l/lib/str"
	libtime "github.com/dcaiafa/bag3l/lib/time"
)

var GlobalPackage = export.Exports{
	{N: "$close_iter", T: export.Func, F: close_iter},
	{N: "$concat", T: export.Func, F: concat},
	{N: "$exec", T: export.Func, F: execExec},
	{N: "$format", T: export.Func, F: format},
	{N: "$home", T: export.Func, F: internalHomeDir},
	// REVISIT: avg uses custom avgAccum type and runtime-based overload resolution
	{N: "avg", T: export.Func, F: avg},
	// REVISIT: count uses custom countAccum type and runtime-based overload resolution.
	{N: "count", T: export.Func, F: count},
	// REVISIT: first uses iter.IterNRet() for dynamic return count; cannot be statically typed in stubgen.
	{N: "first", T: export.Func, F: first},
	{N: "flatten", T: export.Func, F: flatten},
	{N: "get", T: export.Func, F: get},
	{N: "group", T: export.Func, F: group},
	{N: "has", T: export.Func, F: has},
	{N: "hist", T: export.Func, F: hist},
	{N: "in", T: export.Func, F: in},
	{N: "index", T: export.Func, F: index},
	{N: "iterate", T: export.Func, F: iterate},
	{N: "len", T: export.Func, F: Len},
	{N: "map", T: export.Func, F: imap},
	{N: "map_reduce", T: export.Func, F: mapReduce},
	{N: "max", T: export.Func, F: max},
	{N: "min", T: export.Func, F: min},
	{N: "mod", T: export.Func, F: mod},
	{N: "narg", T: export.Func, F: narg},
	{N: "next", T: export.Func, F: next},
	{N: "print", T: export.Func, F: print},
	{N: "print_table", T: export.Func, F: printTable},
	{N: "probe", T: export.Func, F: probe},
	{N: "prompt", T: export.Func, F: prompt},
	{N: "range", T: export.Func, F: range_},
	{N: "read", T: export.Func, F: read},
	{N: "reduce", T: export.Func, F: reduce},
	{N: "regex", T: export.Func, F: regex},
	{N: "sha1", T: export.Func, F: sha1},
	{N: "shuffle", T: export.Func, F: shuffle},
	{N: "skip", T: export.Func, F: skip},
	{N: "skip_until", T: export.Func, F: skipUntil},
	{N: "skip_while", T: export.Func, F: skipWhile},
	{N: "sleep", T: export.Func, F: sleep},
	{N: "sort", T: export.Func, F: sort},
	{N: "stream", T: export.Func, F: stream},
	{N: "sum", T: export.Func, F: sum},
	{N: "take", T: export.Func, F: take},
	{N: "take_until", T: export.Func, F: takeUntil},
	{N: "take_while", T: export.Func, F: takeWhile},
	{N: "type", T: export.Func, F: typep},
	{N: "unique", T: export.Func, F: unique},
}

var BufPackage = export.Exports{
	{N: "cap", T: export.Func, F: bufCap},
	{N: "len", T: export.Func, F: bufLen},
	{N: "new", T: export.Func, F: bufNew},
	{N: "read", T: export.Func, F: bufRead},
	{N: "read_byte", T: export.Func, F: bufReadByte},
	{N: "read_from", T: export.Func, F: bufReadFrom},
	{N: "read_rune", T: export.Func, F: bufReadRune},
	{N: "unread_byte", T: export.Func, F: bufUnreadByte},
}

var CoPackage = export.Exports{
	{N: "run_with_timeout", T: export.Func, F: runWithTimeout},
	{N: "start", T: export.Func, F: start},
}

var ExecPackage = export.Exports{
	{N: "exec", T: export.Func, F: execExec},
	{N: "with_stderr", T: export.Func, F: execWithStderr},
}

var MathPackage = export.Exports{
	{N: "trunc", T: export.Func, F: mathTrunc},
}

type BuiltinRegistry interface {
	RegisterBuiltins(pkgName string, exports export.Exports)
}

func RegisterAll(registry BuiltinRegistry) {
	globalPackage := slices.Concat(GlobalPackage, global.Exports)
	slices.SortFunc(globalPackage, func(a, b export.Export) int {
		return cmp.Compare(a.N, b.N)
	})

	registry.RegisterBuiltins("$global", globalPackage)
	registry.RegisterBuiltins("bool", libbool.Exports)
	registry.RegisterBuiltins("buf", BufPackage)
	registry.RegisterBuiltins("co", CoPackage)
	registry.RegisterBuiltins("crypto", crypto.Exports)
	registry.RegisterBuiltins("encoding/base64", base64.Exports)
	registry.RegisterBuiltins("encoding/csv", libcsv.Exports)
	registry.RegisterBuiltins("encoding/hex", libhex.Exports)
	registry.RegisterBuiltins("encoding/json", json.Exports)
	registry.RegisterBuiltins("exec", ExecPackage)
	registry.RegisterBuiltins("file", file.Exports)
	registry.RegisterBuiltins("float", float.Exports)
	registry.RegisterBuiltins("int", libint.Exports)
	registry.RegisterBuiltins("io", io.Exports)
	registry.RegisterBuiltins("list", list.Exports)
	registry.RegisterBuiltins("maps", maps.Exports)
	registry.RegisterBuiltins("math", MathPackage)
	registry.RegisterBuiltins("os", ospkg.Exports)
	registry.RegisterBuiltins("path/filepath", filepath.Exports)
	registry.RegisterBuiltins("str", str.Exports)
	registry.RegisterBuiltins("time", libtime.Exports)
}

func internalHomeDir(m *vm.VM, args []vm.Value, nret int) ([]vm.Value, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return []vm.Value{vm.NewString(dir)}, nil
}
