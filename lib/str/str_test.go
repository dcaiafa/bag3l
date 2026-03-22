package str_test

import (
	"testing"

	"github.com/dcaiafa/bag3l/internal/btesting"
)

func TestStrFind(t *testing.T) {
	btesting.RunSubO(t, "found", `"banana" | str.find("na") | print `, `2`)
	btesting.RunSubO(t, "not_found", `"banana" | str.find("no") | print `, `<nil>`)
}

func TestStrFindLast(t *testing.T) {
	btesting.RunSubO(t, "found", `"banana" | str.find_last("na") | print `, `4`)
	btesting.RunSubO(t, "not_found", `"banana" | str.find_last("no") | print `, `<nil>`)
}

func TestStrMatch(t *testing.T) {
	btesting.RunSubO(t, "match",
		"var re = r`(\\d+)/.*$`"+`
    "foo/123/bar" | str.match(re) | print`, `
[123/bar 123]
  `)
	btesting.RunSubO(t, "no_match",
		"var re = r`(\\d+)/.*$`"+`
    "foo/12a/bar" | str.match(re) | print`, `
<nil>
  `)
}

func TestStrMatchAll(t *testing.T) {
	btesting.RunSubO(t, "match",
		"var re = r`foo(.?)`"+`
    "seafood fool" | str.match_all(re) | print`, `
[[food d] [fool l]]
  `)
	btesting.RunSubO(t, "no_match",
		"var re = r`foo(.?)`"+`
    "seafeod fiol" | str.match_all(re) | print`, `
<nil>
  `)
}

func TestStrReplace(t *testing.T) {
	btesting.RunSubO(t, "string_all", `"aaa" | str.replace("a", "b") | print`, `bbb`)
	btesting.RunSubO(t, "string_limited", `"aaa" | str.replace("a", "b", 2) | print`, `bba`)
	btesting.RunSubO(t, "regex_string",
		"var re = r`\\d+`"+`
    "abc123def456" | str.replace(re, "X") | print`, `abcXdefX`)
	btesting.RunSubO(t, "regex_func",
		"var re = r`\\d+`"+`
    "abc123def456" | str.replace(re, &m -> f"[{m}]") | print`, `abc[123]def[456]`)
}

func TestStrSplit(t *testing.T) {
	btesting.RunSubO(t, "string", `"a,b,c" | str.split(",") | print`, `[a b c]`)
	btesting.RunSubO(t, "string_limited", `"a,b,c,d" | str.split(",", 3) | print`, `[a b c,d]`)
	btesting.RunSubO(t, "regex",
		"var re = r`\\s+`"+`
    "one  two   three" | str.split(re) | print`, `[one two three]`)
	btesting.RunSubO(t, "regex_limited",
		"var re = r`\\s+`"+`
    "one  two   three" | str.split(re, 2) | print`, `[one two   three]`)
}

func TestStrTrimSpace(t *testing.T) {
	btesting.RunSubO(t, "basic", `"  hello  " | str.trim_space | print`, `hello`)
	btesting.RunSubO(t, "tabs", `"\thello\n" | str.trim_space | print`, `hello`)
	btesting.RunSubO(t, "no_space", `"hello" | str.trim_space | print`, `hello`)
}

func TestStrTrimPrefix(t *testing.T) {
	btesting.RunSubO(t, "match", `"hello world" | str.trim_prefix("hello ") | print`, `world`)
	btesting.RunSubO(t, "no_match", `"hello world" | str.trim_prefix("foo") | print`, `hello world`)
}

func TestStrTrimSuffix(t *testing.T) {
	btesting.RunSubO(t, "match", `"hello.txt" | str.trim_suffix(".txt") | print`, `hello`)
	btesting.RunSubO(t, "no_match", `"hello.txt" | str.trim_suffix(".go") | print`, `hello.txt`)
}

func TestStrToUpper(t *testing.T) {
	btesting.RunSubO(t, "basic", `"hello" | str.to_upper | print`, `HELLO`)
	btesting.RunSubO(t, "mixed", `"Hello World" | str.to_upper | print`, `HELLO WORLD`)
}

