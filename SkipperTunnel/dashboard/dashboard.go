package dashboard

import (
	"SkipperTunnel/constants"
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

func NewDashboardServer() *http.Server {
	tmux := http.NewServeMux()
	tmux.HandleFunc("/", DashboardHandler)
	return &http.Server{
		Handler: tmux,
	}
}

func StartDashboard(srv *http.Server, ctx context.Context, errChan chan error) {
	port, err := checkPorts()
	if err != nil {
		return
	}
	srv.Addr = ":" + port

	defer srv.Shutdown(ctx)

	go func() {
		srv.ListenAndServe()
		if err != nil {
			errChan<- err
			return
		}

	}()	
	<- ctx.Done()
	fmt.Println("ya acabamos el start dashboard")
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HPLA vamos a tener nuestro coso"))
}

func checkPorts() (string, error) {
	var port int
	for i := 3000; i < 9500; i += 24 {
		_, err := net.DialTimeout("tcp", "localhost:"+strconv.Itoa(i), time.Millisecond*500)

		if err != nil {
			// fmt.Println("port found", i)
			port = i
			break
		}
		// fmt.Println("port didint work, searching for other", i)
	}

	fmt.Print("You can see the request dashboard on: ")
	constants.PrintWithColor(constants.Cyan, "localhost:"+strconv.Itoa(port))
	return strconv.Itoa(port), nil
}
