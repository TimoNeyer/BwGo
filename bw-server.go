package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"time"
)

type Server struct {
	port     int
	hostname string
	unlocked float64
	Token    string
}

type ConnectionInit struct {
	Success bool   `json:"success"`
	Token   string `json:"raw"`
}

type BasicConnection struct {
	Success bool `json:"success"`
}

func (i *Server) closeServer(handler chan string) error {
	time.Sleep(time.Duration(i.unlocked))
	serverUrl := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", i.hostname, i.port),
		Path:   "lock",
	}
	req := http.Request{
		Method: "Post",
		URL:    &serverUrl,
	}
	client := http.Client{}
	resp, err := client.Do(&req)
	if err != nil {
		handler <- err.Error()
	}
	ret := BasicConnection{}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		handler <- err.Error()
	}
	err = json.Unmarshal(content, &ret)
	if err != nil {
		handler <- err.Error()
	} else if !ret.Success {
		handler <- "Unsuccessful returned by Server"
	}
	close(handler)
	return nil
}

func (i *Server) start_server(handler chan error) error {
	if i.unlocked >= 0 {
		return nil
	}
	com := exec.Command("/usr/bin/bw", "serve") //, fmt.Sprintf("--hostname %s", i.hostname), fmt.Sprintf("--port %d", i.port))
	go func() {
		if out, err := com.Output(); err != nil {
			fmt.Printf("establishing server code: %s\nstdout: %s\n", err.Error(), com.Stdout)
			handler <- fmt.Errorf("Failed to start bitwarden cli server with\ncode: %s\nstdout: %s\nstderr; %s", err.Error(), out, com.Stderr)
		}
	}()
	handler <- nil
	return nil
}

func (i *Server) unlock_server() error {
	shandler := make(chan error)
	defer close(shandler)
	go i.start_server(shandler)
	time.Sleep(time.Second * 5)
	fmt.Println("after dark")
	if status := <-shandler; status != nil {
		return status
	}
	serverUrl := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", i.hostname, i.port),
		Path:   "unlock",
	}
	passwd, err := get_password()
	if err != nil {
		return err
	}
	bodyMap := map[string]string{"password": string(passwd)}
	bodyBytes, err := json.Marshal(bodyMap)
	if err != nil {
		return fmt.Errorf("unable to encode data")
	}
	body := io.NopCloser(io.Reader(bytes.NewBuffer(bodyBytes)))
	req := http.Request{
		Method: "Post",
		URL:    &serverUrl,
		Body:   body,
		Header: http.Header{},
	}
	req.Header.Add("Content-Type", "application/json")
	passwd = nil
	client := http.Client{}
	resp, err := client.Do(&req)
	if err != nil {
		return fmt.Errorf("UnlockError: %s", err.Error())
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ReadError: %s", err.Error())
	} else if len(content) == 0 {
		return fmt.Errorf("HttpError: returned content with size 0")
	}
	con := ConnectionInit{}
	if err = json.Unmarshal(content, &con); err != nil {
		return fmt.Errorf("ReadError (Marshal): %s", err.Error())
	} else if !con.Success {
		return fmt.Errorf("APIError: Expected Success, received %v", con.Success)
	}
	i.Token = con.Token
	fmt.Println("Successfully unlocked Vault")
	i.unlocked = 600
	chandler := make(chan string)
	go i.closeServer(chandler)
	for message := range chandler {
		return errors.New(message)
	}
	return nil
}
