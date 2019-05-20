package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/powerman/structlog"
)

//nolint:gochecknoglobals
var (
	cmd = strings.TrimSuffix(path.Base(os.Args[0]), ".test")
	ver string // set by ./build
	log = structlog.New()
	cfg struct {
		version  bool
		logLevel string
		port     int
	}
)

// fatalFlagValue report invalid flag values in same way as flag.Parse().
func fatalFlagValue(msg, name string, val interface{}) {
	fmt.Fprintf(os.Stderr, "invalid value %#v for flag -%s: %s\n", val, name, msg)
	flag.Usage()
	os.Exit(2)
}

// Init provides common initialization for both app and tests.
func Init() {
	flag.BoolVar(&cfg.version, "version", false, "print version")
	flag.StringVar(&cfg.logLevel, "log.level", "debug", "log `level` (debug|info|warn|err)")
	flag.IntVar(&cfg.port, "port", 8080, "`port` to listen")

	log.SetDefaultKeyvals(
		structlog.KeyUnit, "main",
	)
}

func main() {
	Init()
	flag.Parse()

	switch {
	case cfg.port <= 0: // Don't support dynamic ports.
		fatalFlagValue("must be > 0", "port", cfg.port)
	case cfg.version: // Must be checked after all other flags for ease testing.
		fmt.Println(cmd, ver, runtime.Version())
		os.Exit(0)
	}

	// Wrong log.level is not fatal, it will be reported and set to "debug".
	structlog.DefaultLogger.SetLogLevel(structlog.ParseLevel(cfg.logLevel))
	log.Info("started", "version", ver)

	http.HandleFunc("/", happy)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func happy(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "smile )")
}
