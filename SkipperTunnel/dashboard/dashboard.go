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

type DashboardServer struct{
	srv *http.Server
	requestChan chan *http.Request
}

func NewDashboardServer(requestChan chan *http.Request) *DashboardServer {
	tmux := http.NewServeMux()
	tmux.HandleFunc("/", DashboardHandler)

	return &DashboardServer{
		srv: &http.Server{
			Handler: tmux,
		},
		requestChan: requestChan,
	}
}

func (d *DashboardServer)StartDashboard(ctx context.Context, errChan chan error) {
	port, err := checkPorts()
	if err != nil {
		return
	}
	d.srv.Addr = ":" + port

	defer d.srv.Shutdown(ctx)

	go func() {
		d.srv.ListenAndServe()
		if err != nil {
			errChan<- err
			return
		}

	}()	
	<- ctx.Done()
	fmt.Println("ya acabamos el start dashboard")
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/dashboard.html")
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
	constants.PrintWithColor(constants.Cyan, "http://localhost:"+strconv.Itoa(port))
	return strconv.Itoa(port), nil
}
