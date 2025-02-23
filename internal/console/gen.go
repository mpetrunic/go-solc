//go:build ignore

package main

import (
	"embed"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mpetrunic/go-solc/internal/console"
)

var (
	//go:embed *.tmpl
	templates embed.FS

	logArgs    = []string{"string", "uint", "address", "bool"}
	logArgsSet = setOf(logArgs)
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// parse templates
	tmpl, err := template.New("").ParseFS(templates, "*")
	if err != nil {
		return err
	}

	signatures := genSignatures()

	args := make([]Args, len(signatures))
	for i, a := range signatures {
		sig := "log(" + strings.Join(a, ",") + ")"
		sel := ([4]byte)(crypto.Keccak256([]byte(sig))[:4])

		args[i] = Args{
			Sig:  sig,
			Sel:  sel,
			Args: a,
		}
	}

	model := &model{
		Addr: console.Addr,
		Args: args,
	}

	if err := gen("args.go", tmpl, model); err != nil {
		return err
	}
	if err := gen("console.sol", tmpl, model); err != nil {
		return err
	}
	return nil
}

func genSignatures() [][]string {
	signatures := [][]string{
		{},
		{"string"},
		{"uint"},
		{"int"},
		{"bool"},
		{"address"},
		{"bytes"},
		{"bytes1"},
		{"bytes2"},
		{"bytes3"},
		{"bytes4"},
		{"bytes5"},
		{"bytes6"},
		{"bytes7"},
		{"bytes8"},
		{"bytes9"},
		{"bytes10"},
		{"bytes11"},
		{"bytes12"},
		{"bytes13"},
		{"bytes14"},
		{"bytes15"},
		{"bytes16"},
		{"bytes17"},
		{"bytes18"},
		{"bytes19"},
		{"bytes20"},
		{"bytes21"},
		{"bytes22"},
		{"bytes23"},
		{"bytes24"},
		{"bytes25"},
		{"bytes26"},
		{"bytes27"},
		{"bytes28"},
		{"bytes29"},
		{"bytes30"},
		{"bytes31"},
		{"bytes32"},
	}

	for _, arg0 := range logArgs {
		for _, arg1 := range logArgs {
			signatures = append(signatures, []string{arg0, arg1})
		}
	}

	for _, arg0 := range logArgs {
		for _, arg1 := range logArgs {
			for _, arg2 := range logArgs {
				signatures = append(signatures, []string{arg0, arg1, arg2})
			}
		}
	}

	for _, arg0 := range logArgs {
		for _, arg1 := range logArgs {
			for _, arg2 := range logArgs {
				for _, arg3 := range logArgs {
					signatures = append(signatures, []string{arg0, arg1, arg2, arg3})
				}
			}
		}
	}

	return signatures
}

func gen(name string, tmpl *template.Template, model *model) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintln(f, "// Code generated by go generate; DO NOT EDIT.")

	return tmpl.ExecuteTemplate(f, name+".tmpl", model)
}

type model struct {
	Addr common.Address
	Args []Args
}

type Args struct {
	Sig  string  // function signature
	Sel  [4]byte // 4 bytes selector
	Args []string
}

func (a Args) SelString() string {
	return fmt.Sprintf("{0x%02x, 0x%02x, 0x%02x, 0x%02x}", a.Sel[0], a.Sel[1], a.Sel[2], a.Sel[3])
}

func (a Args) SignatureArgs() string {
	args := make([]string, len(a.Args))
	for i, arg := range a.Args {
		switch arg {
		case "string", "bytes":
			args[i] = fmt.Sprintf("%s memory p%d", arg, i)
		default:
			args[i] = fmt.Sprintf("%s p%d", arg, i)
		}
	}
	return strings.Join(args, ", ")
}

func (a Args) LogTypeSignature() string {
	return fmt.Sprintf("log%s(%s)", strings.Title(a.Args[0]), a.SignatureArgs())
}

func (a Args) LogSignature() string {
	return fmt.Sprintf("log(%s)", a.SignatureArgs())
}

func (a Args) IsLogType() bool {
	return len(a.Args) == 1
}

func (a Args) IsLog() bool {
	if len(a.Args) == 0 || len(a.Args) >= 2 {
		return true
	}
	_, ok := logArgsSet[a.Args[0]]
	return ok
}

func (a Args) Params() string {
	params := make([]string, len(a.Args))
	for i := range a.Args {
		params[i] = fmt.Sprintf("p%d", i)
	}
	return strings.Join(params, ", ")
}

func (a Args) ArgTypes() string {
	args := make([]string, len(a.Args))
	for i, arg := range a.Args {
		args[i] = fmt.Sprintf("arg%s", strings.Title(arg))
	}
	return strings.Join(args, ", ")
}

func setOf[T comparable](s []T) map[T]struct{} {
	m := make(map[T]struct{}, len(s))
	for _, v := range s {
		m[v] = struct{}{}
	}
	return m
}
