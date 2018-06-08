package multithread

var (
	defaultMaxWorker = 100
)

type Job struct {
	Param []interface{}
	Fn    func(...interface{})
}

func (j *Job) execute() {
	j.Fn(j.Param...)
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

type dispatcher struct {
	jobPool    chan Job
	workerPool []worker
}

func NewDispatcher(jobPool chan Job, maxWorkers int) *dispatcher {
	numWorker := maxWorkers
	if numWorker <= 0 {
		numWorker = defaultMaxWorker
	}

	workers := []worker{}
	for i := 0; i < numWorker; i++ {
		worker := newWorker(jobPool)
		workers = append(workers, worker)
	}
	return &dispatcher{jobPool: jobPool, workerPool: workers}
}

func (d *dispatcher) Start() {
	for i := 0; i < len(d.workerPool); i++ {
		d.workerPool[i].start()
	}
}

func (d *dispatcher) Stop() {
	for i := 0; i < len(d.workerPool); i++ {
		d.workerPool[i].stop()
	}
}
