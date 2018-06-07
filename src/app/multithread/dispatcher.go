package multithread

var (
	defaultMaxWorker = 100
)

type Job interface {
	execute()
}
type worker struct {
	jobPool chan Job
	quit    chan bool
}

func newWorker(jobPool chan Job) worker {
	return worker{
		jobPool: jobPool,
		quit:    make(chan bool)}
}

func (w *worker) start() {
	go func() {
		for {
			select {
			case job, ok := <-w.jobPool:
				if !ok {
					return
				}
				job.execute()
			case <-w.quit:
				return
			}
		}
	}()
}

func (w *worker) stop() {
	go func() {
		w.quit <- true
	}()
}

type Dispatcher struct {
	jobPool    chan Job
	workerPool []worker
}

func NewDispatcher(jobPool chan Job, maxWorkers int) *Dispatcher {
	numWorker := maxWorkers
	if numWorker <= 0 {
		numWorker = defaultMaxWorker
	}

	workers := []worker{}
	for i := 0; i < numWorker; i++ {
		worker := newWorker(jobPool)
		workers = append(workers, worker)
	}
	return &Dispatcher{jobPool: jobPool, workerPool: workers}
}

func (d *Dispatcher) Run() {
	for _, w := range d.workerPool {
		w.start()
	}
}

func (d *Dispatcher) Stop() {
	for _, w := range d.workerPool {
		w.stop()
	}
}
