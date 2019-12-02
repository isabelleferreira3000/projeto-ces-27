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
	"io"
	//"io/ioutil"
	//"strings"
)

// global variables
var electionTimer *time.Timer

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
var parent int
var child1 int
var child2 int

var ch = make(chan int)

type MessageStruct struct {
	Coord int
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

func ReadInts(r io.Reader) ([]int, error) {
    scanner := bufio.NewScanner(r)
    scanner.Split(bufio.ScanWords)
    var result []int
    for scanner.Scan() {
        x, err := strconv.Atoi(scanner.Text())
        if err != nil {
            return result, err
        }
        result = append(result, x)
    }
    return result, scanner.Err()
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
				doSenderJob(msg.Id, "OK", myId)
				startElection()
			}
		} else if msg.Type == "OK" {
			isRunningMyElection = false
		} else if msg.Type == "COORDINATOR" {
			coordinatorId = msg.Coord
			sendCoordinatorMsgs()
			break
		}
	}
}

func doSenderJob(otherProcessID int, msgType string, coordID int) {
	otherProcess := otherProcessID - 1

	var msg MessageStruct
	msg.Type = msgType
	msg.Id = myId
	msg.Coord = coordID

	jsonRequest, err := json.Marshal(msg)
	CheckError(err)

	numberSentMessages ++
	_, err = SendersConn[otherProcess].Write(jsonRequest)
	CheckError(err)

	fmt.Println("Sending msg.type =", msg.Type, "to id =", otherProcessID)

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
	if(child1 >= 0){
		doSenderJob(child1, "COORDINATOR", coordinatorId)
	}
	if(child2 >= 0){
		doSenderJob(child2, "COORDINATOR", coordinatorId)
	}
}

func electionTimerTracker(timer *time.Timer) {
	<-timer.C
	fmt.Println("Election Timer expired")
	if isRunningMyElection {
		coordinatorId = myId
		sendCoordinatorMsgs()
	}
}

func startElection() {
	fmt.Printf("Starting election\n")

	if !isRunningMyElection {
		isRunningMyElection = true
		if(parent >= 0){
			doSenderJob(parent, "ELECTION", myId)
		}
		electionTimer = time.NewTimer(1 * time.Second)
		go electionTimerTracker(electionTimer)
	}
}

func printFinalResults() {
	fmt.Printf("COORDINATOR ID = %d\n", coordinatorId)
	fmt.Printf("END\n")

	f, err := os.OpenFile("results/BullyImproved/" + strconv.Itoa(nPorts) + ".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	heapInfo, _ := os.Open("heap.txt")
	ints, _ := ReadInts(heapInfo)

	initConnections()
	
	for i := 0; i < nPorts; i++ {
    	if(myId == ints[4*i]){
    		parent = ints[4*i + 1]
    		child1 = ints[4*i + 2]
    		child2 = ints[4*i + 3]
    	}
    }
    fmt.Println("parent", parent)
    fmt.Println("child1", child1)
    fmt.Println("child2", child2)

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

	time.Sleep(10 * time.Second)
}
