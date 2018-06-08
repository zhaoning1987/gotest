package multithreadwithorder

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
	jobCh chan Job
	quit  chan bool
}

func newWorker(jobCh chan Job) *worker {
	return &worker{
		jobCh: jobCh,
		quit:  make(chan bool)}
}

func (w *worker) start() {
	go func() {
		for {
			select {
			case job, ok := <-w.jobCh:
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
	JobPool    []chan Job
	workerPool []worker
}

func NewDispatcher(maxWorkers int) *dispatcher {
	numWorker := maxWorkers
	if numWorker <= 0 {
		numWorker = defaultMaxWorker
	}

	workers := []worker{}
	chPool := []chan Job{}
	for i := 0; i < numWorker; i++ {
		ch := make(chan Job)
		worker := newWorker(ch)
		workers = append(workers, *worker)
		chPool = append(chPool, ch)
	}
	return &dispatcher{JobPool: chPool, workerPool: workers}
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
