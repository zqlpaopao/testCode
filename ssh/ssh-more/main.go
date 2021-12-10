package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

/*
https://github.com/ssikdar1/telnetlib-go/blob/master/telnetlib.go
这个可以获取到telnet all
*/

func main() {
	c, err := New("11.xx.xx.27", "xxxx", "@xxxxx", 22)
	if err != nil {
		fmt.Println(err)
	}

	s, err := c.newSession()

	err = c.runTerminalSession(s, "more ./mysql.txt")
	fmt.Println(err)

}

// Cli ...
type Cli struct {
	IP       string      //IP地址
	Username string      //用户名
	Password string      //密码
	Port     int         //端口号
	client   *ssh.Client //ssh客户端
}

// New 创建命令行对象
// ip IP地址
// username 用户名
// password 密码
// port 端口号,默认22
func New(ip string, username string, password string, port ...int) (*Cli, error) {
	cli := new(Cli)
	cli.IP = ip
	cli.Username = username
	cli.Password = password
	if len(port) <= 0 {
		cli.Port = 22
	} else {
		cli.Port = port[0]
	}

	return cli, cli.connect()
}

// Run 执行 shell脚本命令
func (c Cli) Run(shell string) (string, error) {
	session, err := c.newSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	buf, err := session.CombinedOutput(shell)
	return string(buf), err
}

// 此方法不行
func sessino(session *ssh.Session) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err := session.RequestPty("vt100", 80, 200, modes); err != nil {

		fmt.Println(err)
	}
	w, err := session.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	r, err := session.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}

	//通过两个channel与ssh session的输入输出流对接，将针对buffer的操作升级为与channel的交互
	in := make(chan string, 2000)
	out := make(chan string, 2000)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		for cmd := range in {
			fmt.Println(cmd, "cm")
			_, err := w.Write([]byte(cmd))
			if err != nil {
				fmt.Println(err)
				return
			}
			//每条命令间隔500ms，保证设备能充分接收命令。（之前设定的200ms，但是会导致H3C的分页指令不生效）
			//但是这里存在一个风险，就是向设备发送命令停顿500ms，会导致设备响应变长。
			//由于发送操作也是异步，并不阻塞上层调用方法，有可能导致读取设备方法提前退出（只要设备在2秒内有回写就会避免这种情况）。
			time.Sleep(time.Millisecond * 500)
		}
	}()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		var (
			buf [65 * 1024]byte
			cur int
		)
		for {
			fmt.Println(122)
			n, err := r.Read(buf[cur:])
			fmt.Println(245)

			if err != nil {
				if err.Error() != "EOF" {
					fmt.Println(err)
				}
				return
			}
			cur += n
			fmt.Println(string(buf[:cur]), 123)
			out <- string(buf[:cur])
			cur = 0
		}
	}()

	go func() {
		in <- "more ./mysql.txt"
		fmt.Println("in")
	}()

	go func() {
		for {
			select {
			case v, ok := <-out:
				fmt.Println("ok", ok)
				fmt.Println("v", v)

			}
		}

	}()
	time.Sleep(100 * time.Second)
}

// RunTerminal 执行带交互的命令
func (c *Cli) RunTerminal(shell string) error {
	session, err := c.newSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return c.runTerminalSession(session, shell)
}

//-- -------------------------------------------------
//--> @Description  此方法解决ssh 出现more 获取不全问题
//--> @Param
//--> @return
//-- ----------------------------
func sshGetAll(session *ssh.Session, shell string) {
	stdoutPipe, err := session.StdoutPipe()
	fmt.Println(err)
	stderrPipe, err := session.StderrPipe()
	fmt.Println(err)
	outputReader := io.MultiReader(stdoutPipe, stderrPipe)
	outputScanner := bufio.NewScanner(outputReader)

	// Start the session.
	err = session.Start(shell)
	fmt.Println(err)

	// Capture the output asynchronously.
	outputLine := make(chan string)
	outputDone := make(chan bool)
	go func(scan *bufio.Scanner, line chan string, done chan bool) {
		defer close(line)
		defer close(done)
		for scan.Scan() {
			line <- scan.Text()
		}
		done <- true
	}(outputScanner, outputLine, outputDone)

	// Use a custom wait.
	outputBuf := ""
	running := true
	for running {
		select {
		case <-outputDone:
			running = false
		case line := <-outputLine:
			outputBuf += line + "\r\n"
		}
	}
	session.Close()
	fmt.Print(outputBuf)
	// Output the data.

}

// runTerminalSession 执行带交互的命令
func (c *Cli) runTerminalSession(session *ssh.Session, shell string) error {
	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)

	if err != nil {
		panic(err)
	}
	defer terminal.Restore(fd, oldState)

	session.Stdout = os.Stdout
	session.Stderr = os.Stdin
	session.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)

	if err != nil {
		panic(err)
	}
	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          100,    // enable echoing
		ssh.TTY_OP_ISPEED: 144000, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 144000, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
		return err
	}

	session.Run(shell)

	return nil
}

// EnterTerminal 完全进入终端
func (c Cli) EnterTerminal() error {
	session, err := c.newSession()
	if err != nil {
		return err
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer terminal.Restore(fd, oldState)

	session.Stdout = os.Stdout
	session.Stderr = os.Stdin
	session.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	err = session.RequestPty("xterm-256color", termHeight, termWidth, modes)
	if err != nil {
		return err
	}

	err = session.Shell()
	if err != nil {
		return err
	}

	return session.Wait()
}

// Enter 完全进入终端
func (c Cli) Enter(w io.Writer, r io.Reader) error {
	session, err := c.newSession()
	if err != nil {
		return err
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	// oldState, err := terminal.MakeRaw(fd)
	// if err != nil {
	// 	return err
	// }
	// defer terminal.Restore(fd, oldState)

	session.Stdout = w
	session.Stderr = os.Stdin
	session.Stdin = r

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	err = session.RequestPty("xterm-256color", termHeight, termWidth, modes)
	if err != nil {
		return err
	}

	err = session.Shell()
	if err != nil {
		return err
	}

	return session.Wait()
}

// 连接
func (c *Cli) connect() error {
	config := ssh.ClientConfig{
		User: c.Username,
		Auth: []ssh.AuthMethod{ssh.Password(c.Password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", c.IP, c.Port)
	sshClient, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		return err
	}
	c.client = sshClient
	return nil
}

// newSession new session
func (c Cli) newSession() (*ssh.Session, error) {
	if c.client == nil {
		if err := c.connect(); err != nil {
			return nil, err
		}
	}
	session, err := c.client.NewSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}
