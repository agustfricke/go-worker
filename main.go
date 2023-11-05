package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/net"
)

type ServerStats struct {
	MemoryUsage   float64 `json:"memory_usage"`
	CPUUsage      float64 `json:"cpu_usage"`
	NetworkStats  []NetworkInterfaceStats `json:"network_stats"`
}

type NetworkInterfaceStats struct {
	BytesSent  uint64 `json:"bytes_sent"`
	BytesRecv  uint64 `json:"bytes_recv"`
	PacketsSent uint64 `json:"packets_sent"`
	PacketsRecv uint64 `json:"packets_recv"`
}

func getServerStats() ServerStats {
	memInfo, _ := mem.VirtualMemory()
	cpuInfo, _ := cpu.Percent(0, false)
	netInfo, _ := net.IOCounters(false)

	stats := ServerStats{
		MemoryUsage: memInfo.UsedPercent,
		CPUUsage:    cpuInfo[0],
		NetworkStats: make([]NetworkInterfaceStats, len(netInfo)),
	}

	for i, netStat := range netInfo {
		stats.NetworkStats[i] = NetworkInterfaceStats{
			BytesSent:  netStat.BytesSent,
			BytesRecv:  netStat.BytesRecv,
			PacketsSent: netStat.PacketsSent,
			PacketsRecv: netStat.PacketsRecv,
		}
	}

	return stats
}

func serverStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats := getServerStats()
	jsonData, err := json.Marshal(stats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func main() {
	http.HandleFunc("/server-stats", serverStatsHandler)

	go func() {
		for {
			time.Sleep(5 * time.Second)
			fmt.Println("Server is running...")
		}
	}()

	http.ListenAndServe(":8080", nil)
}
