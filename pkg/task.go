package pkg

import (
	"flag"
	"fmt"
	"github.com/e-XpertSolutions/f5-rest-client/f5"
	"github.com/e-XpertSolutions/f5-rest-client/f5/ltm"
	"log"
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"

var (
	WorkerNums      int
	TaskNums        int
	Host            string
	Password        string
	Username        string
	Timeout         time.Duration
	Memberip        string
	VirtualServerIP string
)

func GetRandomString(n int) string {
	bytes := []byte(letterBytes)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < n; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

type VirtualServer struct {
	Virtual_Name      string
	Vs_Destination    string
	Vs_IP_Protocol    string
	Translate_Address string
	Translate_Port    string
	Snat_Type         string
	Persistence       string
	Profiles          string
	Pool_Name         string
	Pool_Member       string
	Pool_Monitor      string
	Pool_Lbmode       string
}

func init() {
	flag.IntVar(&WorkerNums, "w", 10, "The Number of threads to start worker work")
	flag.IntVar(&TaskNums, "n", 12, "The total of task numbers")
	flag.StringVar(&Host, "a", "127.0.0.1", "the remote of host ip")
	flag.StringVar(&Username, "u", "admin", "the username of login host")
	flag.StringVar(&Password, "p", "admin", "the password of login host")
	flag.DurationVar(&Timeout, "t", 60*time.Second, "Set the timeout period for connecting to the host")
	flag.StringVar(&Memberip, "m", "", "specify the ip addess of member")
	//flag.StringVar(&virtualServerIP, "vs", "", "the specfiy  virtual server of ip addess")
}

func getMBIPAddr() string {
	rand.Seed(time.Now().UnixNano())
	ip := fmt.Sprintf("10.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}

func getVSIPAddr() string {
	rand.Seed(time.Now().UnixNano())
	ip := fmt.Sprintf("192.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}

//initiable f5 client
func NewF5Client() (*f5.Client, error) {
	hosts := fmt.Sprintf("https://" + Host)
	client, err := f5.NewBasicClient(hosts, Username, Password)
	//client, err := f5.NewBasicClient("https://192.168.5.134", "admin", "admin")
	client.DisableCertCheck()

	client.SetTimeout(Timeout * time.Second)
	if err != nil {
		fmt.Println(err)
	}
	return client, nil
}

func (vs *VirtualServer) Create(client *f5.Client) (err error) {
	var members []string
	var poolName string

	port := rand.Intn(50000)

	if Memberip == "" {
		memberIPAddr := getMBIPAddr()
		ip := fmt.Sprintf(memberIPAddr+":%d", port)
		members = append(members, ip)
	} else {
		ip := fmt.Sprintf(Memberip+":%d", port)
		members = append(members, ip)
	}

	tx, err := client.Begin()
	if err != nil {
		log.Fatal(err)
	}

	ltmclient := ltm.New(tx)

	str := GetRandomString(8)
	poolName = fmt.Sprintf("Pool_Name_%s", str)
	pool := ltm.Pool{
		Name:              poolName,
		Monitor:           vs.Pool_Monitor,
		LoadBalancingMode: vs.Pool_Lbmode,
		Members:           members,
	}
	//create pool
	if err := ltmclient.Pool().Create(pool); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("pool name %s create success.\n", poolName)

	vsName := fmt.Sprintf("Virtual_Name_%s", str)

	VirtualServerIP = getVSIPAddr()
	vsIP := fmt.Sprintf(VirtualServerIP+":%d", port)

	vss := ltm.VirtualServer{
		Name:                     vsName,
		Destination:              vsIP,
		IPProtocol:               vs.Vs_IP_Protocol,
		TranslateAddress:         vs.Translate_Address,
		TranslatePort:            vs.Translate_Port,
		Profiles:                 []string{vs.Profiles},
		Persistences:             []ltm.Persistence{{Name: vs.Persistence}},
		Pool:                     poolName,
		SourceAddressTranslation: ltm.SourceAddressTranslation{Type: vs.Snat_Type},
	}
	// create virtual server
	if err = ltmclient.Virtual().Create(vss); err != nil {
		log.Fatal(err)
	}
	if err = tx.Commit(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("virtualserver name %s create success.\n", vsName)
	return nil
}

//Job interface
type Job interface {
	Create(client *f5.Client) (err error)
}

// The job processor, where the real business logic is handled, is responsible for receiving and processing tasks,
// and it needs to tell the scheduler if it is ready to receive more tasks.
type worker struct {
	// Multiple workers share a worker queue WorkerQueue, which is used to register their own work channel (chan Job) to WorkerQueue
	// when the worker is idle, so that the worker can receive task requests
	workerQueue chan chan Job
	// jobQueue is an unbuffered job channel that receives Job
	jobQueue chan Job
	quit     chan bool
}

// Initialize a worker thread
func NewWorkers(workPool chan chan Job) *worker {
	return &worker{
		workerQueue: workPool,
		jobQueue:    make(chan Job),
		quit:        make(chan bool),
	}
}

// Define a start method for the thread to indicate that it is listening for a task to begin processing
func (w *worker) Start(client *f5.Client, ch chan struct{}) {
	go func() {
		for {
			w.workerQueue <- w.jobQueue //Register the worker channel to the thread pool
			select {
			case task := <-w.jobQueue:
				if err := task.Create(client); err != nil {
					log.Fatalf("create configure failed :%s", err)
				}
				ch <- struct{}{} // The goroutine ends, then send to signals
			case <-w.quit:
				return
			}
		}
	}()
}

// The thread stops working
func (w *worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

//The task distributor can distribute the tasks in the task queue to the threads in the thread pool one by one for processing
type Dispatcher struct {
	WorkerQueue chan chan Job
	MaxNum      int
	JobQueue    chan Job
}

// Instantiate a task dispatcher
func NewDispatcher(maxWorkerNum int) *Dispatcher {
	return &Dispatcher{
		WorkerQueue: make(chan chan Job, maxWorkerNum),
		MaxNum:      maxWorkerNum,
		JobQueue:    make(chan Job),
	}
}

// Assign Tasks
func (d *Dispatcher) Dispatch() {
	for {
		select {
		// Remove a task from the task queue
		case jobObj := <-d.JobQueue:
			go func(job Job) {
				// Take a thread out of the thread pool
				workChan := <-d.WorkerQueue
				// A task from the task queue is processed by the worker thread
				workChan <- job
			}(jobObj)
		}

	}
}

// Start Task allocator starts running and distributing tasks
func (d *Dispatcher) Run(client *f5.Client, ch chan struct{}) {
	//Creating a new worker Thread
	for i := 0; i < d.MaxNum; i++ {
		workerObj := NewWorkers(d.WorkerQueue)
		//Start the thread
		workerObj.Start(client, ch)
	}
	// Distribute Tasks
	go d.Dispatch()
}
