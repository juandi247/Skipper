package TcpUserClient

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
)

type readResult struct {
	content []byte
	err     error
}

// thats to receive data onn a buffer and readed
func HandleReceive(conn net.Conn, ch chan []byte, ctx context.Context, wg *sync.WaitGroup) {

	readChannel := make(chan readResult)
	for {
		go func() {
			sizeBuf := make([]byte, 4)
			_, err := io.ReadFull(conn, sizeBuf)
			if err != nil {
				// fmt.Print("An error ocurred with the proxy connection", err)
				readChannel <- readResult{err: err}
				return
			}

			length := binary.BigEndian.Uint32(sizeBuf)
			// fmt.Println("TENEMOS UN LENGTH THE REQUEST DE", length)
			msgBuf := make([]byte, length)
			// Si se reciben datos, envíalos al canal
			_, err = io.ReadFull(conn, msgBuf)
			if err != nil {
				// fmt.Println("Error leyendo mensaje:", err)
				readChannel <- readResult{err: err}
			}

			fmt.Printf("Received: %s\n", msgBuf)
			readChannel <- readResult{content: msgBuf}
		}()

		select {
		case <-ctx.Done():
			wg.Done()
			fmt.Println("gtting out of the receivng tcp data")
			return
		case r := <-readChannel:
			if r.err != nil {
				wg.Done()
				return
			}
			ch <- r.content
		}
	}
}

// select{
// 		case <-ctx.Done():
// 			wg.Done()
// 			fmt.Println("gtting out of the receivng tcp data")
// 			return
// 		default:

func HandleSendToTCP(response []byte, conn net.Conn) {
	// fmt.Println("VOY A ENVIAR", response)
	lenght := uint32(len(response))

	buf := new(bytes.Buffer)

	binary.Write(buf, binary.BigEndian, lenght)
	buf.Write(response)

	// todo: check becasue we are reconverting the resopne to bytes, and the resopnse was already in bytes
	// data := []byte(response)
	_, err := conn.Write(buf.Bytes())

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("sendet Hello server message")
}
