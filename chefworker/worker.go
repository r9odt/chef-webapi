package chefworker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	web "github.com/JIexa24/chef-webapi"
	"github.com/JIexa24/chef-webapi/database"
	"github.com/JIexa24/chef-webapi/database/interfaces"
	"github.com/JIexa24/chef-webapi/httpserver/controller"
	"github.com/JIexa24/chef-webapi/logging"
)

// Worker describes a worker object.
type Worker struct {
	Threads               int
	DB                    *database.DBConnector
	Logger                logging.Logger
	Dir                   string
	SSHClient             *ssh.ClientConfig
	sshMux                sync.Mutex
	StopRequest           chan struct{}
	waitGroup             sync.WaitGroup
	mux                   sync.Mutex
	pause                 time.Duration
	fileWatcherPause      time.Duration
	fileWatcherExpireDays int64
	sshUser               string
	sshKeyPath            string
	sshKnownHostsFile     string
}

// CreateSSHConfig configure ssh config for worker.
func (w *Worker) CreateSSHConfig() {
	w.sshMux.Lock()
	defer w.sshMux.Unlock()
	// Initialize ssh client.
	key, err := os.ReadFile(w.sshKeyPath)
	if err != nil {
		w.Logger.Errorf("Unable to read private key %s: %s", w.sshKeyPath, err.Error())
		w.SSHClient = nil
		return
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		w.Logger.Errorf("Unable to parse private key %s: %s", w.sshKeyPath, err.Error())
		w.SSHClient = nil
		return
	}

	// Initialize the SSH config.
	config := &ssh.ClientConfig{
		User: w.sshUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	w.SSHClient = config
}

func (w *Worker) getSSHConfig() *ssh.ClientConfig {
	w.sshMux.Lock()
	defer w.sshMux.Unlock()
	return w.SSHClient
}

// NewChefWorker return a new Worker.
func NewChefWorker(l logging.Logger, db *database.DBConnector, directory, sshUser,
	sshKeyPath, sshKnownHostsFile string, expireWorkerLogDays int64) *Worker {

	if expireWorkerLogDays <= 1 {
		expireWorkerLogDays = 2
	}
	worker := &Worker{
		Threads:               runtime.NumCPU(),
		DB:                    db,
		Logger:                l,
		SSHClient:             nil,
		Dir:                   directory,
		StopRequest:           make(chan struct{}),
		waitGroup:             sync.WaitGroup{},
		mux:                   sync.Mutex{},
		sshMux:                sync.Mutex{},
		pause:                 time.Second * 2,
		fileWatcherPause:      time.Hour * 24,
		sshUser:               sshUser,
		sshKeyPath:            sshKeyPath,
		sshKnownHostsFile:     sshKnownHostsFile,
		fileWatcherExpireDays: expireWorkerLogDays,
	}
	worker.CreateSSHConfig()
	// Return a worker instance.
	return worker
}

// Start launches Worker.
func (w *Worker) Start() {
	// Updating the taskbase at startup so that there are no hoped tasks,
	// everyone put ERROR.
	if err := w.DB.UpdateTaskStatusAtStartup(); err != nil {
		w.Logger.Errorf(
			"[worker] worker fatal [w.DB.UpdateTaskStatusAtStartup]: %s",
			err.Error())
		return
	}

	// Check the existence of a directory for Worker logs.
	_, err := os.Stat(w.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(w.Dir, 0755)
			if err != nil {
				w.Logger.Errorf(
					"[worker] worker fatal [os.MkdirAll]: %s",
					err.Error())
				return
			}
		} else {
			w.Logger.Errorf(
				"[worker] worker fatal [os.Stat]: %s",
				err.Error())
			return
		}
	}

	// Starting the thread pool.
	for thread := 0; thread < w.Threads; thread++ {
		w.waitGroup.Add(1)
		go w.worker(thread)
	}
	w.waitGroup.Add(1)
	go w.fileWatcher()
}

