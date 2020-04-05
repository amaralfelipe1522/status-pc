package stats

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"

	linuxproc "github.com/c9s/goprocinfo/linux"
)

//Stats armazenará as informações retornadas de todos os métodos presente em main.go
type Stats struct {
	MemAlloc  uint64     `json:"memoryalloc"`
	MemSys    uint64     `json:"memorysys"`
	OSCPU     [8]uint64  `json:"oscpu"`
	UserCPU   [8]uint64  `json:"usercpu"`
	DiskName  [12]string `json:"diskname"`
	DiskMajor [12]int    `json:"diskmajor"`
}

//PrintMemory bla bla bla
func (st *Stats) PrintMemory() {
	mem := new(runtime.MemStats)
	runtime.ReadMemStats(mem)
	//fmt.Printf("Memória alocada no momento é %v bytes (%vmb).\n", mem.Alloc, bToMb(mem.Alloc))
	st.MemAlloc = mem.Alloc
	//fmt.Printf("Memória utilizada no Sistema Operacional no momento é %v bytes (%vmb).\n", mem.Sys, bToMb(mem.Sys))
	st.MemSys = mem.Sys
}

//PrintCPU bla bla
func (st *Stats) PrintCPU() {
	status, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Fatal("Falha na leitura")
	}

	for i, s := range status.CPUStats {
		//fmt.Printf("O uso de CPU atual no Sistema Operacional é %v bytes (%vmb).\n", s.System, bToMb(s.System))
		st.OSCPU[i] = s.System
		//fmt.Printf("O uso de CPU atual do usuário é %v bytes (%vmb).\n", s.User, bToMb(s.User))
		st.UserCPU[i] = s.User
	}
}

//PrintDisc bla bla
func (st *Stats) PrintDisc() {
	status, err := linuxproc.ReadDiskStats("/proc/diskstats")
	if err != nil {
		log.Fatal("Falha na leitura")
	}

	for i, s := range status {
		if i <= 11 {
			//fmt.Printf("Major %v.\n", s.Major)
			st.DiskMajor[i] = s.Major
			// fmt.Printf("Minor %v.\n", s.Minor)
			//fmt.Printf("Name %v.\n", s.Name)
			st.DiskName[i] = s.Name
			// fmt.Printf("ReadIOs %v.\n", s.ReadIOs)
		}
	}
	/*
		status, err := linuxproc.ReadDisk("/proc/diskstats")
		if err != nil {
			log.Fatal("Falha na leitura")
		}
		fmt.Println("Iniciando leitura do disco...")
		fmt.Printf("Uso de disco atual é %v.\n", status.Used)
		fmt.Printf("Qtd de disco livre é %v.\n", status.Free)
		fmt.Printf("Total de disco é %v.\n", status.All)
		fmt.Println("Leitura de disco finalizada.")
	}*/
}

//BToMb bla bla
func BToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// PrintStatus bla bla
func PrintStatus() string {
	stats := Stats{}
	stats.PrintMemory()
	stats.PrintCPU()
	stats.PrintDisc()

	tojson, err := json.Marshal(stats)
	if err != nil {
		log.Fatal(err)
	}
	return string(tojson)
}

//PrintStatusHTTP bla bla
func PrintStatusHTTP(w http.ResponseWriter, r *http.Request) {
	stats := Stats{}
	stats.PrintMemory()
	stats.PrintCPU()
	stats.PrintDisc()

	tojson, err := json.Marshal(stats)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(tojson))
}

//Print bla blaa
func Print() {
	http.HandleFunc("/", PrintStatusHTTP)
}
