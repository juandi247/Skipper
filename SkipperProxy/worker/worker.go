package worker

import (
	"SkipperProxy/connectionmanager"
	"SkipperProxy/gen"
	"fmt"
	"google.golang.org/protobuf/proto"
)

func StartWorker(workerId int, wpChan chan []byte, cm *connectionmanager.ConnectionManager) {
	for msg := range wpChan {
		fmt.Println("WORKER EXECUTING THE Response flwo", workerId)
		// fmt.Println("message received", string(msg))
		var response gen.Response
		err := proto.Unmarshal(msg, &response)
		if err != nil {
			fmt.Println("error parsing resopnee")
			continue
		}

		// fmt.Println("le vamos a enviarrrrr")

		cm.Mu.Lock()
		ch, exists := cm.GlobalResponseChannel[response.RequestId]

		// fmt.Println("REQUEST ID FROM RESOPNSEEEE", response.RequestID)
		if exists {
			// Enviar la respuesta al channel que espera el HTTP handler
			ch <- msg

			fmt.Println("si existio mensaje le enviamoss!!")

			// Opcional: cerrar el channel y borrarlo del mapa para limpiar
			close(ch)
			delete(cm.GlobalResponseChannel, response.RequestId)
		}
		fmt.Println("QUE PASO NO PASO NADAAAAAA")
		cm.Mu.Unlock()
	}
}
