package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)


var (
	newline = []byte{'\n'}
	space   = []byte{' '}
	counter=0
	

	
)

var (
	//测试需要修改参数之处
	//程序运行后按ctrl+C退出
	//发送播放进度间隔
	processPeriod=10*time.Second
	//起始和结束userID
    beginUserID=1
	endUserID=6000
	wsAddr="ws://112.74.99.240:8080/ws"
	//wsAddr="ws://127.0.0.1:8080/ws"
	hostName="罗勇胜"
)






func main() {
	if !parseCommandLine(){
		return
	}
	h := new(Host)
	h.init(hostName, beginUserID, endUserID,wsAddr)
	go h.run() //消息总线
	go h.createPeers() //创建用户，且并发执行
	quit:=make(chan os.Signal)
	signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)
	<-quit
	
	fmt.Println("shutting down....")
	h.deactAllPeers()	


	time.Sleep(1*time.Second)
	
	
}

func mainnn() {
	if !parseCommandLine(){
		return
	}
	h := new(Host)
	h.init(hostName, beginUserID, endUserID,wsAddr)
	go h.run() //消息总线
	h.createPeersOnly() //创建用户
	h.actAllPeers()
	quit:=make(chan os.Signal)
	signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)
	<-quit
	
	fmt.Println("shutting down....")
	h.deactAllPeers()	
	

	time.Sleep(1*time.Second)
	
	
}

func printNextID(h *Host){
	for {
		if ID:=h.nextPeerID();ID!=-1 {
			fmt.Println(ID)
		}else{
			return
		}
	}
	
}

func parseCommandLine()bool{
	
	if len(os.Args)==2 && os.Args[1]=="vscodetest" {
		return true
	}
	if len(os.Args) !=6 {
		fmt.Println("需要6个参数。")
		printTips()
		return false
	}	
	
	hostName=os.Args[1]
	ID,err:=strconv.Atoi(os.Args[2])
	if err!=nil{
		fmt.Println("第二个参数必须是数字。")
		printTips()
		return false
	}
	beginUserID=ID
	ID,err=strconv.Atoi(os.Args[3]);
	if  err!=nil{
		fmt.Println("第三个参数必须是数字。")
		printTips()
		return false
	}
	endUserID=ID
	var t int64
	t,err=strconv.ParseInt(os.Args[4],10,64);
	
	if  err!=nil{
		fmt.Println("第四个参数必须是数字。")
		printTips()
		return false
	}
	processPeriod=time.Duration(t)*time.Second
	wsAddr=os.Args[5]
	return true
}
func printTips()  {
	fmt.Println("按ctrl+c退出")
	fmt.Println("用法：goTester 测试者姓名  模拟用户起始ID 模拟用户截止ID  模拟用户发送数据间隔秒数 完整服务器地址")
	fmt.Println("模拟用户起始ID 模拟用户截止ID:起始ID可按指定来，截止ID要小于或等于指定的截止ID")
	fmt.Println("例：goTester lyshsd 1 1000 1 ws://127.0.0.1:8080/ws")
}