func TestStrToLower(t *testing.T) {
	btesting.RunSubO(t, "basic", `"HELLO" | str.to_lower | print`, `hello`)
	btesting.RunSubO(t, "mixed", `"Hello World" | str.to_lower | print`, `hello world`)
}

func TestStrHasPrefix(t *testing.T) {
	btesting.RunSubO(t, "true", `"hello world" | str.has_prefix("hello") | print`, `true`)
	btesting.RunSubO(t, "false", `"hello world" | str.has_prefix("world") | print`, `false`)
}

func TestStrHasSuffix(t *testing.T) {
	btesting.RunSubO(t, "true", `"hello world" | str.has_suffix("world") | print`, `true`)
	btesting.RunSubO(t, "false", `"hello world" | str.has_suffix("hello") | print`, `false`)
}

func TestStrContains(t *testing.T) {
	btesting.RunSubO(t, "true", `"hello world" | str.contains("lo wo") | print`, `true`)
	btesting.RunSubO(t, "false", `"hello world" | str.contains("xyz") | print`, `false`)
	btesting.RunSubO(t, "empty", `"hello" | str.contains("") | print`, `true`)
}

func TestStrContainsAny(t *testing.T) {
	btesting.RunSubO(t, "true", `"hello" | str.contains_any("aeiou") | print`, `true`)
	btesting.RunSubO(t, "false", `"hello" | str.contains_any("xyz") | print`, `false`)
}

func TestStrCount(t *testing.T) {
	btesting.RunSubO(t, "multiple", `"banana" | str.count("an") | print`, `2`)
	btesting.RunSubO(t, "none", `"banana" | str.count("xyz") | print`, `0`)
}

func TestStrEqualFold(t *testing.T) {
	btesting.RunSubO(t, "true", `str.equal_fold("Hello", "hello") | print`, `true`)
	btesting.RunSubO(t, "false", `str.equal_fold("Hello", "world") | print`, `false`)
}

func TestStrCut(t *testing.T) {
	btesting.RunSubO(t, "found", `
    var before, after, found = str.cut("hello=world", "=")
    print(before, after, found)`, `hello world true`)
	btesting.RunSubO(t, "not_found", `
    var before, after, found = str.cut("hello", "=")
    print(before, after, found)`, `hello  false`)
}

func TestStrTrim(t *testing.T) {
	btesting.RunSubO(t, "basic", `"##hello##" | str.trim("#") | print`, `hello`)
	btesting.RunSubO(t, "multi_chars", `"¡¡¡Hello, Gophers!!!" | str.trim("!¡") | print`, `Hello, Gophers`)
}

func TestStrTrimLeft(t *testing.T) {
	btesting.RunSubO(t, "basic", `"##hello##" | str.trim_left("#") | print`, `hello##`)
}

func TestStrTrimRight(t *testing.T) {
	btesting.RunSubO(t, "basic", `"##hello##" | str.trim_right("#") | print`, `##hello`)
}

func TestStrToTitle(t *testing.T) {
	btesting.RunSubO(t, "basic", `"hello world" | str.to_title | print`, `HELLO WORLD`)
}

func TestStrFields(t *testing.T) {
	btesting.RunSubO(t, "basic", `"  one  two   three  " | str.fields | print`, `[one two three]`)
	btesting.RunSubO(t, "single", `"hello" | str.fields | print`, `[hello]`)
}

func TestStrRepeat(t *testing.T) {
	btesting.RunSubO(t, "basic", `"ab" | str.repeat(3) | print`, `ababab`)
	btesting.RunSubO(t, "zero", `"ab" | str.repeat(0) | print`, ``)
}

func TestStrInto(t *testing.T) {
	btesting.RunSubO(t, "int", `str.into(42) | print`, `42`)
	btesting.RunSubO(t, "bool", `str.into(true) | print`, `true`)
	btesting.RunSubO(t, "float", `str.into(3.14) | print`, `3.14`)
	btesting.RunSubO(t, "string", `str.into("hello") | print`, `hello`)
}
