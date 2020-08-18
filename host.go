package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

//Host 表示测试用主机
type Host struct {
	peersBeginID  int
	peersEndID    int
	name          string
	currentPeerID int
	//peerID自增器
	IDChan    chan int
	serverURL string
	peers map[*Peer]bool
	reg chan *Peer
	unreg chan *Peer

}

func (h *Host) init(name string, beginID, endID int, serverURL string) {
	h.name = name
	h.peersBeginID = beginID
	h.peersEndID = endID
	h.currentPeerID = beginID
	h.IDChan = make(chan int, endID-beginID+1)
	h.serverURL = serverURL
	h.peers=make(map[*Peer]bool)
	h.reg=make(chan *Peer,4)

	for i:=beginID;i<=endID;i++{
		h.IDChan<-i
	}
/*尝试内存换CPU
	go func(hst *Host) {
		for {
			select {
			case hst.IDChan <- hst.currentPeerID:
				if hst.currentPeerID < hst.peersEndID {
					hst.currentPeerID++
				} else {
					return
				}
			default:

			}
		}
	}(h)
*/
}
func (h *Host) nextPeerID() int {
	for {
		select {
		case ID := <-h.IDChan:
			if ID== h.peersEndID {
				close(h.IDChan)
			}
			return ID
		// default:
		// 	if h.currentPeerID >= h.peersEndID {
		// 		return -1
		// 	}
		}
	}
}

func (h *Host) getCourseID() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(10) + 1
}

func (h *Host) getPartID() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(7) + 1

}

func (h *Host)run(){
	for {
		select {
		case peer :=<-h.reg:
			h.peers[peer]=true
		case peer :=<-h.unreg:
			if _,ok :=h.peers[peer]; ok{
				delete (h.peers,peer)
				close(peer.send)
			}
		//default:

		}
	}
}


func (h *Host) createOnePeer() *Peer {
	p := new(Peer)
	c, _, err := websocket.DefaultDialer.Dial(h.serverURL, nil)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	p.userID = h.nextPeerID()
	if p.userID == -1 {
		p = nil
		return nil
	}
	
	p.conn = c
	p.courseID = h.getCourseID()
	p.partID = h.getPartID()
	p.send = make(chan []byte, 256)
	counter++
	fmt.Println("创建成功",counter,":",p.userID,",",p.courseID,",",p.partID,",",p.process)
	return p
}

func (h *Host) createPeers() {
	for i:=h.peersBeginID;i<=h.peersEndID;i++ {
		//连续三次创建不成功，退出
		var p *Peer
		if p = h.createOnePeer(); p == nil{
			if p = h.createOnePeer(); p == nil{
				if p = h.createOnePeer(); p == nil{
					return
				}
			}
		}
		//登记
		h.reg<-p
		h.actPeer(p)
	
	}
}
func (h *Host) createPeersOnly() {
	for i:=h.peersBeginID;i<=h.peersEndID;i++ {
		//连续三次创建不成功，退出
		var p *Peer
		if p = h.createOnePeer(); p == nil{
			if p = h.createOnePeer(); p == nil{
				if p = h.createOnePeer(); p == nil{
					return
				}
			}
		}
		
		//登记
		// h.reg<-p
		// h.actPeer(p)
	
	}
}


func(h *Host) actAllPeers(){
	counter:=0
	for p := range h.peers {
		h.actPeer(p)
		counter++
		fmt.Printf("%d开始模拟播放，每%d秒发送播放进度.共有%d模拟中\n",p.userID,processPeriod/time.Second,counter)
	}
}

func(h *Host)deactAllPeers(){
	//程序可能有并发问题
	for p := range h.peers {
					close(p.send)
					// if p.conn !=nil {
					// 	p.conn.WriteMessage(websocket.CloseMessage, []byte{})
					// 	p.conn.Close()
					// }
					delete(h.peers,p)
				

			}
}

func(h *Host)actPeer(p *Peer){
	//模拟一个用户
	go p.playingVideo()//模拟播放视频
	go p.writePump() //向服务器发送本机的消息
	go p.readPump()  //读取服务端送来的消息
}