// getWaitingTask Returns one Task out of the queue. Function - broker.
func (w *Worker) getWaitingTask() *interfaces.TaskEntry {
	w.mux.Lock()
	defer w.mux.Unlock()
	moduleName := "ChefWorker"
	module, err := w.DB.GetAppModuleByName(moduleName)
	if err != nil {
		w.Logger.Errorf(
			"[worker task manager] getWaitingTask [w.DB.GetAppModuleByName]: %s",
			err.Error())
		return nil
	}
	if module == nil {
		w.Logger.Errorf("Missing required module %s - recreating modules",
			moduleName)
		err := w.DB.CreateAppModules()
		if err != nil {
			w.Logger.Errorf(
				"[worker task manager] getWaitingTask [w.DB.CreateAppModules]: %s",
				err.Error())
		}
		return nil
	}
	if !module.IsON {
		return nil
	}

	// Request to the database.
	task, err := w.DB.GetWaitingTask()
	if err != nil {
		w.Logger.Errorf(
			"[worker task manager] getWaitingTask [w.DB.GetWaitingTask]: %s",
			err.Error())
		return nil
	}

	// If you have a task - update status.
	if task != nil {
		if err := w.DB.UpdateTaskStatusByID(task.ID, "InProgress"); err != nil {
			w.Logger.Errorf(
				"[worker task manager] getWaitingTask [w.DB.UpdateTaskStatusByID]: %s",
				err.Error())
			return nil
		}
	}

	return task
}

// Worker This is a Worker Function.
func (w *Worker) worker(workerID int) {
	defer w.waitGroup.Done()
	w.Logger.Infof("Worker %d started", workerID)

	// Flow check request for completion of work.
	exit := false
	go func() {
		<-w.StopRequest
		exit = true
	}()

	timeoutChannel := make(chan bool)
	// Worker's work cycle.
	for !exit {
		// Get the task.
		task := w.getWaitingTask()
		// If there is a task - perform.
		if task != nil {
			w.Logger.Infof("[worker %d] Deploy task for %s: %s", workerID, task.Resource, task.Name)

			var resource string
			switch task.Resource {
			default:
				resource = "unknown"
			case "roles":
				resource = "role"
			case "nodes":
				resource = "name"
			}
			// Request with API list of nodes for processing.
			list, err := controller.GetResourceBySearch("node",
				resource+":"+task.Name, true)
			if err != nil {
				w.Logger.Errorf(
					"[worker %d] worker [controller.GetResourceBySearch]: %v", workerID,
					err.Error())
				if err := w.DB.UpdateTaskStatusByID(task.ID, "Error"); err != nil {
					w.Logger.Errorf(
						"[worker %d] worker [w.DB.UpdateTaskStatusByID]: %s",
						err.Error())
				}
				continue
			}

			// Sign of the presence of an error.
			taskHasError := false

			// Log file template
			resourceFile := fmt.Sprintf("worker-%s-%s-%s.log", task.ID, task.Resource, task.Name)
			path := filepath.Join(w.Dir, resourceFile)
			file, err := os.OpenFile(path,
				os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				w.Logger.Errorf(
					"[worker %d] writeLog [os.OpenFile]: %v",
					workerID, err.Error())
				if err := w.DB.UpdateTaskStatusByID(task.ID, "Error"); err != nil {
					w.Logger.Errorf(
						"[worker %d] worker [w.DB.UpdateTaskStatusByID]: %s",
						err.Error())
				}
				continue
			}

			// The cycle of execution of the task in nodes.
			for _, v := range list {
				// Call Processing Nodes. The sign accumulates.
				taskHasError = taskHasError || w.taskProcessing(workerID, resource,
					task, file, v)
			}

			// If the list is empty.
			if len(list) <= 0 {
				taskHasError = true
				err = fmt.Errorf(
					"deploy nodes not found")
				w.writeLog(file, "Unknown", err.Error())
			}
			// If there is a sign of error - Error status.
			if taskHasError {
				if err := w.DB.UpdateTaskStatusByID(task.ID, "Error"); err != nil {
					w.Logger.Errorf(
						"[worker %d] worker [w.DB.UpdateTaskStatusByID]: %s",
						err.Error())
				}
			} else {
				if err := w.DB.UpdateTaskStatusByID(task.ID, "Complete"); err != nil {
					w.Logger.Errorf(
						"[worker %d] worker [w.DB.UpdateTaskStatusByID]: %s",
						err.Error())
				}
			}
			w.Logger.Infof("[worker %d] Complete task for %s: %s", workerID,
				task.Resource, task.Name)
			// Close the log file.
			file.Close()
			err = w.DB.UploadFile(path, resourceFile)
			if err != nil {
				w.Logger.Errorf(
					"[worker %d] writeLog [w.DB.UploadFile]: %v",
					workerID, err.Error())
			}
		}
		go func() {
			// Waiting one of two signals.
			select {
			// Pause not to load the CPU.
			case <-time.After(w.pause):
			case <-w.StopRequest:
			}
			timeoutChannel <- true
		}()
		// Waiting for signal.
		<-timeoutChannel
	}
	// Stop Worker.
	w.Logger.Infof("Worker %d stopped", workerID)
}

