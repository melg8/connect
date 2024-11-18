package connect

import (
	"encoding/hex"
	"fmt"
	"net"
	"testing"
)

func InitPacketData() []byte {
	packet := "9f0000d1bbdc3821c600000311866bd523c9cadaca10db1bd8782447d22d7a6fbb5ba84e933b54fc6ecb1e8d32b5f2a7a582776052374c722bae5dd2ee0805da8dcb8328fd3ba4e2a57a52ca5f99fa9c1f5d21e2877a9d3266e27c9c34d6c75f556138bcb985f3019bb62644d6fc93d14d4c54091e5d1826f57db183ad3944ee2f04d4710cd8aabeb15df14e95dd29fc9cc37720b6ad97f7e0bd0700000000"
	data, err := hex.DecodeString(packet)
	if err != nil {
		panic(err)
	}
	return data
}

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
			_, err = conn.Write(InitPacketData())
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

	// t.Run("explicit shutdown", func(t *testing.T) {
	// 	ln, err := net.Listen("tcp", "localhost:2112")
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	defer ln.Close()
	// 	go func() {
	// 		conn, err := ln.Accept()
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		defer conn.Close()
	// 	}()
	// 	multiAuth := StartBotsAt(ln.Addr().String(), 1)
	// 	multiAuth.Stop()
	// })
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
