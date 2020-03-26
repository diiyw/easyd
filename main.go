package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
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

	// deploy with file
	upload, uploadHeader, err := req.FormFile("file")
	if err != nil {
		log.Println(err)
	}
	fi, err := os.Create(dir + "/data/" + uploadHeader.Filename + time.Now().Format("20060102.zip"))
	if err != nil {
		panic(err)
	}
	if _, err := io.Copy(fi, upload); err != nil {
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Println(err)
		}
		return
	}

	cmd := &exec.Cmd{}
	cmd = exec.Command(dir + "/scripts/" + script + ".sh")
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
