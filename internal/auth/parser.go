package auth

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

//go:embed css.wasm
var cssWasm []byte

type tokenParser struct {
	rt  wazero.Runtime
	mod api.Module
	cdx api.Function
	rdx api.Function
	bdx api.Function
	ndx api.Function
	mdx api.Function
}

func newTokenParser() (*tokenParser, error) {
	ctx := context.Background()
	rt := wazero.NewRuntime(ctx)

	compiled, err := rt.CompileModule(ctx, cssWasm)
	if err != nil {
		_ = rt.Close(ctx)
		return nil, fmt.Errorf("compile wasm: %w", err)
	}
	mod, err := rt.InstantiateModule(ctx, compiled, wazero.NewModuleConfig())
	if err != nil {
		_ = rt.Close(ctx)
		return nil, fmt.Errorf("instantiate wasm: %w", err)
	}

	getExport := func(name string) (api.Function, error) {
		f := mod.ExportedFunction(name)
		if f == nil {
			return nil, fmt.Errorf("export %q not found", name)
		}
		return f, nil
	}

	cdx, err := getExport("cdx")
	if err != nil {
		_ = rt.Close(ctx)
		return nil, err
	}
	rdx, err := getExport("rdx")
	if err != nil {
		_ = rt.Close(ctx)
		return nil, err
	}
	bdx, err := getExport("bdx")
	if err != nil {
		_ = rt.Close(ctx)
		return nil, err
	}
	ndx, err := getExport("ndx")
	if err != nil {
		_ = rt.Close(ctx)
		return nil, err
	}
	mdx, err := getExport("mdx")
	if err != nil {
		_ = rt.Close(ctx)
		return nil, err
	}

	return &tokenParser{
		rt:  rt,
		mod: mod,
		cdx: cdx, rdx: rdx, bdx: bdx, ndx: ndx, mdx: mdx,
	}, nil
}

func (p *tokenParser) close(ctx context.Context) error {
	return p.rt.Close(ctx)
}

type tokenIndices struct {
	access  []int
	refresh []int
}

func (p *tokenParser) indicesFromSalts(s [5]int) (tokenIndices, error) {
	ctx := context.Background()

	call5 := func(f api.Function, a, b, c, d, e int) (int, error) {
		res, err := f.Call(ctx,
			uint64(uint32(a)), uint64(uint32(b)),
			uint64(uint32(c)), uint64(uint32(d)), uint64(uint32(e)),
		)
		if err != nil {
			return 0, err
		}
		return int(int32(res[0])), nil
	}

	s1, s2, s3, s4, s5 := s[0], s[1], s[2], s[3], s[4]

	n, err := call5(p.cdx, s1, s2, s3, s4, s5)
	if err != nil {
		return tokenIndices{}, err
	}
	l, err := call5(p.rdx, s1, s2, s4, s3, s5)
	if err != nil {
		return tokenIndices{}, err
	}
	o, err := call5(p.bdx, s1, s2, s4, s3, s5)
	if err != nil {
		return tokenIndices{}, err
	}
	pIdx, err := call5(p.ndx, s1, s2, s4, s3, s5)
	if err != nil {
		return tokenIndices{}, err
	}
	q, err := call5(p.mdx, s1, s2, s4, s3, s5)
	if err != nil {
		return tokenIndices{}, err
	}

	a, err := call5(p.cdx, s2, s1, s3, s5, s4)
	if err != nil {
		return tokenIndices{}, err
	}
	b, err := call5(p.rdx, s2, s1, s3, s4, s5)
	if err != nil {
		return tokenIndices{}, err
	}
	c, err := call5(p.bdx, s2, s1, s4, s3, s5)
	if err != nil {
		return tokenIndices{}, err
	}
	d, err := call5(p.ndx, s2, s1, s4, s3, s5)
	if err != nil {
		return tokenIndices{}, err
	}
	e, err := call5(p.mdx, s2, s1, s4, s3, s5)
	if err != nil {
		return tokenIndices{}, err
	}

	return tokenIndices{
		access:  []int{n, l, o, pIdx, q},
		refresh: []int{a, b, c, d, e},
	}, nil
}
