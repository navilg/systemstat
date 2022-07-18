package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	sigar "github.com/cloudfoundry/gosigar"
	"github.com/gorilla/mux"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/moby/sys/mountinfo"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
)

func main() {
	prefix := flag.String("prefix", "", "--prefix=/systeminfo Default: /")
	flag.Parse()

	handleRequests(prefix)
}

func handleRequests(prefix *string) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/"+*prefix, systemstat).Methods("GET")
	fmt.Println("Starting server at port 30000")
	log.Fatal(http.ListenAndServe(":30000", router))
}

func systemstat(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Uptime:")
	fmt.Fprintf(w, "\n\n")
	uptime := sigar.Uptime{}
	uptime.Get()

	tw := table.NewWriter()
	tw.AppendHeader(table.Row{"Uptime"})
	tw.AppendRow(table.Row{uptime.Format()})
	tw.SetStyle(table.StyleLight)
	fmt.Fprintf(w, tw.Render())
	fmt.Fprintf(w, "\n\n")

	fmt.Fprintf(w, "CPU Utilisation (higher than real values due to loading this page):")
	fmt.Fprintf(w, "\n\n")

	cpupercentage, _ := cpu.Percent(time.Second, true)

	tw = table.NewWriter()
	tw.AppendHeader(table.Row{"Cores", "Percentage"})
	for i, cp := range cpupercentage {
		tw.AppendRow(table.Row{i + 1, fmt.Sprintf("%.2f", cp) + " %%"})
	}
	tw.SetStyle(table.StyleLight)
	fmt.Fprintf(w, tw.Render())
	fmt.Fprintf(w, "\n\n")

	fmt.Fprintf(w, "Load Average:")
	fmt.Fprintf(w, "\n\n")

	cpuinfo, _ := load.Avg()
	tw = table.NewWriter()
	tw.AppendHeader(table.Row{"1 Min", "5 Min", "15 Min"})
	tw.AppendRow(table.Row{cpuinfo.Load1, cpuinfo.Load5, cpuinfo.Load15})
	tw.SetStyle(table.StyleLight)
	fmt.Fprintf(w, tw.Render())
	fmt.Fprintf(w, "\n\n")

	info, err := mountinfo.GetMounts(mountinfo.FSTypeFilter("ext4", "fuseblk"))
	if err != nil {
		os.Exit(1)
	}

	usage := sigar.FileSystemUsage{}

	fmt.Fprintf(w, "Disk Usage:")
	fmt.Fprintf(w, "\n\n")
	tw = table.NewWriter()

	tw.AppendHeader(table.Row{"Filesystem", "Total (GB)", "Used (GB)", "Free (GB)", "Percentage Used"})
	for _, j := range info {
		usage.Get(j.Mountpoint)
		totaldisk := float64(usage.Total) / float64(1024*1024)
		useddisk := float64(usage.Used) / float64(1024*1024)
		freedisk := float64(usage.Free) / float64(1024*1024)

		useddiskpercent := (useddisk * 100) / totaldisk

		tw.AppendRow(table.Row{j.Mountpoint, fmt.Sprintf("%.2f", totaldisk), fmt.Sprintf("%.2f", useddisk), fmt.Sprintf("%.2f", freedisk), fmt.Sprintf("%.2f", useddiskpercent) + " %%"})
		// fmt.Printf("%s:Total: %.2f GB, Used: %.2f GB, Free: %.2f GB\n", j.Mountpoint, float64(disk.All)/float64(1024*1024*1024), float64(disk.Used)/float64(1024*1024*1024), float64(disk.Free)/float64(1024*1024*1024))
	}

	tw.SetStyle(table.StyleLight)
	fmt.Fprintf(w, tw.Render())
	fmt.Fprintf(w, "\n\n")

	fmt.Fprintf(w, "Memory Usage:")
	fmt.Fprintf(w, "\n\n")

	mem := sigar.Mem{}
	swap := sigar.Swap{}

	mem.Get()
	swap.Get()

	totalmem := float64(mem.Total) / float64(1024*1024*1024)
	usedmem := float64(mem.ActualUsed) / float64(1024*1024*1024)
	freemem := float64(mem.ActualFree) / float64(1024*1024*1024)
	usedmempercent := (usedmem * 100) / totalmem

	totalswap := float64(swap.Total) / float64(1024*1024*1024)
	usedswap := float64(swap.Used) / float64(1024*1024*1024)
	freeswap := float64(swap.Free) / float64(1024*1024*1024)
	usedswappercent := (usedswap * 100) / totalswap

	tw = table.NewWriter()
	tw.AppendHeader(table.Row{"", "Total (GB)", "Used (GB)", "Free (GB)", "Percentage Used"})
	tw.AppendRow(table.Row{"Memory", fmt.Sprintf("%.2f", totalmem), fmt.Sprintf("%.2f", usedmem), fmt.Sprintf("%.2f", freemem), fmt.Sprintf("%.2f", usedmempercent) + " %%"})
	tw.AppendRow(table.Row{"Swap", fmt.Sprintf("%.2f", totalswap), fmt.Sprintf("%.2f", usedswap), fmt.Sprintf("%.2f", freeswap), fmt.Sprintf("%.2f", usedswappercent)})
	tw.SetStyle(table.StyleLight)
	fmt.Fprintf(w, tw.Render())
	fmt.Fprintf(w, "\n\n")

}
