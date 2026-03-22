package fiber

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// gate provides deterministic control over a single blocking operation.
type gate struct {
	entered chan struct{}
	release chan struct{}
}

func newGate() gate {
	return gate{
		entered: make(chan struct{}),
		release: make(chan struct{}),
	}
}

func (g gate) blockFn(ctx context.Context) {
	close(g.entered)
	select {
	case <-g.release:
	case <-ctx.Done():
	}
}

func TestSimpleFiber(t *testing.T) {
	var ran bool
	s := NewScheduler()

	main := s.NewFiber(func() {
		ran = true
	})

	s.Run(main)
	require.True(t, ran)
}

func TestSwitchToNew(t *testing.T) {
	var res strings.Builder
	s := NewScheduler()

	main := s.NewFiber(func() {
		res.WriteString("main:1\n")
		child := s.NewFiber(func() {
			res.WriteString("child:1\n")
		})
		s.SwitchToNew(child)
		res.WriteString("main:2\n")
	})

	s.Run(main)

	require.Equal(t, "main:1\nchild:1\nmain:2", strings.TrimSpace(res.String()))
}

func TestNestedSwitchToNew(t *testing.T) {
	var res strings.Builder
	s := NewScheduler()

	main := s.NewFiber(func() {
		res.WriteString("a\n")
		b := s.NewFiber(func() {
			res.WriteString("b\n")
			c := s.NewFiber(func() {
				res.WriteString("c\n")
			})
			s.SwitchToNew(c)
			res.WriteString("b:end\n")
		})
		s.SwitchToNew(b)
		res.WriteString("a:end\n")
	})

	s.Run(main)

	require.Equal(t, "a\nb\nc\na:end\nb:end", strings.TrimSpace(res.String()))
}

func TestBlock(t *testing.T) {
	var res strings.Builder
	s := NewScheduler()

	g := newGate()

	main := s.NewFiber(func() {
		res.WriteString("before\n")
		s.Block(context.Background(), g.blockFn)
		res.WriteString("after\n")
	})

	done := make(chan struct{})
	go func() {
		s.Run(main)
		close(done)
	}()

	<-g.entered
	close(g.release)
	<-done

	require.Equal(t, "before\nafter", strings.TrimSpace(res.String()))
}

func TestBlockYieldsToReady(t *testing.T) {
	var res strings.Builder
	s := NewScheduler()

	g := newGate()

	main := s.NewFiber(func() {
		child := s.NewFiber(func() {
			res.WriteString("child:1\n")
			s.Block(context.Background(), g.blockFn)
			res.WriteString("child:2\n")
		})
		s.SwitchToNew(child)
		res.WriteString("main:1\n")
	})

	done := make(chan struct{})
	go func() {
		s.Run(main)
		close(done)
	}()

	<-g.entered
	close(g.release)
	<-done

	require.Equal(t, "child:1\nmain:1\nchild:2", strings.TrimSpace(res.String()))
}

func TestMultipleFibersInterleaving(t *testing.T) {
	var res strings.Builder
	s := NewScheduler()

	gcGates := [3]gate{newGate(), newGate(), newGate()}
	c2Gates := [3]gate{newGate(), newGate(), newGate()}
	gcDone := make(chan struct{})
	c2Done := make(chan struct{})

	main := s.NewFiber(func() {
		res.WriteString("main started\n")
		child1 := s.NewFiber(func() {
			res.WriteString("child1 started\n")
			grandChild1 := s.NewFiber(func() {
				for i := range 3 {
					fmt.Fprintf(&res, "grandChild1 block %v\n", i)
					s.Block(context.Background(), gcGates[i].blockFn)
				}
				res.WriteString("grandChild1 finished\n")
				close(gcDone)
			})
			s.SwitchToNew(grandChild1)
			res.WriteString("child1 ended\n")
		})
		s.SwitchToNew(child1)
		child2 := s.NewFiber(func() {
			for i := range 3 {
				fmt.Fprintf(&res, "child2 block %v\n", i)
				s.Block(context.Background(), c2Gates[i].blockFn)
			}
			res.WriteString("child2 finished\n")
			close(c2Done)
		})
		s.SwitchToNew(child2)
		res.WriteString("main ended\n")
	})

	done := make(chan struct{})
	go func() {
		s.Run(main)
		close(done)
	}()

	// Both initial blocks are entered.
	<-gcGates[0].entered
	<-c2Gates[0].entered

	// Release child2 first, then grandChild1.
	close(c2Gates[0].release)
	<-c2Gates[1].entered

	close(gcGates[0].release)
	<-gcGates[1].entered

	close(c2Gates[1].release)
	<-c2Gates[2].entered

	close(c2Gates[2].release)
	<-c2Done

	close(gcGates[1].release)
	<-gcGates[2].entered

	close(gcGates[2].release)
	<-gcDone

	<-done

	expected := `main started
child1 started
grandChild1 block 0
child2 block 0
child1 ended
main ended
child2 block 1
grandChild1 block 1
child2 block 2
child2 finished
grandChild1 block 2
grandChild1 finished`
	require.Equal(t, expected, strings.TrimSpace(res.String()))
}

func TestCancelBlocked(t *testing.T) {
	s := NewScheduler()

	g := newGate()
	var canceled bool

	main := s.NewFiber(func() {
		s.Block(context.Background(), func(ctx context.Context) {
			close(g.entered)
			<-ctx.Done()
			canceled = true
		})
	})

	done := make(chan struct{})
	go func() {
		s.Run(main)
		close(done)
	}()

	<-g.entered
	s.CancelBlocked()
	<-done

	require.True(t, canceled)
}

func TestForEachFiber(t *testing.T) {
	s := NewScheduler()
	g := newGate()

	var count int

	main := s.NewFiber(func() {
		child := s.NewFiber(func() {
			s.Block(context.Background(), g.blockFn)
		})
		s.SwitchToNew(child)
		// child is blocked, main is active.
		s.ForEachFiber(func(f *Fiber) {
			count++
		})
		close(g.release)
	})

	s.Run(main)
	require.Equal(t, 2, count)
}
