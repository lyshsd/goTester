package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

//Peer ---------------peer---------------------
type Peer struct{
	userID int
	courseID int
	partID int
	process int
	conn *websocket.Conn 
	send chan []byte
}
//todo: 应全面模拟实际场景
// type Video struct {
// 	VideoID 

// }


func(p *Peer)playingVideo(){
	ticker := time.NewTicker(processPeriod)
	defer func() {
		ticker.Stop()
		fmt.Println(p.userID,"停止播放.....")
	}()

	for {
		select {
		
			//每个processPeriod发送一次消息到通道send
		case <-ticker.C:
			p.send<-p.getProcessMessage()
			if p.process>=3600 {
				fmt.Println(p.userID,"用户1小时视频播放测试完成")
				return
			}
		
		case <-shutdownPlaying:
				return
			
		    
			
		// default:
		// 	//超过1小时退出
		// 	if p.process>3600 {
		// 		return
		// 	}
		}
	}	
}

func(p *Peer)getProcessMessage() []byte{
	p.process+=10
	message:=fmt.Sprintf("{\"UserIdCourseId\":\"%d_%d\",\"partid\":\"%d\",\"currentTime\":%d}",p.userID,p.courseID,p.partID,p.process)
	return []byte(message)
}


func(p *Peer)readPump(){
	defer p.conn.Close()
	for {
		_, message, err := p.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		fmt.Println("收到消息：",string(message))
	}
}

func (p *Peer) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		p.conn.Close()
		fmt.Println(p.userID,",writePump退出了")
	}()

	for {
		select {
		case message, ok := <-p.send:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// channel was closed.
				p.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			
			w, err := p.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				fmt.Printf("%d NextWriter  Error:%s\n",p.userID, err.Error())
				
				return
			}
			w.Write(message)
			fmt.Printf("用户%d发送消息:%s\n",p.userID,string(message))

			// Add queued chat messages to the current websocket message.
			n := len(p.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-p.send)
			}

			if err := w.Close(); err != nil {
				fmt.Printf("%d w.Close Error:%s\n",p.userID, err.Error())
				return

			}
			
		case <-ticker.C:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Printf("%d ticker.C Error:%s\n",p.userID, err.Error())
				
				return
			}
		// case <-shutdown_writePump:
		// 	p.conn.WriteMessage(websocket.CloseMessage, []byte{})
		// 		return

		}
	}
}

//IntToBytes convert a int to []byte
// func IntToBytes(n int) []byte {
//   x := int32(n)
//   bytesBuffer := bytes.NewBuffer([]byte{})
//   binary.Write(bytesBuffer, binary.BigEndian, x)
//   return bytesBuffer.Bytes()
// }