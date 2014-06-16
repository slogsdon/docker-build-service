package build

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

func Create(app, lang string, code []byte) {
	var ext string

	switch lang {
	case "go":
		ext = "go"
	case "elixir":
		ext = "ex"
	case "ruby":
		ext = "rb"
	case "python":
		ext = "py"
	case "racket":
		ext = "rkt"
	}

	os.MkdirAll("./builds/"+app, 0755)
	ioutil.WriteFile("./builds/"+app+"/main."+ext, code, 0755)

	dockerfile, _ := ioutil.ReadFile("./dockerfiles/" + lang + ".Dockerfile")
	ioutil.WriteFile("./builds/"+app+"/Dockerfile", dockerfile, 0755)
	return
}

func Build(app string) (response, error) {
	return callCmd(exec.Command("docker", "build", "-t", app+"/app", "./builds/"+app))
}

func Run(app string) (response, error) {
	return callCmd(exec.Command("docker", "run", "-t", app+"/app"))
}

func GetFullResp(app string, buildresp, runresp response) fullResp {
	return fullResp{
		AppId: app,
		Build: buildresp,
		Run:   runresp,
	}
}

func callCmd(cmd *exec.Cmd) (response, error) {
	var (
		stdout, stderr bytes.Buffer
		err            error
	)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	done := make(chan error)
	go func() {
		done <- cmd.Run()
	}()
	select {
	case <-time.After(5 * time.Minute):
		err = cmd.Process.Kill()
		<-done // allow goroutine to exit
		if err == nil {
			err = errorTimeout{s: "timeout"}
		}
	case err := <-done:
		err = err
	}

	return response{Success: err == nil, Result: stdout.String()}, err
}

type errorTimeout struct {
	s string
}

func (e errorTimeout) Error() string {
	return e.s
}

type response struct {
	Success bool   `json:"success,omitempty"`
	Result  string `json:"result,omitempty"`
}

type fullResp struct {
	AppId string   `json:"app_id,omitempty"`
	Build response `json:"success,omitempty"`
	Run   response `json:"result,omitempty"`
}
