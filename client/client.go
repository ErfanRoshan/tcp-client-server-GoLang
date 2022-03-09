package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func rFS(c net.Conn) {
	for {
		// fmt.Print("rfS")
		message, _ := bufio.NewReader(c).ReadString('\n')
		fmt.Print("->: " + message)
	}
}
func wTS(c net.Conn) {
	for {
		// fmt.Print("wTS")
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
	go rFS(c)
	go wTS(c)

	time.Sleep(5000000 * time.Second)
	// for {
	// 	reader := bufio.NewReader(os.Stdin)
	// 	fmt.Print(">> ")
	// 	text, _ := reader.ReadString('\n')
	// 	fmt.Fprintf(c, text+"\n")

	// 	message, _ = bufio.NewReader(c).ReadString('\n')
	// 	fmt.Print("->: " + message)
	// 	if strings.TrimSpace(string(text)) == "STOP" {
	// 		fmt.Println("TCP client exiting...")
	// 		return
	// 	}

	// }
}
