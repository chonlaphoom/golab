package throttler

import (
	"context"
	"log"
	"sync"
)

type Task struct {
	fn   func(ctx context.Context) error
	name string
}

type Throttler struct {
	tasks       chan Task
	wg          *sync.WaitGroup
	workerCount int
	result      chan int16 // assuming some result type
	ctx         context.Context
}

func NewThrottler(workerCount int, bufferSize int, ctx context.Context) (*Throttler, error) {
	return &Throttler{
		tasks:       make(chan Task, bufferSize),
		result:      make(chan int16, bufferSize),
		wg:          &sync.WaitGroup{},
		workerCount: workerCount,
		ctx:         ctx,
	}, nil
}

func (t *Throttler) Submit(task Task) {
	select {
	case t.tasks <- task:
	case <-t.ctx.Done():
		log.Println("Context cancelled, cannot submit task:", task.name)
	}
}

func (t *Throttler) Close() {
	close(t.tasks) // No more tasks will be sent
}

func (t *Throttler) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-t.tasks:
			if !ok {
				return // channel closed
			}
			err := task.fn(ctx)
			if err != nil {
				log.Printf("Error: %v", err)
				continue
			}
			select {
			case t.result <- 1:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (t *Throttler) Start() <-chan int16 /* read-only channel */ {
	t.wg.Add(t.workerCount)
	for i := 0; i < t.workerCount; i++ {
		// limited number of goroutines
		go func() {
			defer t.wg.Done()
			t.worker(t.ctx)
		}()
	}

	go func() {
		t.wg.Wait()
		close(t.result)
	}()

	return t.result
}

func Execute() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	throttler, err := NewThrottler(3, 10, ctx)
	if err != nil {
		log.Fatalf("Failed to create throttler: %v", err)
	}

	tasks := []Task{
		{
			fn: func(ctx context.Context) error {
				select {
				case <-ctx.Done():
					return nil
				default:
					log.Println("Executing Task 1")
					return nil
				}
			},
			name: "Task 1",
		},
		{
			fn: func(ctx context.Context) error {
				select {
				case <-ctx.Done():
					return nil
				default:
					log.Println("Executing Task 2")
					return nil
				}
			},
			name: "Task 2",
		},
	}

	results := throttler.Start()

	for _, task := range tasks {
		throttler.Submit(task)
	}

	throttler.Close()

	for res := range results {
		log.Printf("Task completed with result: %d", res)
	}
}
