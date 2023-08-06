package processor

import "time"

type BackgroundProcess struct {
	state     string
	progress  int
	control   chan bool
	listeners []chan int
}

func New() *BackgroundProcess {
	return &BackgroundProcess{
		state:     "ready",
		progress:  0,
		control:   make(chan bool),
		listeners: make([]chan int, 0),
	}
}

func (bp *BackgroundProcess) Start() {
	if bp.state == "ready" {
		go func() {
			defer func() {
				for _, listener := range bp.listeners {
					close(listener)
				}
				bp.listeners = make([]chan int, 0)
				bp.state = "ready"
			}()

			bp.state = "running"
			bp.progress = 0

			for bp.progress < 100 {
				select {
				case <-bp.control:
					bp.state = "ready"
					return
				default:
					bp.progress += 10
					for _, listener := range bp.listeners {
						listener <- bp.progress
					}
					time.Sleep(time.Second)
				}
			}
		}()
	}
}

func (bp *BackgroundProcess) Stop() {
	if bp.state == "running" {
		bp.control <- true
	}
}

func (bp *BackgroundProcess) GetState() string {
	return bp.state
}

func (bp *BackgroundProcess) GetProgress() int {
	return bp.progress
}

func (bp *BackgroundProcess) Listen() chan int {
	listener := make(chan int)
	bp.listeners = append(bp.listeners, listener)
	return listener
}
