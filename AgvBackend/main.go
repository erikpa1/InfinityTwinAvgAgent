package main

import (
	"crypto/tls"
	"fmt"
	"mime"
	"net/http"
	"os/exec"
	"strings"
	"time"
	"turtle/core/lgr"
	"turtle/core/usersApi"
	"turtle/credentials"
	"turtle/server"
	"turtle/vfs"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite" //Must be kvoli registracii drivera na sqlite3
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	return r
}

func dev_main() {

	lgr.Ok("Working directory: ", vfs.GetWorkingDirectory())

	r := setupRouter()

	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".xyz", "text/plain")
	mime.AddExtensionType(".gzip", "application/x-gzip-compressed")
	mime.AddExtensionType(".gz", "application/x-gzip-compressed")

	usersApi.InitUsersApi(r)

	r.Use(static.Serve("/", static.LocalFile("./static", true)))
	//r.NoRoute(tools.ProxyMiddleware2())

	prefix := "http://"
	port := "8001"

	if credentials.RunHttps() {
		prefix = "https://"
	}

	addr := "0.0.0.0:" + port

	server.RunMyioServer(r)

	lgr.Info("Running in: %s", vfs.GetExeFile())

	srv := &http.Server{
		Addr:         addr, // Change to your desired port
		Handler:      r,
		ReadTimeout:  1000 * time.Second, // StringSet read timeout
		WriteTimeout: 1000 * time.Second, // StringSet write timeout
		IdleTimeout:  1000 * time.Second, // StringSet idle timeout
	}

	if credentials.RunHttps() {
		srv.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	lgr.Ok("------------------------------------")
	lgr.Ok("Mode: %s", gin.Mode())

	lgr.Ok("Running server at: %s", fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", prefix+addr, prefix+addr))
	lgr.Ok("Access at: %s", fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", prefix+"127.0.0.1:"+port, prefix+"127.0.0.1:"+port))
	lgr.Ok("Access at: %s", fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", prefix+"localhost:"+port, prefix+"localhost:"+port))

	if vfs.IsLinux() {
		cmd := exec.Command("hostname", "-I")

		// Capture the output
		output, err := cmd.Output()
		if err != nil {
			lgr.Error(err.Error())

		} else {
			ips := strings.Fields(string(output))

			for _, ip := range ips {
				lgr.Ok("Access at: ", fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", prefix+ip+":"+port, prefix+ip+":"+port))

			}
		}

	}

	if credentials.RunHttps() {
		lgr.Ok("Started HTTP(S) branch")
		error := srv.ListenAndServeTLS("cert.pem", "key.pem")

		if error != nil {
			lgr.ErrorJson(error)
		}

	} else {
		lgr.Ok("Started HTTP branch")
		error := srv.ListenAndServe()

		if error != nil {
			lgr.Error("%s", error)
		}
	}

	lgr.Error("Execution ended")
}

func main() {

	lgr.Info("Starting infinity twin application")
	lgr.Info("DbName: ", credentials.GetDBName())

	dev_main()
}
