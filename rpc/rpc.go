package rpc

import (
	"context"
	"fmt"

	"github.com/letgoapp/go-bloomfilter/rotate"
)

type Config struct {
	rotate.Config
	Port int
}

type BloomfilterRPC int

type Bloomfilter struct {
	BloomfilterRPC
}

func New(ctx context.Context, cfg Config) *Bloomfilter {
	if bf != nil {
		bf.Close()
	}

	bf = rotate.New(ctx, cfg.Config)

	return new(Bloomfilter)
}

type AddInput struct {
	Elems [][]byte
}

type AddOutput struct {
	Count int
}

func (b *BloomfilterRPC) Add(in AddInput, out *AddOutput) error {
	fmt.Println("add:", in.Elems)
	defer func() { fmt.Println("added elements:", out.Count) }()

	if bf == nil {
		out.Count = 0
		return ErrNoBloomfilterInitialized
	}

	k := 0
	for _, elem := range in.Elems {
		bf.Add(elem)
		k++
	}
	out.Count = k

	return nil
}

type CheckInput struct {
	Elems [][]byte
}

type CheckOutput struct {
	Checks []bool
}

func (b *BloomfilterRPC) Check(in CheckInput, out *CheckOutput) error {
	fmt.Println("check:", in.Elems)
	defer func() { fmt.Println("checked elements:", out.Checks) }()

	checkRes := make([]bool, len(in.Elems))

	if bf == nil {
		out.Checks = checkRes
		return ErrNoBloomfilterInitialized
	}

	for i, elem := range in.Elems {
		checkRes[i] = bf.Check(elem)
	}
	out.Checks = checkRes

	return nil
}

type UnionInput struct {
	BF *rotate.Bloomfilter
}

type UnionOutput struct {
	Capacity float64
}

func (b *BloomfilterRPC) Union(in UnionInput, out *UnionOutput) error {
	fmt.Println("union:", in.BF)
	defer func() { fmt.Println("union resulting capacity:", out.Capacity) }()

	if bf == nil {
		out.Capacity = 0
		return ErrNoBloomfilterInitialized
	}

	var err error
	out.Capacity, err = bf.Union(in.BF)

	return err
}

func (b *Bloomfilter) Close() {
	if bf != nil {
		bf.Close()
	}
}

func (b *Bloomfilter) Bloomfilter() *rotate.Bloomfilter {
	return bf
}

var (
	ErrNoBloomfilterInitialized = fmt.Errorf("Bloomfilter not initialized")
	bf                          *rotate.Bloomfilter
)
