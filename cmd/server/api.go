package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"

	"strings"

	// bolt "github.com/coreos/bbolt"
	gin_gzip "github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/contrib/secure"
	"github.com/gin-gonic/gin"

	"github.com/jsnathan/getascension/config"
	"github.com/jsnathan/getascension/utils/log"
)

const (
	errInvalidRequest = "invalid-request-format"
)

type MainTemplate struct {
	SiteName  string
	Prerender template.HTML
	Data      template.JS
	Route     string
	IconFont  template.CSS
	PaceJS    template.JS
	PaceCSS   template.CSS
	LogoSVG   template.HTML
}

func setupMainHandler(r *gin.Engine) {
	iconFontCSS := string(MustAsset("assets/css/icons.css"))
	iconFontCSS = strings.Replace(iconFontCSS, "../font/", "/assets/font/", -1)
	iconFontWOFF2 := base64.StdEncoding.EncodeToString(MustAsset("assets/font/icons.woff2"))
	regex := regexp.MustCompile(`url\('[./\w]+\.woff2(?:\?\d+)?'\) format\('woff2'\)`)
	iconFont := template.CSS(regex.ReplaceAllString(iconFontCSS, fmt.Sprintf(
		"url('data:font/woff2;base64,%s') format('woff2')", iconFontWOFF2,
	)))
	logoSVG := template.HTML(MustAsset("assets/logo.svg"))

	siteName := config.MustGetString("ServerTitle")

	mainPageHandler := func(c *gin.Context) {
		// If the user is logged in, pre-resolve the usual data,
		//   and embed it in the HTML
		var data template.JS
		var mainHTML template.HTML
		var imagePaths []string

		// Push assorted files needed to render the page, if HTTP/2 is in use
		pusher := c.Writer.Pusher()
		if pusher != nil {
			opts := &http.PushOptions{
				Header: http.Header{
					// "Accept-Encoding": r.Header["Accept-Encoding"],
				},
			}
			staticPaths := []string{
				"/app/styles.css",
				// "/assets/css/icons.css",
				// "/assets/font/icons.woff2",
				"/app/build.js",
			}
			numStaticPushed := 0
			numImagesPushed := 0
			for _, path := range staticPaths {
				err := pusher.Push(path, opts)
				if err != nil {
					log.Log.Debugw("Failed to execute HTTP/2 PUSH",
						"func", "mainPageHandler",
						"path", path,
						"error", err,
					)
				} else {
					numStaticPushed++
				}
			}
			for _, imgPath := range imagePaths {
				err := pusher.Push(imgPath, opts)
				if err != nil {
					log.Log.Debugw("Failed to execute HTTP/2 PUSH",
						"func", "mainPageHandler",
						"path", imgPath,
						"error", err,
					)
				} else {
					numImagesPushed++
				}
			}
			log.Log.Debugw("Pushed assets via HTTP/2 PUSH",
				"func", "mainPageHandler",
				"num-images-pushed", numImagesPushed,
				"num-images-total", len(imagePaths),
				"num-static-pushed", numStaticPushed,
				"num-static-total", len(staticPaths),
			)
		} else {
			//	log.Log.Debugw("http.Pusher not available on c.Writer",
			//		"func", "mainPageHandler",
			//	)
		}

		c.HTML(200, "main.tmpl", &MainTemplate{
			SiteName:  siteName,
			Route:     c.Request.URL.String(),
			Data:      data,
			Prerender: mainHTML,
			IconFont:  iconFont,
			LogoSVG:   logoSVG,
		})
	}
	noRouteHandler := func(c *gin.Context) {
		path := c.Request.URL.Path
		path = strings.TrimPrefix(path, "/")
		if strings.HasPrefix(path, "assets/") {
			c.AbortWithStatus(404)
			return
		}
		mainPageHandler(c)
	}

	mainPage := r.Group("/")
	mainPage.Use(gin_gzip.Gzip(gin_gzip.DefaultCompression))
	mainPage.GET("/", mainPageHandler)
	// redirectToMainPage := func(c *gin.Context) {
	// 	c.Redirect(302, "/")
	// }
	// _ = redirectToMainPage
	r.NoRoute(noRouteHandler)
}

