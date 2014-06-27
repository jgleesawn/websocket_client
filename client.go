package main

import (
	"fmt"
	"net/http"
	"os"
	"bufio"
	"github.com/gorilla/websocket"
	"ECC_Conn"

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
		fmt.Println(err)
		return
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
	dh_conn := new(ECC_Conn.ECC_Conn)
	dh_conn.Connect(conn)
	fmt.Println("Outside diffie.")
	buf := make([]byte,dh_conn.PacketSize)
	n,_ := dh_conn.Read(buf)
	fmt.Println(string(buf[:n]))

	go process(dh_conn)

	reader := bufio.NewReader(os.Stdin)
	for {
		var q *Quest
		q = &Quest{0,"Test","Testing Function","Test",false,0,[]int{0},[]string{""}}
		addQuest(dh_conn,q)
		var u *User
		u = &User{"testing_username","test","test",100,[]int{0},[]string{""}}
		addUser(dh_conn,u)
		q.Name = "update"
		updateQuest(dh_conn,q)
		u.Firstname = "update"
		updateUser(dh_conn,u)

		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		str := []byte(line[0:len(line)-1])
		err = conn.WriteMessage(websocket.TextMessage,str)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
func addQuest(conn *ECC_Conn.ECC_Conn, q *Quest){
	data,err := json.Marshal(q)
	sep := []string{"add Quest",string(data)}
	jn := []byte(strings.Join(sep,";"))
	_,err = conn.Write(jn)
	if err != nil {
		fmt.Println("Sending Message failed.")
	}
}
func addUser(conn *ECC_Conn.ECC_Conn, u *User){
	data,err := json.Marshal(u)
	sep := []string{"add User",string(data)}
	jn := []byte(strings.Join(sep,";"))
	_,err = conn.Write(jn)
	if err != nil {
		fmt.Println("Sending Message failed.")
	}
}
func updateQuest(conn *ECC_Conn.ECC_Conn, q *Quest){
	data,err := json.Marshal(q)
	sep := []string{"update Quest",string(data)}
	jn := []byte(strings.Join(sep,";"))
	_,err = conn.Write(jn)
	if err != nil {
		fmt.Println("Sending Message failed.")
	}
}
func updateUser(conn *ECC_Conn.ECC_Conn, u *User){
	data,err := json.Marshal(u)
	sep := []string{"update User",string(data)}
	jn := []byte(strings.Join(sep,";"))
	_,err = conn.Write(jn)
	if err != nil {
		fmt.Println("Sending Message failed.")
	}
}

func process(conn *ECC_Conn.ECC_Conn) {
	data := make([]byte,conn.PacketSize)
	for {
		n,err := conn.Read(data)
		if n > 0 {
			fmt.Println(string(data[:n]))
		}
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
