package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

func rFS(c net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		message, _ := bufio.NewReader(c).ReadString('\n')
		if strings.TrimSpace(string(message)) == "STOP" {
			return
		} else {
			fmt.Print("->: " + message)
		}
	}
}
func wTS(c net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(c, text+"\n")
		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	CONNECT := "localhost:1313"
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print("Enter phone number:\n>> ")
	var pNum string
	var name string
	fmt.Scan(&pNum)
	fmt.Fprintf(c, pNum+"\n")
	message, _ := bufio.NewReader(c).ReadString('\n')
	if strings.TrimSpace(string(message)) == "username" {
		fmt.Print("Enter username:\n>> ")
		fmt.Scan(&name)
		fmt.Fprintf(c, name+"\n")
		message, _ = bufio.NewReader(c).ReadString('\n')
		fmt.Print("->: " + message)
	} else {
		fmt.Print("->: " + message)
	}
	go rFS(c, &wg)
	go wTS(c, &wg)
	wg.Wait()
}