func setupHTTPSRedirectHandler(r *gin.Engine) {
	serverFQDN := config.MustGetString("ServerFQDN")
	r.Use(secure.Secure(secure.Options{
		AllowedHosts:          []string{serverFQDN, "www." + serverFQDN},
		SSLRedirect:           true,
		SSLHost:               serverFQDN,
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
	}))
}

func setupAssetsHandlers(r *gin.Engine, socialRoot string) {
	r.Static("/assets", filepath.Join(socialRoot, "assets"))
	// todo: enable this optimization,
	//   if the lifetime of the JS is expected to be as long as server uptime
	// r.StaticFS("/assets", &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "assets"})
	staticNoCache := func(r *gin.Engine, fn string) {
		path := "/app/" + fn
		r.GET(path, func(c *gin.Context) {
			raw, err := ioutil.ReadFile(filepath.Join(socialRoot, "assets", fn))
			if err != nil {
				c.String(500, "")
				return
			}
			rawGz := new(bytes.Buffer)
			gzw := gzip.NewWriter(rawGz)
			if _, err = io.Copy(gzw, bytes.NewReader(raw)); err == nil {
				if err = gzw.Close(); err == nil {
					c.Header("Content-Encoding", "gzip")
					c.Header("Vary", "Accept-Encoding")
					raw = rawGz.Bytes()
				}
			}
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			if filepath.Ext(fn) == ".js" {
				c.Data(200, "application/javascript", raw)
			} else if filepath.Ext(fn) == ".png" {
				c.Data(200, "image/png", raw)
			} else if filepath.Ext(fn) == ".jpg" {
				c.Data(200, "image/jpeg", raw)
			} else if filepath.Ext(fn) == ".svg" {
				c.Data(200, "image/svg+xml", raw)
			} else if filepath.Ext(fn) == ".css" {
				c.Data(200, "text/css", raw)
			} else if filepath.Ext(fn) == ".woff" {
				c.Data(200, "application/font-woff", raw)
			} else {
				c.String(404, "")
			}
		})
	}
	staticNoCache(r, "logo.svg")
	staticNoCache(r, "bg.jpg")
	staticNoCache(r, "styles.css")
	staticNoCache(r, "build.js")
}

// TODO: the functions below should take a request and use it to log
//         the returned status codes, and the errors (if any)

// sendErrJSON will send a failure reply to the client, indicating
//   the error that occurred during the API call
func sendErrJSON(c *gin.Context, text string) {
	fmt.Println("sending error:", text)
	c.JSON(500, gin.H{
		"status": "error",
		"error":  text,
	})
}

// sendOkJSON will send a simple success reply to the client
func sendOkJSON(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

// sendDataJSON will send a success reply to the client which
//   contains a data field with arbitrary result data
func sendDataJSON(c *gin.Context, data interface{}) {
	raw, err := json.Marshal(map[string]interface{}{
		"status": "ok",
		"data":   data,
	})
	if err != nil {
		sendErrJSON(c, fmt.Sprintf("failed to marshal json response: %s", err))
		return
	}
	rawGz := new(bytes.Buffer)
	gzw := gzip.NewWriter(rawGz)
	if _, err = io.Copy(gzw, bytes.NewReader(raw)); err == nil {
		if err = gzw.Close(); err == nil {
			c.Header("Content-Encoding", "gzip")
			c.Header("Vary", "Accept-Encoding")
			raw = rawGz.Bytes()
		}
	}
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Data(200, "application/json", raw)
	// c.JSON(200, gin.H{
	// 	"status": "ok",
	// 	"data":   data,
	// })
}
