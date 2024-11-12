package connect

import (
	"fmt"
	"net"
	"testing"
)

func TestConnect(t *testing.T) {
	t.Run("successful connection", func(t *testing.T) {
		ln, err := net.Listen("tcp", "localhost:2112")
		finishedChan := make(chan bool)
		if err != nil {
			t.Fatal(err)
		}
		defer ln.Close()

		go func() {
			conn, err := ln.Accept()
			if err != nil {
				panic(err)
			}
			defer conn.Close()
			_, err = conn.Write([]byte("Hello, client!"))
			if err != nil {
				panic(err)
			}
			finishedChan <- true
		}()

		fmt.Println("Address: " + ln.Addr().String())
		StartBotsAt(ln.Addr().String(), 1)
		<-finishedChan
	})

	t.Run("failed connection", func(t *testing.T) {
		StartBotsAt("localhost:2111", 1)
	})

	t.Run("invalid address", func(t *testing.T) {
		StartBotsAt("invalid-address", 1)
	})
}

// 	// Test connection failure due to timeout
// 	t.Run("timeout", func(t *testing.T) {
// 		Connect("localhost:12345")
// 		// Check that an error message was printed
// 		var buf bytes.Buffer
// 		fmt.Println = func(a ...interface{}) (n int, err error) {
// 			buf.WriteString(fmt.Sprint(a...))
// 			return
// 		}
// 		if !bytes.Contains(buf.Bytes(), []byte("Error connecting to server:")) {
// 			t.Errorf("expected error message, got '%s'", buf.String())
// 		}
// 	})

// 	// Test reading data from the server
// 	t.Run("read data", func(t *testing.T) {
// 		ln, err := net.Listen("tcp", "localhost:0")
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		defer ln.Close()

// 		go func() {
// 			conn, err := ln.Accept()
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			defer conn.Close()

// 			// Send some data to the client
// 			_, err = conn.Write([]byte("Hello, client!"))
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 		}()

// 		// Connect to the server
// 		Connect(ln.Addr().String())

// 		// Check that the client received the data
// 		var buf bytes.Buffer
// 		ShowAsHexAndAscii = func(data []byte) {
// 			buf.Write(data)
// 		}
// 		if buf.String() != "Hello, client!" {
// 			t.Errorf("expected 'Hello, client!', got '%s'", buf.String())
// 		}
// 	})

// 	// Test handling of EOF from the server
// 	t.Run("eof", func(t *testing.T) {
// 		ln, err := net.Listen("tcp", "localhost:0")
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		defer ln.Close()

// 		go func() {
// 			conn, err := ln.Accept()
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			defer conn.Close()

// 			// Close the connection to simulate EOF
// 			conn.Close()
// 		}()

// 		// Connect to the server
// 		Connect(ln.Addr().String())

// 		// Check that the client handled the EOF correctly
// 		var buf bytes.Buffer
// 		fmt.Println = func(a ...interface{}) (n int, err error) {
// 			buf.WriteString(fmt.Sprint(a...))
// 			return
// 		}
// 		if !bytes.Contains(buf.Bytes(), []byte("Server doesn't send any more data")) {
// 			t.Errorf("expected 'Server doesn't send any more data', got '%s'", buf.String())
// 		}
// 	})

// 			// Test handling of errors while reading from the server
// 			t.Run("read error", func(t *testing.T) {
// 				ln, err := net.Listen("tcp", "localhost:0")
// 				if err != nil {
// 					t.Fatal(err)
// 				}
// 				defer ln.Close()

// 				go func() {
// 					conn, err := ln.Accept()
// 					if err != nil {
// 						t.Fatal(err)
// 					}
// 					defer conn.Close()

// 					// Close the connection to simulate an error
// 					conn.Close()
// 				}()

// 				// Connect to the server
// 				Connect(ln.Addr().String())

// 				// Check that the client handled the error correctly
// 				var buf bytes.Buffer
// 				fmt.Println = func(a ...interface{}) (n int, err error) {
// 					buf.WriteString(fmt.Sprint(a...))
// 					return
// 				}
// 				if !bytes.Contains(buf.Bytes(), []byte("Error reading from connection:")) {
// 					t.Errorf("expected error message, got '%s'", buf.String())
// 				}
// 			})
// 		})
// 	}
// }
