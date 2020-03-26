package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"
	"time"
)

var DeployToken string

func main() {
	u, _ := user.Current()
	log.Println("run with:", u.Username)

	// specified deploy token
	if len(os.Args) > 1 {
		DeployToken = os.Args[1]
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/deploy", deployHandler)
	log.Println(http.ListenAndServe(":8082", mux))
}

func deployHandler(w http.ResponseWriter, req *http.Request) {

	// runtime directory
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return
	}

	token := req.FormValue("token")
	if token != DeployToken {
		_, err := w.Write([]byte("deploy token error."))
		if err != nil {
			log.Println(err)
		}
		return
	}
	script := req.FormValue("script")
	if script == "" {
		_, err := w.Write([]byte("script not found."))
		if err != nil {
			log.Println(err)
		}
		return
	}

	var cacheID = time.Now().Format("2006010215040506")
	// deploy with file
	upload, uploadHeader, err := req.FormFile("file")
	if err == nil {
		defer upload.Close()
		var filename = cacheID + uploadHeader.Filename
		fi, err := os.Create(dir + "/data/" + filename)
		if err != nil {
			log.Println(err)
		}
		_ = fi.Close()
		if _, err := io.Copy(fi, upload); err != nil {
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				log.Println(err)
			}
			return
		}
		// unzip
		fileExt := path.Ext(filename)
		if err := Unzip(dir+"/data/"+filename, dir+"/cache/"+strings.TrimSuffix(filename, fileExt)); err != nil {
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				log.Println(err)
			}
			return
		}
	}

	cmd := &exec.Cmd{}
	cmd = exec.Command(dir+"/scripts/"+script+".sh", cacheID)
	cmd.Stdout = w
	if err := cmd.Start(); err != nil {
		_, err := w.Write([]byte(fmt.Sprintf("run script with error: %s ", err)))
		if err != nil {
			log.Println(err)
		}
		return
	}
	if err := cmd.Wait(); err != nil {
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Println(err)
		}
		return
	}
	_, _ = w.Write([]byte("success"))
}
