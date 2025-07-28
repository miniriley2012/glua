package main

import (
	"flag"
	"fmt"
	"go/format"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/rangetable"
)

type Range struct {
	Low  uint32
	High uint32
}

func main() {
	var inputs flags
	var pkg string
	var output string
	flag.Var(&inputs, "i", "input path")
	flag.StringVar(&pkg, "p", "", "package name")
	flag.StringVar(&output, "o", "", "output path")
	flag.Parse()

	props := flag.Args()

	var contents string
	for _, input := range inputs {
		data, err := os.ReadFile(input)
		if err != nil {
			panic(err)
		}
		contents += string(data)
	}

	lines := map[string][]Range{}

	for line := range strings.Lines(contents) {
		for _, prop := range props {
			if strings.Contains(line, "; "+prop) {
				field := strings.TrimSpace(strings.SplitN(line, ";", 2)[0])
				low, high, isRange := strings.Cut(field, "..")
				if isRange {
					low, err := strconv.ParseUint(low, 16, 32)
					if err != nil {
						panic(err)
					}
					high, err := strconv.ParseUint(high, 16, 32)
					if err != nil {
						panic(err)
					}
					lines[prop] = append(lines[prop], Range{uint32(low), uint32(high)})
				} else {
					n, err := strconv.ParseUint(low, 16, 32)
					if err != nil {
						panic(err)
					}
					lines[prop] = append(lines[prop], Range{uint32(n), uint32(n)})
				}
			}
		}
	}

	const t = 0x10000

	tables := map[string]*unicode.RangeTable{}
	for prop, ranges := range lines {
		tbl := &unicode.RangeTable{}
		for _, rng := range ranges {
			if rng.Low >= t {
				tbl.R32 = append(tbl.R32, unicode.Range32{
					Lo:     rng.Low,
					Hi:     rng.High,
					Stride: 1,
				})
			} else if rng.High < t {
				tbl.R16 = append(tbl.R16, unicode.Range16{
					Lo:     uint16(rng.Low),
					Hi:     uint16(rng.High),
					Stride: 1,
				})
			} else {
				tbl.R16 = append(tbl.R16, unicode.Range16{
					Lo:     uint16(rng.Low),
					Hi:     uint16(0xFFFF),
					Stride: 1,
				})
				tbl.R32 = append(tbl.R32, unicode.Range32{
					Lo:     t,
					Hi:     rng.High,
					Stride: 1,
				})
			}
		}
		tbl = rangetable.Merge(tbl)
		tables[prop] = tbl
	}

	var builder strings.Builder

	builder.WriteString("package ")
	builder.WriteString(pkg)
	builder.WriteString("\n\nimport \"unicode\"\n\n")

	for _, k := range slices.Sorted(maps.Keys(tables)) {
		builder.WriteString(fmt.Sprintf("var %s = %#v\n", k, tables[k]))
	}

	data, err := format.Source([]byte(builder.String()))
	if err != nil {
		panic(err)
	}

	if output == "" {
		fmt.Printf("%s", data)
	} else {
		err = os.WriteFile(output, data, 0o666)
		if err != nil {
			panic(err)
		}
	}
}

type flags []string

func (f *flags) String() string {
	if f != nil {
		return strings.Join(*f, ",")
	}
	return ""
}

func (f *flags) Set(s string) error {
	*f = append(*f, s)
	return nil
}