func (w *Worker) taskProcessing(workerID int, resource string,
	task *interfaces.TaskEntry,
	file *os.File, v map[string]interface{}) bool {

	// Collect the array of candidates.Order IS FQDN, IP or Name of Node.
	candidates := make([]string, 0)
	if v["fqdn"] != nil {
		candidates = append(candidates, v["fqdn"].(string))
	}
	if v["ipaddress"] != nil {
		candidates = append(candidates, v["ipaddress"].(string))
	}
	if v["name"] != nil {
		candidates = append(candidates, v["name"].(string))
	}

	// Collect the command to execute.
	chefCmd := "sudo chef-client"
	switch task.Resource {
	case "nodes":
		// selective resource deplosity.
		if task.SelectedResource {
			resources := make([]string, 0)
			err := json.Unmarshal([]byte(task.Resources), &resources)
			if err == nil {
				if len(resources) <= 0 {
					w.writeLog(file, task.Name, "Nothing to deploy")
					return true
				}
				// Overload of the Resource List.
				runlist := web.ConcatenateStringWithDelimeter(",", resources)
				chefCmd = fmt.Sprintf("%s -o %s", chefCmd, runlist)
			}
		}
	case "roles":
		// If only role is rolled out.
		if task.OnlyResource {
			chefCmd = fmt.Sprintf("%s -o role[%s]", chefCmd, task.Name)
		}
		// If selective nodes roll out.
		if task.SelectedResource {
			resources := make([]string, 0)
			err := json.Unmarshal([]byte(task.Resources), &resources)
			if err == nil {
				if len(resources) <= 0 {
					w.writeLog(file, task.Name, "Nothing to deploy")
					return true
				}
				nodeFound := false
				// If there is no name nodes in the processing requests - then we skip.
				// This is not considered an error.
				for i := range resources {
					if v["name"].(string) == resources[i] {
						nodeFound = true
					}
				}
				if !nodeFound {
					return false
				}
			}
		}
	}

	sshConfig := w.getSSHConfig()

	if sshConfig == nil {
		err := fmt.Errorf(
			"[worker %d] taskProcessing %v: SSH Client does not configured",
			workerID, candidates)
		w.Logger.Errorf(err.Error())
		w.writeLog(file, "Unknown", err.Error())
		return true
	}

	// Cycle connections to node on candidates.
	for _, host := range candidates {
		w.Logger.Infof(
			"[worker %d] taskProcessing [Try host]: %v", workerID, host)

		// Connection unit.
		client, err := ssh.Dial("tcp", host+":22", sshConfig)
		if err != nil {
			w.Logger.Errorf(
				"[worker %d] taskProcessing [ssh.Dial]: (%s) %v", workerID, host,
				err.Error())
			continue
		}

		// Open the session, request the host name.
		b, err := w.clientSession(client, "hostname -f")
		if err != nil {
			w.Logger.Errorf(
				"[worker %d] taskProcessing [w.clientSession]: (%s) %v", workerID, host,
				err.Error())
			client.Close()
			continue
		}

		// Checking that went on the node.
		hostname := strings.TrimRight(b.String(), "\n")
		if hostname != v["fqdn"] && hostname != v["name"] {
			err = fmt.Errorf(
				"(%s) Hostname (%s) does not equal to name (%s) or fqdn (%s) node",
				host, hostname, v["name"], v["fqdn"])
			w.writeLog(file, hostname, err.Error())
			w.Logger.Errorf("[worker %d] %v", workerID, err.Error())
			client.Close()
			continue
		}

		// If you went to that server - the call command.
		b, err = w.clientSession(client, chefCmd)
		if err != nil {
			switch err.(type) {
			default:
			case *ssh.ExitError, *ssh.ExitMissingError:
				w.writeLog(file, hostname, b.String())
				return true
			}
			w.Logger.Errorf(
				"[worker %d] worker [w.clientSession]: (%s) %v", workerID, host,
				err.Error())
			client.Close()
			continue
		}
		// Save the log to the file.
		w.writeLog(file, hostname, b.String())
		// Close the session.
		client.Close()
		return false
	}

	err := fmt.Errorf(
		"[worker %d] taskProcessing %v: Task processing complete without success",
		workerID, candidates)
	w.Logger.Errorf(err.Error())
	w.writeLog(file, "Unknown", err.Error())
	return true
}

