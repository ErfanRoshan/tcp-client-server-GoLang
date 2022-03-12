package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const address = "users.txt"

var userList []User
var connS []Conn

type User struct {
	PhoneNumber string
	Name        string
}

type Conn struct {
	Connection  *(net.Conn)
	PhoneNumber string
}

func SaveJson(fileAddress string, users *[]User) {
	f, err := os.Create(fileAddress)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	data, _ := json.Marshal(users)
	_, err2 := f.WriteString(string(data) + "\n")
	if err2 != nil {
		log.Fatal(err2)
	}
}

func LoadJson(fileAddress string) []User {
	file, err := os.Open(fileAddress)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	var texts []string

	for scanner.Scan() {
		texts = append(texts, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var users []User
	json.Unmarshal([]byte(texts[0]), &users)
	return users
}

func removeC(s []Conn, i int) []Conn {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func appendC(c net.Conn, pNum string) Conn {
	var conn Conn
	conn.Connection = &c
	conn.PhoneNumber = pNum
	connS = append(connS, conn)
	return conn
}

func handleConnection(c net.Conn) {
	phoneData, err1 := bufio.NewReader(c).ReadString('\n')
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	pNum := strings.TrimSpace(string(phoneData))
	name := ""

	for _, user := range userList {
		if user.PhoneNumber == pNum {
			name = user.Name
			break
		}
	}

	if name == "" {
		c.Write([]byte(string("username\n")))
		nameData, err1 := bufio.NewReader(c).ReadString('\n')
		if err1 != nil {
			fmt.Println(err1)
			return
		}
		name = strings.TrimSpace(string(nameData))
		var user User
		user.Name = name
		user.PhoneNumber = pNum
		userList = append(userList, user)
		SaveJson(address, &userList)
	}

	c.Write([]byte(string(name + " : connected\n")))

	fmt.Println(name + " : connected")
	conn := appendC(c, pNum)

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		temp := strings.TrimSpace(string(netData))
		if temp == "STOP" {
			idx := -1
			for i := 0; i < len(connS); i++ {
				if connS[i].Connection == conn.Connection {
					idx = i
					break
				}
			}
			connS = removeC(connS, idx)
			fmt.Println(name + " : disconnected")
			c.Write([]byte(string("STOP\n")))
			break
		}
		for _, connection := range connS {
			if connection.PhoneNumber != conn.PhoneNumber {
				cOthers := *(connection.Connection)
				cOthers.Write([]byte(string(name + " : " + temp + "\n")))
			}
		}
	}
	c.Close()
}

func main() {
	PORT := ":1313"
	_, error := os.Stat(address)
	if os.IsNotExist(error) {
		f, err := os.Create(address)

		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()
	} else {
		userList = LoadJson(address)
	}
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}
