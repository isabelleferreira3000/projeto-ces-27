package main

// por enqnto só tem o código da tarefa 2 do lab 1
import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	_ "strconv"
	"time"
)

// global variables
var err error

var myPort string
var myId int

var nPorts int
var isCandidate bool
var numberSentMessages int
var SendersConn []*net.UDPConn
var ReceiversConn *net.UDPConn

var ch = make(chan int)

type MessageStruct struct {
	Id int
	Type string
}

// auxiliary functions
func readFileParameters(filepath string) {
	file, err := os.Open(filepath)
	CheckError(err)

	defer file.Close()

	reader := bufio.NewReader(file)

	// reading number of ports
	line, _, err := reader.ReadLine()
	CheckError(err)
	nPorts, err = strconv.Atoi(string(line))
}

func max(x int, y int) int {
	if x >= y {
		return x
	} else {
		return y
	}
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Erro: ", err)
		os.Exit(0)
	}
}

func doReceiverJob() {
	buf := make([]byte, 1024)

	n, _, err := ReceiversConn.ReadFromUDP(buf)
	CheckError(err)

	var otherLogicalClock MessageStruct
	err = json.Unmarshal(buf[:n], &otherLogicalClock)
	CheckError(err)

	fmt.Println("Received", otherLogicalClock)
	myId = logicalClock.Id
	myClocks := logicalClock.Clocks
	otherProcessClocks := otherLogicalClock.Clocks

	// updating clocks
	logicalClock.Clocks[myId-1]++
	for i := 0; i < nPorts; i++ {
		logicalClock.Clocks[i] = max(otherProcessClocks[i], myClocks[i])
	}

	fmt.Println("logicalClock atualizado:", logicalClock)
}

func doSenderJob(otherProcessID int, msgType string) {
	otherProcess := otherProcessID - 1

	var msg MessageStruct
	msg.Type = msgType
	msg.Id = 

	jsonRequest, err := json.Marshal(MessageStruct)
	CheckError(err)

	numberSentMessages ++
	_, err = SendersConn[otherProcess].Write(jsonRequest)
	CheckError(err)

	time.Sleep(time.Second * 1)
}

func initConnections() {
	numberSentMessages = 0

	// getting my Id
	myId, err := strconv.Atoi(os.Args[1])
	CheckError(err)

	// getting if is candidate
	isCandidateAux, err := strconv.Atoi(os.Args[2])
	CheckError(err)
	if isCandidateAux == 1 {
		isCandidate = true
	} else {
		isCandidate = false
	}
	fmt.Print("isCandidate: ", isCandidate, "\n")

	// getting my port
	myPort = ":" + strconv.Itoa(10000+myId)

	// Server
	ReceiverAddr, err := net.ResolveUDPAddr("udp", myPort)
	CheckError(err)
	ReceiversConn, err = net.ListenUDP("udp", ReceiverAddr)
	CheckError(err)

	// Clients
	for i := 0; i < nPorts; i++ {
		// getting each port
		aPort := ":" + strconv.Itoa(10001+i)
		fmt.Printf("aPort: %s\n", aPort)
		//aPort := os.Args[i+2]

		ServerAddr, err := net.ResolveUDPAddr("udp","127.0.0.1" + aPort)
		CheckError(err)

		LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
		CheckError(err)

		auxConn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
		SendersConn = append(SendersConn, auxConn)
		CheckError(err)
	}
}

func main() {
	readFileParameters("params.txt")
	fmt.Printf("nPorts: %d\n", nPorts)

	initConnections()

	defer ReceiversConn.Close()
	for i := 0; i < nPorts; i++ {
		defer SendersConn[i].Close()
	}

	for {
		// Server
		go doReceiverJob()

		select {
		case processID, valid := <-ch:
			if valid {
				// Clients
				if processID == myId {
					fmt.Printf("logicalClock atualizado: %d \n", logicalClock)
				} else {
					fmt.Printf("logicalClock enviado: %d \n", logicalClock)
					go doSenderJob(processID)
				}

			} else {
				fmt.Println("Channel closed!")
			}
		default:
			time.Sleep(time.Second * 1)
		}
	}
}
