package sd

import (
	"fmt"
	h "hr-server/handler"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

// @Summary Shows OK as the ping-pong result
// @Description Shows OK as the ping-pong result
// @Tags sd
// @Accept  json
// @Produce  json
// @Success 200 {string} plain "OK"
// @Router /v1/sd/health [get]
func HealthCheck(c *gin.Context) {
	_, disk := readDisk()
	_, ram := readRAM()
	_, cpu := readCPU()

	h.SendResponse(c, nil, CreateResponse{Disk: disk, CPU: cpu, RAM: ram})
}

func readDisk() (status int, health Health) {
	u, _ := disk.Usage("/")

	health.UsedMB = int(u.Used) / MB
	health.UsedGB = int(u.Used) / GB
	health.TotalMB = int(u.Total) / MB
	health.TotalGB = int(u.Total) / GB
	health.UsedPercent = int(u.UsedPercent)

	status = http.StatusOK
	//text := "OK"

	if health.UsedPercent >= 95 {
		status = http.StatusOK
		//text = "CRITICAL"
	} else if health.UsedPercent >= 90 {
		status = http.StatusTooManyRequests
		//text = "WARNING"
	}

	//message := fmt.Sprintf("%s - Free space: %dMB (%dGB) / %dMB (%dGB) | Used: %d%%", text, usedMB, usedGB, totalMB, totalGB, usedPercent)
	return status, health
}

func readRAM() (status int, health Health) {
	u, _ := mem.VirtualMemory()

	health.UsedMB = int(u.Used) / MB
	health.UsedGB = int(u.Used) / GB
	health.TotalMB = int(u.Total) / MB
	health.TotalGB = int(u.Total) / GB
	health.UsedPercent = int(u.UsedPercent)

	status = http.StatusOK

	if health.UsedPercent >= 95 {
		status = http.StatusInternalServerError
	} else if health.UsedPercent >= 90 {
		status = http.StatusTooManyRequests
	}

	//message := fmt.Sprintf("%s - Free space: %dMB (%dGB) / %dMB (%dGB) | Used: %d%%", text, usedMB, usedGB, totalMB, totalGB, usedPercent)
	return status, health
}

func readCPU() (status int, health Health) {
	a, _ := load.Avg()
	health.Load1 = a.Load1
	health.Load5 = a.Load5
	health.Load15 = a.Load15

	return status, health
}

// @Summary Checks the disk usage
// @Description Checks the disk usage
// @Tags sd
// @Accept  json
// @Produce  json
// @Success 200 {string} plain "OK - Free space: 17233MB (16GB) / 51200MB (50GB) | Used: 33%"
// @Router /v1/sd/disk [get]
func DiskCheck(c *gin.Context) {
	status, data := readDisk()
	h.SendResponse(c, nil, CreateResponse{Status: status, Disk: data})
}

func getCPUSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}

// @Summary Checks the cpu usage
// @Description Checks the cpu usage
// @Tags sd
// @Accept  json
// @Produce  json
// @Success 200 {string} plain "CRITICAL - Load average: 1.78, 1.99, 2.02 | Cores: 2"
// @Router /v1/sd/cpu [get]
func CPUCheck(c *gin.Context) {
	idle0, total0 := getCPUSample()
	time.Sleep(3 * time.Second)
	idle1, total1 := getCPUSample()

	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)
	cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks

	health := Health{}
	health.CpuUsage = cpuUsage
	health.CpuTotal = totalTicks
	h.SendResponse(c, nil, CreateResponse{CPU: health})
}

// @Summary Checks the ram usage
// @Description Checks the ram usage
// @Tags sd
// @Accept  json
// @Produce  json
// @Success 200 {string} plain "OK - Free space: 402MB (0GB) / 8192MB (8GB) | Used: 4%"
// @Router /v1/sd/ram [get]
func RAMCheck(c *gin.Context) {
	status, data := readRAM()
	h.SendResponse(c, nil, CreateResponse{Status: status, RAM: data})
}
