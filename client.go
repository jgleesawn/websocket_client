package main

import (
	"fmt"
	"net/http"
	"os"
	"bufio"
	"github.com/gorilla/websocket"

//	"io"
	"encoding/json"
	"strings"
)

type Quest struct {
	Questid		int
	Name		string
	Description	string
	Category	string
	Recurring	bool
	Xpvalue		int
	Requiredquests	[]int
	Attributes	[]string
}
func (quest Quest) New(qid int, name string, desc string, cat string, rec bool, xpval int, reqquests []int, attr []string) {
	quest.Questid = qid
	quest.Name= name
	quest.Description = desc
	quest.Category = cat
	quest.Recurring = rec
	quest.Xpvalue = xpval
	quest.Requiredquests = reqquests
	quest.Attributes = attr
}

type User struct {
	Username	string	`db:"Username"`
	Firstname	string	`db:"Firstname"`
	Lastname	string	`db:"Lastname"`
	Xp		int	`db:"Xp"`
	Completedquests	[]int	`db:"Completedquests"`
	Attributes	[]string `db:"Attributes"`
}
func (user User) New(u string,f string, l string, a []string) {
	user.Username = u
	user.Firstname = f
	user.Lastname = l
	user.Xp = 0
	user.Completedquests = make([]int,0)
	user.Attributes = a
}


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
	resp.Header.Add("Upgrade","websocket")
	resp.Header.Add("Connection","Upgrade")
	resp.Header.Add("Host",url)
	resp.Header.Add("Origin","http://"+url+port)
	resp.Header.Add("Sec-WebSocket-Protocol","chat, superchat")
	resp.Header.Add("Sec-WebSocket-Extensions","permessage-deflate; client_max_window_bits, x-webkit-deflate-frame")
	resp.Header.Add("Sec-WebSocket-Key","z1VdBz6K3WZTV3rMw2QUFw==")
	resp.Header.Add("Sec-WebSocket-Version","13")
	resp.Header.Add("Cache-Control","no-cache")
	resp.Header.Add("Pragma","no-cache")
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
		var q *Quest
		q = &Quest{0,"Test","Testing Function","Test",false,0,[]int{0},[]string{}}
		addQuest(conn,q)
		var u *User
		u = &User{"test","test","test",100,[]int{0},[]string{""}}
		addUser(conn,u)

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
func addQuest(conn *websocket.Conn, q *Quest){
	data,err := json.Marshal(q)
	mt := websocket.TextMessage
	sep := []string{"add Quest",string(data)}
	err = conn.WriteMessage(mt,[]byte(strings.Join(sep,";")))
	if err != nil {
		fmt.Println("Sending Message failed.")
	}
}
func addUser(conn *websocket.Conn, u *User){
	data,err := json.Marshal(u)
	mt := websocket.TextMessage
	sep := []string{"add User",string(data)}
	err = conn.WriteMessage(mt,[]byte(strings.Join(sep,";")))
	if err != nil {
		fmt.Println("Sending Message failed.")
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
