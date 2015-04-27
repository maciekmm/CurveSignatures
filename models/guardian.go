package models

import (
	"sync"
	"sync/atomic"
)

type Guardian struct {
	lock  *sync.RWMutex
	files map[string]*File
}

type File struct {
	group      sync.WaitGroup
	Lock       sync.RWMutex
	OnceCaller Once
}

// This is "stolen" from sync.Once
type Once struct {
	m    sync.Mutex
	done uint32
}

func (o *Once) Do(f func()) bool {
	if atomic.LoadUint32(&o.done) == 1 {
		return false
	}
	// Slow-path.
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
		return true
	}
	return false
}

// End of things stolen

// Returns new instance of Guardian
func New() *Guardian {
	mut := &sync.RWMutex{}
	files := make(map[string]*File)
	gr := &Guardian{mut, files}
	return gr
}

// Gets or creates a lock for file
// You have to call Done() in order to free resources
func (g *Guardian) GetOrCreate(fil string) *File {
	g.lock.RLock()
	val, ok := g.files[fil]
	if ok {
		val.group.Add(1)
	}
	g.lock.RUnlock()
	if !ok {
		val = &File{sync.WaitGroup{}, sync.RWMutex{}, Once{}}
		g.lock.Lock()
		defer g.lock.Unlock()
		val.group.Add(1)
		g.files[fil] = val
		go func() {
			// When all tasks are finished delete from map
			val.group.Wait()
			g.lock.Lock()
			defer g.lock.Unlock()
			delete(g.files, fil)
		}()

	}
	return val
}

func (f *File) Done() {
	f.group.Done()
}
