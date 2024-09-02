/*
 *@author  chengkenli
 *@project pjstack
 *@package main
 *@file    ConnectSSH
 *@date    2024/9/2 13:58
 */

package proc

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
)

type ConSSH struct {
	User           string
	Host           string
	Port           int
	PrivateKeyFile string
	Command        string
}

// Connect 根据公钥链接远程服务器并执行命令返回结果
func Connect(c *ConSSH) interface{} {
	// 连接到SSH服务器
	client, err := connectSSH(c.User, c.Host, c.Port, publicKeyAuthFunc(c.PrivateKeyFile))
	if err != nil {
		fmt.Println(err.Error())
	}
	defer client.Close()
	// 执行命令
	output := runCommand(client, c.Command)
	return output
}

// PublicKeyAuthFunc 连接到SSH服务器
func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	// 读取私钥文件
	key, err := ioutil.ReadFile(kPath)
	if err != nil {
		fmt.Println(err.Error())
	}
	// 创建SSH签名器
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		fmt.Println(err.Error())
	}

	return ssh.PublicKeys(signer)
}

// connectSSH 建立SSH客户端连接
func connectSSH(user, host string, port int, auth ssh.AuthMethod) (*ssh.Client, error) {
	var addr = fmt.Sprintf("%s:%d", host, port)
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 在生产环境中应该更安全地处理HostKey
	}

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// runCommand 执行远程命令
func runCommand(client *ssh.Client, command string) string {
	sess, err := client.NewSession()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer sess.Close()

	output, err := sess.CombinedOutput(command)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(output)
}
