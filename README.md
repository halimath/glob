# glob

Advanced filesystem glob for golang

![CI Status][ci-img-url] 
[![Go Report Card][go-report-card-img-url]][go-report-card-url] 
[![Package Doc][package-doc-img-url]][package-doc-url] 
[![Releases][release-img-url]][release-url]

[ci-img-url]: https://github.com/halimath/glob/workflows/CI/badge.svg
[go-report-card-img-url]: https://goreportcard.com/badge/github.com/halimath/glob
[go-report-card-url]: https://goreportcard.com/report/github.com/halimath/glob
[package-doc-img-url]: https://img.shields.io/badge/GoDoc-Reference-blue.svg
[package-doc-url]: https://pkg.go.dev/github.com/halimath/glob
[release-img-url]: https://img.shields.io/github/v/release/halimath/glob.svg
[release-url]: https://github.com/halimath/glob/releases

`glob` provides an advanced file system glob language, a superset of the 
pattern language provided by that of the golang standard lib's `fs` package.

# Installation

`glob` is provided as a go module and requires go >= 1.18.

```shell
go get github.com/halimath/glob@main
```

# Usage

`glob` provides a type `Pattern` which can be created using the `New` function:

```go
pat, err := glob.New("**/*_test.go")
```

A `Pattern` may then be used to search for matches in a `fs.FS`. If you want all
matches, simply use the `GlobFS` method:

```go
files, err := pat.GlobFS(fsys, "")
```

# Pattern language

The pattern language used by `glob` works similar to the 
[pattern format of `.gitignore`](https://git-scm.com/docs/gitignore). It is
completely compatible with the pattern format used by `os.Glob` or `fs.Glob`
and extends it.

The format is specified as the following EBNF:

```ebnf
pattern = term, { '/', term };

term        = '**' | name;
name        = { charSpecial | group | escapedChar | '*' | '?' };
charSpecial = (* any unicode rune except '/', '*', '?', '[' and '\' *);
char        = (* any unicode rune *);
escapedChar = '\\', char;
group       = '[', [ '^' ] { escapedChar | groupChar | range } ']';
groupChar   = (* any unicode rune except '-' and ']' *);
range       = ( groupChar | escapedChar ), '-', (groupChar | escapedChar);
```

The format operators have the following meaning:

* any character (rune) matches the exactly this rune - with the following
  exceptions
* `/` works as a directory separator. It matches directory boundarys of the
  underlying system independently of the separator char used by the OS.
* `?` matches exactly one non-separator char
* `*` matches any number of non-separator chars - including zero
* `\` escapes a character's special meaning allowing `*` and `?` to be used
  as regular characters.
* `**` matches any number of nested directories. If anything is matched it
  always extends until a separator or the end of the name.
* Groups can be defined using the `[` and `]` characters. Inside a group the
  special meaning of the characters mentioned before is disabled but the
  following rules apply
    * any character used as part of the group acts as a choice to pick from
    * if the group's first character is a `^` the whole group is negated
    * a range can be defined using `-` matching any rune between low and high
      inclusive
    * Multiple ranges can be given. Ranges can be combined with choices.
    * The meaning of `-` and `]` can be escacped using `\`

# Performance

`glob` separates pattern parsing and matching. This can create a 
performance benefit when applied repeatedly. When reusing a precompiled pattern
to match filenames `glob` outperforms `filepath.Match` with both simple
and complex patterns. When not reusing the parsed pattern, `filepath` works
much faster (but lacks the additional features).

Test | Execution time `[ns/op]` | Memory usage `[B/op]` | Allocations per op
-- | --: | --: | --:
`filepath` simple pattern                        |   15.5 | 0    | 0
`glob` simple pattern (reuse)               |    3.9 | 0    | 0
`glob` simple pattern (noreuse)             |  495.0 | 1112 | 5
`filepath` complex pattern                       |  226.2 |    0 | 0
`glob` complex pattern (reuse)              |  108.1 |    0 | 0
`glob` complex pattern (noreuse)            | 1103.0 | 2280 | 8
`glob` directory wildcard pattern (reuse)   |  111.7 |    0 | 0
`glob` directory wildcard pattern (noreuse) | 1229.0 | 2280 | 8

# License

Copyright 2022 Alexander Metzner.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