func (w *Worker) writeLog(file *os.File, hostname, log string) {
	fmt.Fprintf(file, "----------!!![%s]!!!----------\n", hostname)
	fmt.Fprintf(file, "%s\n", log)
}

func (w *Worker) clientSession(client *ssh.Client, cmd string) (bytes.Buffer, error) {
	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	var b bytes.Buffer
	if err != nil {
		w.Logger.Errorf(
			"clientSession [client.NewSession]: %s",
			err.Error())
		return b, err
	}
	defer session.Close()
	// Save the output stream to the buffer
	session.Stdout = &b
	// Run the command.
	if err := session.Run(cmd); err != nil {
		w.Logger.Errorf(
			"clientSession [session.Run(%s)]: %v", cmd,
			err.Error())
		return b, err
	}
	return b, nil
}

// Stop is a function of stopping the Worker.
func (w *Worker) Stop() {
	close(w.StopRequest)
	w.waitGroup.Wait()
}

func (w *Worker) fileWatcher() {
	defer w.waitGroup.Done()
	// Flow check request for completion of work.
	exit := false
	go func() {
		<-w.StopRequest
		exit = true
	}()
	timeoutChannel := make(chan bool)
	for !exit {
		// Find all Worker logs.
		files, err := os.ReadDir(w.Dir)
		if err != nil {
			w.Logger.Error(err.Error())
		}
		for _, file := range files {
			fileInfo, err := file.Info()
			if err != nil {
				continue
			}
			if !(file.Type().IsRegular() && !strings.HasPrefix(file.Name(), ".") &&
				fileInfo.ModTime().Before(time.Now().AddDate(0, 0, -int(w.fileWatcherExpireDays)))) {
				continue
			}
			if w.DB.CheckFile(file.Name()) {
				if err := os.Remove(filepath.Join(w.Dir, file.Name())); err != nil {
					w.Logger.Errorf("fileWatcher [os.Remove]: %s", err.Error())
				}
			} else {
				if err := w.DB.UploadFile(filepath.Join(w.Dir, file.Name()), file.Name()); err != nil {
					w.Logger.Errorf("fileWatcher [w.DB.UploadFile]: [%s %s] %s", w.Dir, file.Name(), err.Error())
				}
			}
		}
		go func() {
			// Waiting one of two signals.
			select {
			// Pause not to load the CPU.
			case <-time.After(w.fileWatcherPause):
			case <-w.StopRequest:
			}
			timeoutChannel <- true
		}()
		<-timeoutChannel
	}
}
