package tcpserver

type TcpMessage struct {
	Target string `json:"target"`
	Data   []byte `json:"data"`  // Este será el JSON con la solicitud HTTP
}