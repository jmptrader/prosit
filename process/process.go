package process

import (
	"container/list"
	"prosit/cerr"
	"sync"
)

var l *list.List
var lMutex sync.RWMutex

func init() {
	l = list.New()
}

type Process struct {
	Id        string `json:"id"`
	Pid       int    `json:"pid"`
	RunAs     string `json:"runAs"`
	Run       string `json:"run"`
	Folder    string `json:"folder"`
	Error     string `json:"error"`
	Started   int64  `json:"started"`
	IsRunning bool   `json:"isRunning"`
	AlertID   string `json:"alertID"`
}

func AddProcess(id, run, folder, alertID, runAs string) error {

	lMutex.Lock()
	defer lMutex.Unlock()

	// we create an internal process
	intProc, err := newInternalProcess(id, run, folder, alertID, runAs)

	if err != nil {
		return err
	}

	// we add to the list
	l.PushBack(intProc)

	// we start the Process
	intProc.start()

	return nil
}

func ProcessExists(id string) bool {

	lMutex.RLock()
	defer lMutex.RUnlock()

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		var tmpProc = e.Value.(*internalProcess)

		if tmpProc.id == id {
			return true
		}
	}

	return false
}

func StopProcess(id string) error {

	lMutex.RLock()
	defer lMutex.RUnlock()

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		var tmpProc = e.Value.(*internalProcess)

		if tmpProc.id == id {
			return tmpProc.stop()
		}
	}

	return nil
}

func RestartProcess(id string) error {

	lMutex.RLock()
	defer lMutex.RUnlock()

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		var tmpProc = e.Value.(*internalProcess)

		if tmpProc.id == id {
			return tmpProc.restart()
		}
	}

	return nil
}

func ListProcesses() []Process {

	lMutex.RLock()
	defer lMutex.RUnlock()

	ret := make([]Process, 0)

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		var tmpIntProcess = e.Value.(*internalProcess)

		var tmpProcess = &Process{}
		tmpProcess.Id = tmpIntProcess.id
		tmpProcess.Pid = tmpIntProcess.pid
		tmpProcess.Run = tmpIntProcess.fullPath
		tmpProcess.RunAs = tmpIntProcess.runAs
		tmpProcess.Folder = tmpIntProcess.folder
		if tmpIntProcess.err != nil {
			tmpProcess.Error = tmpIntProcess.err.Error()
		}
		tmpProcess.Started = tmpIntProcess.lastStarted
		tmpProcess.IsRunning = tmpIntProcess.isRunning
		tmpProcess.AlertID = tmpIntProcess.alertID

		ret = append(ret, *tmpProcess)
	}

	return ret
}

func GetProcessLogs(id string) ([]LogItem, error) {

	lMutex.RLock()
	defer lMutex.RUnlock()

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		var tmpIntProcess = e.Value.(*internalProcess)

		if tmpIntProcess.id == id {
			return tmpIntProcess.stdout.LogList(), nil
		}
	}

	return nil, cerr.NewBadRequestError("Process '%s' not found", id)
}

func GetProcessErrors(id string) ([]LogItem, error) {

	lMutex.RLock()
	defer lMutex.RUnlock()

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		var tmpIntProcess = e.Value.(*internalProcess)

		if tmpIntProcess.id == id {
			return tmpIntProcess.stderr.LogList(), nil
		}
	}

	return nil, cerr.NewBadRequestError("Process '%s' not found", id)
}
