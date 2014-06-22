package main

import (
	"fmt"
	"net/http"
	"os"
	"bufio"
	"github.com/gorilla/websocket"
)

func main() {
	var port string
	var url string
	if len(os.Args) < 3 {
		//default port is port 80
		//using port = ":80" in the origin can break the handshake
		//give Bad Origin.
		port = ""
	} else {
		//add : here for simple appending of variable
		//can append empty string for default port now.
		port = ":"+os.Args[2]
	}
	if len(os.Args) < 2 {
		url = "onelyfe.herokuapp.com"
	} else {
		url = os.Args[1]
	}

//Look into WebSocket Keys, not sure using the same one every time is good.
	resp, err := http.Get("http://"+url+port)
	if err != nil {
		panic(err)
	}
	for i := range resp.Header {
		fmt.Println(i+":"+resp.Header[i][0])
	}
	fmt.Println()
	for i := range resp.Header {
		resp.Header.Del(i)
	}
	resp.Header.Add("Cache-Control","no-cache")
	resp.Header.Add("Connection","Upgrade")
	resp.Header.Add("Host",url)
	resp.Header.Add("Origin","http://"+url+port)
	resp.Header.Add("Pragma","no-cache")
	resp.Header.Add("Sec-WebSocket-Extensions","permessage-deflate; client_max_window_bits, x-webkit-deflate-frame")
	resp.Header.Add("Sec-WebSocket-Key","z1VdBz6K3WZTV3rMw2QUFw==")
	resp.Header.Add("Sec-WebSocket-Version","13")
	resp.Header.Add("Upgrade","websocket")
	resp.Header.Add("User-Agent","Mozilla/5.0")

	/*
	resp.Header.Add("Host",url)
	resp.Header.Add("Upgrade","websocket")
	resp.Header.Add("Connection","Upgrade")
	resp.Header.Add("Sec-WebSocket-Key","x3JJHMbDL1EzLkh9GBhXDw==")
	resp.Header.Add("Sec-WebSocket-Protocol","chat, superchat")
	resp.Header.Add("Sec-WebScoket-Version","13")
	//resp.Header.Add("Origin","http://"+url+":"+port)
	*/
	for i := range resp.Header {
		fmt.Println(i+":"+resp.Header[i][0])
	}
	fmt.Println()


//middle argument is response header.
	var DefaultDialer *websocket.Dialer
	conn, resp,err := DefaultDialer.Dial("ws://"+url+port + "/ws",resp.Header)
	for i := range resp.Header {
		fmt.Println(i+":"+resp.Header[i][0])
	}
	fmt.Println()
	if err != nil {
		panic(err)
	}

	go process(conn)

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		str := []byte(line[0:len(line)-1])
		err = conn.WriteMessage(websocket.TextMessage,str)
		if err != nil {
			panic(err)
		}
	}
}

func process(conn *websocket.Conn) {
	for {
		mt,data,err := conn.ReadMessage()
		if len(data) > 0 {
			if mt == websocket.TextMessage{
				fmt.Println(string(data))
			} else {
				fmt.Println(data)
			}
		}
		if err != nil {
			panic(err)
		}
	}
}
