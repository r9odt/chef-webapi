package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	web "github.com/JIexa24/chef-webapi"
	"github.com/JIexa24/chef-webapi/chefworker"
	"github.com/JIexa24/chef-webapi/httpserver"
	"github.com/JIexa24/chef-webapi/logging"
)

var (
	// Version of application.
	Version = "unknown"
	// GoVersion is go version.
	GoVersion = "unknown"
	// GitCommit is git commit.
	GitCommit = "unknown"
)

var (
	chefClientName      string
	chefClientKey       string
	chefURL             string
	sshUser             string
	sshKeyPath          string
	sshKnownHostsFile   string
	address             string
	port                string
	env                 string
	sessionProvider     string
	databaseProvider    string
	sessionExpire       int64
	workerDir           string
	enableProfiler      bool
	profilerAddress     string
	profilerPort        string
	expireWorkerLogDays int64
	databaseName        string
	databaseUser        string
	databasePassword    string
	databaseHost        string
	databasePort        string
	appKeyPath          string
	printVersion        bool
)

// Usage is the Flag Processing Function "-h".
func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [params]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

// parseFlags is the function of processing application parameter flags.
func parseFlags() {
	flag.Usage = usage

	flag.BoolVar(&printVersion, "version", false, "Print version and exit")
	flag.StringVar(&env, "env", "development",
		"Application mode [development, production]. (env: ENV)")
	flag.StringVar(&address, "address", "0.0.0.0",
		"Listen Address. (env: LISTEN_ADDR)")
	flag.StringVar(&port, "port", "3000",
		"Listen Port. (env: LISTEN_PORT)")
	flag.StringVar(&workerDir, "worker-dir", "/app/content/worker",
		"Worker directory path. (env: WORKER_DIR)")
	flag.StringVar(&appKeyPath, "app-key", "/app/keys.key",
		"App private key path. (env: APP_KEY)")
	flag.StringVar(&sshKeyPath, "ssh-key", "/app/id_rsa",
		"SSH key path. (env: SSH_KEY)")
	flag.StringVar(&sshUser, "ssh-user", "web",
		"SSH username. (env: SSH_USER)")
	flag.StringVar(&sshKnownHostsFile, "ssh-known-hosts", "/app/known_hosts",
		"SSH known hosts file. (env: SSH_KNOWN_HOSTS)")
	flag.StringVar(&chefClientKey, "chef-key",
		"/app/chef-key.pem",
		"Chef-Client key path. (env: CHEF_KEY)")
	flag.StringVar(&chefClientName, "chef-client", "web",
		"Chef-Client name. (env: CHEF_CLIENT)")
	flag.StringVar(&chefURL, "chef-url",
		"https://chef.example.com/",
		"Chef server url. (env: CHEF_URL)")
	flag.StringVar(&databaseUser, "database-user",
		`chef`,
		`Database user.	(env: DATABASE_USER)`)
	flag.StringVar(&databasePassword, "database-password",
		`chef`,
		`Database password. (env: DATABASE_Password)`)
	flag.StringVar(&databaseHost, "database-host",
		`127.0.0.1`,
		`Database host.	(env: DATABASE_HOST)`)
	flag.StringVar(&databasePort, "database-port",
		`27017`,
		`Database port.	(env: DATABASE_PORT)`)
	flag.StringVar(&databaseProvider, "database-provider", "mongo",
		`Database provider.
		 (env: DATABASE_PROVIDER)`)
	flag.StringVar(&databaseName, "database-name", "ChefWebApp",
		`Database name.
							(env: DATABASE_NAME)`)
	flag.StringVar(&sessionProvider, "sessions-provider", "database",
		`Sessions provider [database].
		 (env: SESSION_PROVIDER)`)
	flag.Int64Var(&sessionExpire, "session-expire", 86400,
		"Expire session in seconds. (env: SESSION_EXPIRE)")
	flag.Int64Var(&expireWorkerLogDays, "logs-expire", 30,
		"Expire worker logs in days. (env: LOGS_EXPIRE)")
	flag.BoolVar(&enableProfiler, "pprof", false,
		"Enable pprof. (env: PPROF)")
	flag.StringVar(&profilerAddress, "pprof-address", "0.0.0.0",
		"PPROF listen Address. (env: PPROF_ADDR)")
	flag.StringVar(&profilerPort, "pprof-port", "8080",
		"PPROF listen Port. (env: PPROF_PORT)")
	flag.Parse()
	if printVersion {
		fmt.Println(os.Args[0])
		fmt.Println("Version:", Version)
		fmt.Println("Git Commit:", GitCommit)
		fmt.Println("Go Version:", GoVersion)
		os.Exit(0)
	}
	parseEnvs()
}

