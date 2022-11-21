package pkg

import (
	"flag"
	"fmt"
	"github.com/e-XpertSolutions/f5-rest-client/f5"
	"github.com/e-XpertSolutions/f5-rest-client/f5/ltm"
	"github.com/rs/xid"
	"log"
	"math/rand"
	"time"
)

var (
	WorkerNums      int
	TaskNums        int
	Host            string
	Password        string
	Username        string
	Timeout         time.Duration
	MemberIP        string
	VirtualServerIP string
	File            string
	Partition       string
)

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
	flag.IntVar(&TaskNums, "n", 10, "The total of task numbers")
	flag.StringVar(&Host, "a", "127.0.0.1", "the remote of host ip")
	flag.StringVar(&Username, "u", "admin", "the username of login host")
	flag.StringVar(&Password, "p", "admin", "the password of login host")
	flag.DurationVar(&Timeout, "t", 60*time.Second, "Set the timeout period for connecting to the host")
	flag.StringVar(&MemberIP, "m", "", "Specify the ip addess of member, If you don't specify an IP address, an IP address of 10.0.0.0/8 will be generated randomly.")
	flag.StringVar(&File, "f", "", "Specify the file location of the output results")
	flag.StringVar(&VirtualServerIP, "vs", "", "the specfiy  virtual server of ip addess, If you don't specify an IP address, an IP address of 192.0.0.0/8 will be generated randomly.")
	//unknown finished
	flag.StringVar(&Partition, "P", "Common", "the specfiy  the location of partition")
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
	host := fmt.Sprintf("https://" + Host)
	client, err := f5.NewBasicClient(host, Username, Password)
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

	rand.Seed(time.Now().UnixNano())
	port := rand.Intn(50000)

	if MemberIP == "" {
		// the ip address will be generated randomly.
		memberIPAddr := getMBIPAddr()
		ip := fmt.Sprintf(memberIPAddr+":%d", port)
		members = append(members, ip)
	} else {
		//Add IP manually
		ips := ParseIP(MemberIP)
		ip := fmt.Sprintf(ips+":%d", port)
		members = append(members, ip)
	}

	tx, err := client.Begin()
	if err != nil {
		log.Fatalf("client open transaction: %s", err)
	}

	ltmclient := ltm.New(tx)

	uid := xid.New()
	poolName = fmt.Sprintf("Pool_Name_%s", uid)
	pool := ltm.Pool{
		Name:              poolName,
		Monitor:           vs.Pool_Monitor,
		LoadBalancingMode: vs.Pool_Lbmode,
		Members:           members,
	}
	//create pool
	if err := ltmclient.Pool().Create(pool); err != nil {
		log.Fatalf("create pool failed: %s", err)
	}

	if File == "" {
		fmt.Printf("pool name %s create success.\n", poolName)
	} else {
		poolResult := fmt.Sprintf("pool name %s create success.\n", poolName)
		WriteFile(poolResult, File)
	}

	vsName := fmt.Sprintf("Virtual_Name_%s", uid)

	var vsIP string
	if VirtualServerIP == "" {
		// the ip address will be generated randomly.
		VirtualServerIP = getVSIPAddr()
		vsIP = fmt.Sprintf(VirtualServerIP+":%d", port)
	} else {
		//Add IP manually
		ips := ParseIP(VirtualServerIP)
		vsIP = fmt.Sprintf(ips+":%d", port)
	}

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
		log.Fatalf("create virtualserver failed: %s", err)
	}
	if err = tx.Commit(); err != nil {
		log.Fatalf("client open transaction: %s", err)
	}

	if File == "" {
		fmt.Printf("virtualserver name %s create success.\n", vsName)
	} else {
		vsResult := fmt.Sprintf("virtualserver name %s create success.\n", vsName)
		WriteFile(vsResult, File)
	}

	return nil
}

func CreatePartition(client *f5.Client) (err error) {
	tx, err := client.Begin()
	if err != nil {
		log.Fatalf("client open transaction: %s", err)
	}
	cmd := fmt.Sprintf("tmsh create auth partition " + Partition)
	tx.Exec(cmd)
	if err = tx.Commit(); err != nil {
		log.Fatalf("client commits transaction: %s", err)
	}
	return nil
}
