package main

import (
	"fmt"
	"syscall/js"
	"time"

	"golang.org/x/crypto/ssh"
)

type Client struct {
	Name string
}

func NewClient(this js.Value, args []js.Value) any {
	fmt.Println("NewClient")
	c := &Client{Name: args[0].String()}
	return map[string]any{
		"connect": js.FuncOf(c.Connect),
		"onData":  js.FuncOf(c.OnData),
	}
}

func add() {
	fmt.Println("add")
}

func (c *Client) Connect(this js.Value, args []js.Value) any {
	fmt.Println("Connecting")
	fmt.Println(c.Name)

	cli, sess, err := connectToHost("gustavo", "localhost:22")
	fmt.Println(cli)
	fmt.Println(sess)
	fmt.Println(err)

	return nil
}

func connectToHost(user, host string) (*ssh.Client, *ssh.Session, error) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password("a")},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, session, nil
}

func (c *Client) OnData(this js.Value, args []js.Value) any {
	if args[0].Type() != js.TypeFunction {
		panic("OnData panic")
	}

	go func() {
		for {
			args[0].Invoke("toma")
			time.Sleep(time.Second)
		}
	}()

	return nil
}

func main() {
	fmt.Println("Hello2")

	js.Global().Set("NewClient", js.FuncOf(NewClient))

	select {}
}
