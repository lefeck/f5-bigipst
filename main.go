package main

import (
	"f5-bigipst/pkg"
	"flag"
)

func main() {
	flag.Parse()
	ch := make(chan struct{})
	client, _ := pkg.NewF5Client()
	// Initializes the task dispatcher
	dispatcher := pkg.NewDispatcher(pkg.WorkerNums)
	dispatcher.Run(client, ch)
	vs := &pkg.VirtualServer{
		Vs_IP_Protocol:    "tcp",
		Translate_Address: "enabled",
		Translate_Port:    "enabled",
		Snat_Type:         "automap",
		Persistence:       "source_addr",
		Pool_Lbmode:       "round-robin",
		Pool_Monitor:      "gateway_icmp",
		Profiles:          "tcp-mobile-optimized",
	}

	go func() {
		for i := 1; i <= pkg.TaskNums; i++ {
			dispatcher.JobQueue <- vs
		}
	}()

	for range ch {
		// Each time data is received from the ch, it indicates the end of an active concurrent process
		pkg.TaskNums--
		// Close the pipeline when all active concurrent processes are finished
		if pkg.TaskNums == 0 {
			close(ch)
		}
	}
}