// parseEnvs is the function of processing application environment.
func parseEnvs() {
	var arg string
	if arg = os.Getenv("LISTEN_ADDR"); arg != "" {
		address = arg
	}
	if arg = os.Getenv("LISTEN_PORT"); arg != "" {
		port = arg
	}
	if arg = os.Getenv("WORKER_DIR"); arg != "" {
		workerDir = arg
	}
	if arg = os.Getenv("APP_KEY"); arg != "" {
		appKeyPath = arg
	}
	if arg = os.Getenv("SSH_KEY"); arg != "" {
		sshKeyPath = arg
	}
	if arg = os.Getenv("SSH_USER"); arg != "" {
		sshUser = arg
	}
	if arg = os.Getenv("SSH_KNOWN_HOSTS"); arg != "" {
		sshKnownHostsFile = arg
	}
	if arg = os.Getenv("CHEF_KEY"); arg != "" {
		chefClientKey = arg
	}
	if arg = os.Getenv("CHEF_CLIENT"); arg != "" {
		chefClientName = arg
	}
	if arg = os.Getenv("CHEF_URL"); arg != "" {
		chefURL = arg
	}
	if arg = os.Getenv("DATABASE_PROVIDER"); arg != "" {
		databaseProvider = arg
	}
	if arg = os.Getenv("DATABASE_NAME"); arg != "" {
		databaseName = arg
	}
	if arg = os.Getenv("DATABASE_USER"); arg != "" {
		databaseUser = arg
	}
	if arg = os.Getenv("DATABASE_PASSWORD"); arg != "" {
		databasePassword = arg
	}
	if arg = os.Getenv("DATABASE_HOST"); arg != "" {
		databaseHost = arg
	}
	if arg = os.Getenv("DATABASE_PORT"); arg != "" {
		databasePort = arg
	}
	if arg = os.Getenv("SESSION_PROVIDER"); arg != "" {
		sessionProvider = arg
	}
	if arg = os.Getenv("SESSION_EXPIRE"); arg != "" {
		intval, err := strconv.Atoi(arg)
		if err == nil {
			sessionExpire = int64(intval)
		}
	}
	if arg = os.Getenv("LOGS_EXPIRE"); arg != "" {
		intval, err := strconv.Atoi(arg)
		if err == nil {
			expireWorkerLogDays = int64(intval)
		}
	}
	if arg = os.Getenv("PPROF"); arg != "" {
		bval, err := strconv.ParseBool(arg)
		if err == nil {
			enableProfiler = bval
		}
	}
	if arg = os.Getenv("PPROF_ADDR"); arg != "" {
		profilerAddress = arg
	}
	if arg = os.Getenv("PPROF_PORT"); arg != "" {
		profilerPort = arg
	}
	if arg = os.Getenv("ENV"); arg != "" {
		env = arg
	}
}

// main is main ^).
func main() {
	parseFlags()
	logger, err := logging.ConfigureLog("stdout", "info", "web")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not configure log: %s\n", err.Error())
		os.Exit(1)
	}
	reloadChannel := make(chan struct{})
	stop := make(chan struct{})
	if enableProfiler {
		pprofListen := fmt.Sprintf("%s:%s", profilerAddress, profilerPort)
		go func() {
			if err := http.ListenAndServe(pprofListen, nil); err != nil && err != http.ErrServerClosed {
				logger.Errorf("Profiler [http.ListenAndServe]: %s", err.Error())
			}
		}()
		logger.Infof("PPROF started on %s:%s",
			profilerAddress, profilerPort)
	}

	web.NewApplication(env, reloadChannel)
	web.App.ConfigureLogger(logger)
	if err = web.App.ConfigureChefClient(chefClientName, chefURL, chefClientKey); err != nil {
		logger.Errorf("Cannot configure application. %s", err.Error())
	}
	if err = web.App.ConfigureDatabase(databaseProvider,
		sessionProvider, databaseName, databaseUser, databasePassword,
		databaseHost, databasePort); err != nil || web.App.DB == nil {
		if err == nil {
			err = fmt.Errorf("database pointer is nil")
		}
		logger.Fatalf("Cannot configure database. %s", err.Error())
	}
	web.App.ConfigureApp(workerDir, appKeyPath, sshKeyPath, sessionExpire)

	chefworker := chefworker.NewChefWorker(logger, web.App.DB, workerDir,
		sshUser, sshKeyPath, sshKnownHostsFile, expireWorkerLogDays)
	chefworker.Start()
	defer chefworker.Stop()
	go reloader(chefworker, stop, reloadChannel)
	server := httpserver.New(address, port, logger)
	server.Listen()
	defer server.Stop()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	logger.Info(fmt.Sprint(<-ch))
	close(stop)
	logger.Info("Service shutting down.")
}

func reloader(worker *chefworker.Worker, stop, reload chan struct{}) {
	for {
		select {
		case <-stop:
			return
		case <-reload:
			worker.CreateSSHConfig()
			_ = web.App.CreateChefClientConfig()
		}
	}
}
