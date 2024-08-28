package main

import (
	_ "embed"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

//go:embed about.md
var about []byte

var tags []string

const description = "I do nothing :/"

// Render generates smth.
func Render(msg []byte, conn net.Conn, n int) error {

	_, err := conn.Write(about)
	return err
}

func main() {
	tags = append(tags, "mock", "test", "tag1")
	// Create a Unix domain socket and listen for incoming connections.
	socket, err := net.Listen("unix", "./plugin.sock")
	if err != nil {
		log.Fatal(err)
	}

	// Cleanup the sockfile.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove("./plugin.sock")
		os.Exit(1)
	}()

	for {
		// Accept an incoming connection.
		conn, err := socket.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Handle the connection in a separate goroutine.
		go func(conn net.Conn) {
			defer conn.Close()
			// Create a buffer for incoming data.
			buf := make([]byte, 4096)

			// Read data from the connection.
			n, err := conn.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(buf[:n])

			// Generate smth.
			if string(buf[:n]) == "doc" {
				err = Render(buf, conn, n)
				if err != nil {
					log.Fatal(err)
				}
			}

			if string(buf[:n]) == "tag" {
				for i := range tags {
					_, err := conn.Write([]byte(tags[i]))
					if err != nil {
						log.Fatal(err)
					}
				}
			}

			if string(buf[:n]) == "desc" {
				_, err := conn.Write([]byte(description))
				if err != nil {
					log.Fatal(err)
				}

			}
		}(conn)
	}
}
