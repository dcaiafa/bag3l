package fiber

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

type fiberState int

const (
	stateNew fiberState = iota
	stateRunning
	stateReady
	stateBlocking
	stateTerminated
)

// Fiber is a cooperative unit of execution managed by a Scheduler. Each fiber
// runs in its own goroutine but only one fiber is logically active at a time.
// Fibers yield control explicitly via Block or SwitchToNew.
type Fiber struct {
	// Data is an arbitrary user-data slot for attaching per-fiber state.
	Data any

	le    fiberListElem
	id    uint32
	sched *Scheduler
	cv    *sync.Cond
	state fiberState
	fn    func()

	blockCancel context.CancelFunc
}

func (f *Fiber) expectState(state fiberState) {
	if f.state != state {
		panic(fmt.Sprintf("Expected state %v but it was %v", state, f.state))
	}
}

func newFiber(sched *Scheduler, id uint32, fn func()) *Fiber {
	f := &Fiber{
		id:    id,
		sched: sched,
		state: stateNew,
		fn:    fn,
	}
	f.cv = sync.NewCond(&sched.mutex)
	f.le.Self = f
	return f
}

// Scheduler implements cooperative (non-preemptive) scheduling of fibers. At
// most one fiber is active at any time; fibers yield control by calling Block
// or SwitchToNew. Blocking operations run concurrently in their fiber's
// goroutine, allowing I/O parallelism while maintaining single-threaded
// logical execution for non-blocking code.
type Scheduler struct {
	mutex   sync.Mutex
	cv      *sync.Cond
	ready   fiberList
	blocked fiberList
	active  *Fiber
	lastID  uint32
}

// NewScheduler creates a new fiber scheduler.
func NewScheduler() *Scheduler {
	s := &Scheduler{}
	s.cv = sync.NewCond(&s.mutex)
	s.ready.Init()
	s.blocked.Init()
	return s
}

// Run starts the scheduler loop with the given fiber as the initial fiber.
// It blocks until all fibers (active, ready, and blocked) have terminated.
func (s *Scheduler) Run(main *Fiber) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	main.state = stateRunning
	s.active = main
	s.runFiber(main)

	for {
		// Wait while a fiber is currently running or while there is nothing ready.
		for s.active != nil || (s.ready.Empty() && !s.blocked.Empty()) {
			s.cv.Wait()
		}

		if s.ready.Empty() && s.blocked.Empty() {
			break
		}

		f := s.ready.Begin().Self
		f.le.Remove()

		s.active = f
		f.state = stateRunning

		s.mutex.Unlock()
		f.cv.Signal()
		s.mutex.Lock()
	}
}

// Block suspends the active fiber, allowing other fibers to run, while fn
// executes concurrently in this fiber's goroutine. When fn returns, the fiber
// is re-enqueued and waits to be re-scheduled. The context passed to fn is
// derived from ctx and can be cancelled via CancelBlocked.
func (s *Scheduler) Block(ctx context.Context, fn func(ctx context.Context)) {
	s.mutex.Lock()

	me := s.active
	me.expectState(stateRunning)

	ctx, me.blockCancel = context.WithCancel(ctx)
	defer me.blockCancel()

	s.active = nil
	me.state = stateBlocking
	me.le.InsertBefore(s.blocked.End())

	s.mutex.Unlock()
	s.cv.Signal()

	fn(ctx)

	s.mutex.Lock()
	me.state = stateReady
	me.le.Remove()
	me.le.InsertBefore(s.ready.End())
	s.mutex.Unlock()
	s.cv.Signal()

	s.mutex.Lock()
	for me.state != stateRunning {
		me.cv.Wait()
	}
	s.mutex.Unlock()

	me.blockCancel = nil
}

// NewFiber creates a new fiber that will execute fn when started. The fiber is
// not started until it is passed to SwitchToNew or Run.
func (s *Scheduler) NewFiber(fn func()) *Fiber {
	id := atomic.AddUint32(&s.lastID, 1)
	return newFiber(s, id, fn)
}

// SwitchToNew parks the active fiber in the ready queue and immediately starts
// f. The caller resumes when the scheduler re-schedules it (FIFO).
func (s *Scheduler) SwitchToNew(f *Fiber) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	f.expectState(stateNew)

	parked := s.active
	parked.state = stateReady
	parked.le.InsertBefore(s.ready.End())

	s.active = f
	f.state = stateRunning

	s.runFiber(f)

	for parked.state != stateRunning {
		parked.cv.Wait()
	}
}

// ForEachFiber calls h for every fiber in the scheduler: the active fiber (if
// any), all ready fibers, and all blocked fibers.
func (s *Scheduler) ForEachFiber(h func(fiber *Fiber)) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.active != nil {
		h(s.active)
	}
	for le := s.ready.Begin(); le != s.ready.End(); le = le.Next {
		h(le.Self)
	}
	for le := s.blocked.Begin(); le != s.blocked.End(); le = le.Next {
		h(le.Self)
	}
}

// CancelBlocked cancels the context of every currently blocked fiber. This
// causes the context passed to each Block callback to be cancelled, which
// should cause the callback to return promptly.
func (s *Scheduler) CancelBlocked() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for le := s.blocked.Begin(); le != s.blocked.End(); le = le.Next {
		le.Self.blockCancel()
	}
}

func (s *Scheduler) runFiber(f *Fiber) {
	go func() {
		defer func() {
			s.mutex.Lock()
			s.active = nil
			f.state = stateTerminated
			s.mutex.Unlock()
			s.cv.Signal()
		}()

		f.fn()
	}()
}
