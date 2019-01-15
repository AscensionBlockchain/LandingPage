package main

// TODO: enable this in production //go:generate build_web -q
//go:generate build_web -q
//go:generate go run github.com/AscensionBlockchain/LandingPage/cmd/server/tools/build_deploy

import (
	"crypto/tls"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/integrii/flaggy"

	"github.com/AscensionBlockchain/LandingPage/config"
	"github.com/AscensionBlockchain/LandingPage/utils"
	"github.com/AscensionBlockchain/LandingPage/utils/log"
)

var socialRoot string

func main() {
	defer log.Log.Sync()

	var insecureMode, showVersion, noCache bool
	var addr, port string = "127.0.0.1", "6050"
	flaggy.DefaultParser.ShowVersionWithVersionFlag = false
	flaggy.Bool(&insecureMode, "", "insecure", "Run server without TLS (only for testing)")
	flaggy.String(&port, "p", "port", "Select port for insecure mode")
	flaggy.String(&addr, "a", "addr", "Select addr for insecure mode")
	flaggy.Bool(&showVersion, "", "version", "Show the version information for Ascension Server")
	flaggy.Bool(&noCache, "nc", "no-cache", "disable redis cache use in lookup stage")
	flaggy.Parse()

	if showVersion {
		fmt.Printf("Version: %s", config.Version)
		return
	}

	socialRoot := os.Getenv("SOCIALDIR")
	if socialRoot == "" {
		log.Log.Fatal("Please set SOCIALDIR environment variable before starting this server; the database will be stored there")
	}
	err := os.MkdirAll(socialRoot, 0755)
	if err != nil {
		log.Log.Fatal("Unable to ensure directory SOCIALDIR at", socialRoot, "because:", err)
	}

	restoreAssets(socialRoot)
	config.LoadFrom(socialRoot)

	go func() {
		var addr string
		if config.GetBool("InsideDocker") {
			addr = "0.0.0.0:7070"
		} else {
			addr = "127.0.0.1:7070"
		}
		// serve pprof endpoints until program exit
		log.Log.Warnw("pprof endpoint died",
			"func", "main",
			"error", http.ListenAndServe(addr, nil),
		)
	}()

	buf := make([]byte, 2e6)
	dumpStackTrace := func() {
		buf = buf[:runtime.Stack(buf[:cap(buf)], true)]
		ts := strings.Replace(utils.Timestamp(), ":", "-", -1)
		fn := fmt.Sprintf("goroutine_stacks.%s.log", ts)
		f, err := os.Create(filepath.Join(socialRoot, fn))
		if err != nil {
			log.Log.Warnw("failed to store goroutine stacktraces", "error", err)
			return
		}
		f.Write(buf)
		f.Close()
	}

	go func() {
		for {
			time.Sleep(60 * time.Second)
			dumpStackTrace()
		}
	}()

	//	sigChan := make(chan os.Signal, 1)
	//	go func() {
	//		for {
	//			<-sigChan
	//			dumpStackTrace()
	//		}
	//	}()
	//	signal.Notify(sigChan, syscall.SIGUSR1)

	r := gin.Default()
	// r.Use(gzip.Gzip(gzip.DefaultCompression))
	// r.Use(limit.MaxAllowed(20))

	loadAllTemplates(r, socialRoot)
	setupAssetsHandlers(r, socialRoot)

	setupMainHandler(r)

	if insecureMode {
		log.Log.Fatal(
			r.Run(addr + ":" + port),
		)
	} else {
		redirector := gin.Default()
		setupHTTPSRedirectHandler(redirector)
		go redirector.Run(":80")

		certFiles := map[string][]string{
			"getascension.com": []string{
				filepath.Join("/root/getascension.com/fullchain.pem"),
				filepath.Join("/root/getascension.com/privkey.pem"),
			},
		}

		certMap := map[string]*tls.Certificate{}
		for domain, files := range certFiles {
			cert, err := tls.LoadX509KeyPair(files[0], files[1])
			if err != nil {
				panic(err)
			}
			certMap[domain] = &cert
		}

		tlsConfig := &tls.Config{
			PreferServerCipherSuites: true,
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.X25519,
			},
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				if strings.Contains(info.ServerName, "getascension.com") {
					return certMap["getascension.com"], nil
				}
				return nil, fmt.Errorf("no certificate for SNI : <%s>", info.ServerName)
			},
		}

		server := &http.Server{
			Addr:              ":443",
			Handler:           r,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      120 * time.Second,
			IdleTimeout:       120 * time.Second,
			TLSConfig:         tlsConfig,
			MaxHeaderBytes:    1 << 20,
		}

		log.Log.Fatal(
			server.ListenAndServeTLS("", ""),
		)
	}
}

func loadAllTemplates(r *gin.Engine, dataDir string) {
	r.LoadHTMLGlob(filepath.Join(dataDir, "templates", "http", "*.tmpl"))
}

func restoreAssets(socialRoot string) {
	os.MkdirAll(socialRoot, 0755)

	err := RestoreAssets(socialRoot, "templates")
	if err != nil {
		log.Log.Fatalw("Unable to restore 'templates/' in temporary directory; cannot boot",
			"func", "restoreAssets",
			"error", err,
		)
	}

	err = RestoreAssets(socialRoot, "assets")
	if err != nil {
		log.Log.Fatalw("Unable to restore 'assets/' in temporary directory; cannot boot",
			"func", "restoreAssets",
			"error", err,
		)
	}

	err = RestoreAsset(socialRoot, "config.toml")
	if err != nil {
		log.Log.Fatalw("Unable to restore 'config.toml' in temporary directory; cannot boot",
			"func", "restoreAssets",
			"error", err,
		)
	}
}
