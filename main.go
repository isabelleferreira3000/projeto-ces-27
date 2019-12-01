package main

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
var timer1 *time.Timer

var err error

var myPort string
var myId int

var nPorts int
var isCandidate bool
var isRunningMyElection bool
var coordinatorId int
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

func CheckError(err error) {
	if err != nil {
		fmt.Println("Erro: ", err)
		os.Exit(0)
	}
}

func doReceiverJob() {
	for coordinatorId == -1 {
		buf := make([]byte, 1024)

		n, _, err := ReceiversConn.ReadFromUDP(buf)
		if err != nil {
			break
		}
		CheckError(err)

		var msg MessageStruct
		err = json.Unmarshal(buf[:n], &msg)
		CheckError(err)

		fmt.Println("Received msg.type =", msg.Type, "from id =", msg.Id)

		if msg.Type == "ELECTION" {
			if msg.Id < myId {
				doSenderJob(msg.Id, "OK")
				startElection()
			}
		} else if msg.Type == "OK" {
			isRunningMyElection = false
		} else if msg.Type == "COORDINATOR" {
			coordinatorId = msg.Id
			break
		}
	}
}

func doSenderJob(otherProcessID int, msgType string) {
	otherProcess := otherProcessID - 1

	var msg MessageStruct
	msg.Type = msgType
	msg.Id = myId

	jsonRequest, err := json.Marshal(msg)
	CheckError(err)

	numberSentMessages ++
	_, err = SendersConn[otherProcess].Write(jsonRequest)
	CheckError(err)

	fmt.Println("Sending msg.type =", msg.Type, "from id =", msg.Id)

	time.Sleep(time.Second * 1)
}

func initConnections() {
	numberSentMessages = 0
	coordinatorId = -1

	// getting my Id
	myId, err = strconv.Atoi(os.Args[1])
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

		ServerAddr, err := net.ResolveUDPAddr("udp","127.0.0.1" + aPort)
		CheckError(err)

		LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
		CheckError(err)

		auxConn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
		SendersConn = append(SendersConn, auxConn)
		CheckError(err)
	}
}

func sendCoordinatorMsgs() {
	for otherProcessId := 1; otherProcessId < nPorts; otherProcessId++ {
		if otherProcessId != myId {
			doSenderJob(otherProcessId, "COORDINATOR")
		}
	}
}

func timerTracker(timer *time.Timer) {
	<-timer.C
	fmt.Println("Timer expired")
	if isRunningMyElection {
		coordinatorId = myId
		sendCoordinatorMsgs()
	}
}

func startElection() {
	fmt.Printf("Starting election\n")

	if !isRunningMyElection {
		isRunningMyElection = true
		for otherProcessId := myId + 1; otherProcessId < nPorts+1; otherProcessId++ {
			doSenderJob(otherProcessId, "ELECTION")
		}
		timer1 = time.NewTimer(2 * time.Second)
		go timerTracker(timer1)
	}
}

func printFinalResults() {
	fmt.Printf("COORDINATOR ID = %d\n", coordinatorId)
	fmt.Printf("END\n")

	f, err := os.OpenFile("results.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	msgToPrint := "Proccess " + strconv.Itoa(myId) + ": coodinator = " + strconv.Itoa(coordinatorId) + 
		", with " + strconv.Itoa(numberSentMessages) + " messages sent\n"

	if _, err := f.WriteString(msgToPrint); err != nil {
		fmt.Println(err)
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

	if isCandidate {
		startElection()
	}
	go doReceiverJob()

	for coordinatorId == -1 {}

	printFinalResults()

	time.Sleep(1 * time.Second)
}
