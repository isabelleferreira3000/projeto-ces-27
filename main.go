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

var nPorts int
var isCandidate int

var SendersConn []*net.UDPConn
var ReceiversConn *net.UDPConn

var ch = make(chan int)

type ClockStruct struct {
	Id int
	Clocks []int
}
var logicalClock ClockStruct

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

	var otherLogicalClock ClockStruct
	err = json.Unmarshal(buf[:n], &otherLogicalClock)
	CheckError(err)

	fmt.Println("Received", otherLogicalClock)
	myId := logicalClock.Id
	myClocks := logicalClock.Clocks
	otherProcessClocks := otherLogicalClock.Clocks

	// updating clocks
	logicalClock.Clocks[myId-1]++
	for i := 0; i < nPorts; i++ {
		logicalClock.Clocks[i] = max(otherProcessClocks[i], myClocks[i])
	}

	fmt.Println("logicalClock atualizado:", logicalClock)
}

func doSenderJob(otherProcessID int) {
	otherProcess := otherProcessID - 1

	jsonRequest, err := json.Marshal(logicalClock)
	CheckError(err)

	_, err = SendersConn[otherProcess].Write(jsonRequest)
	CheckError(err)

	time.Sleep(time.Second * 1)
}

func initConnections() {
	//nPorts = len(os.Args) - 2

	// getting my Id
	myId, err := strconv.Atoi(os.Args[1])
	CheckError(err)

	// getting my port
	myPort = os.Args[myId + 1]

	// creating logicalClock
	var clocks []int
	for i := 0; i < nPorts; i++ {
		clocks = append(clocks, 0)
	}
	logicalClock = ClockStruct{
		myId,
		clocks,
	}

	// Server
	ServerAddr, err := net.ResolveUDPAddr("udp", myPort)
	CheckError(err)
	aux, err := net.ListenUDP("udp", ServerAddr)
	ReceiversConn = aux
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
	fmt.Printf("nPorts: %d\nisCandidate: %d\n", nPorts, isCandidate)

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
				// updating my clock
				myId := logicalClock.Id
				logicalClock.Clocks[myId-1]++
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
