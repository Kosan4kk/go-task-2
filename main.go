package main

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	serverURL := "http://srv.msk01.gigacorp.local/_stats"
	errorCount := 0

	for {
		resp, err := http.Get(serverURL)
		if err != nil {
			errorCount++
			fmt.Printf("Error making request: %v\n", err)
		} else {
			if resp.StatusCode != http.StatusOK {
				errorCount++
				fmt.Printf("Server returned non-200 status: %d\n", resp.StatusCode)
			} else {
				scanner := bufio.NewScanner(resp.Body)
				if scanner.Scan() {
					data := scanner.Text()
					values := strings.Split(data, ",")
					if len(values) != 7 {
						errorCount++
						fmt.Printf("Invalid data format: expected 7 values, got %d\n", len(values))
					} else {
						errorCount = 0
						
						loadAvg, _ := strconv.ParseFloat(values[0], 64)
						totalMem, _ := strconv.ParseUint(values[1], 10, 64)
						usedMem, _ := strconv.ParseUint(values[2], 10, 64)
						totalDisk, _ := strconv.ParseUint(values[3], 10, 64)
						usedDisk, _ := strconv.ParseUint(values[4], 10, 64)
						totalNet, _ := strconv.ParseUint(values[5], 10, 64)
						usedNet, _ := strconv.ParseUint(values[6], 10, 64)
						
						// 1. Load Average
						if loadAvg > 30 {
							fmt.Printf("Load Average is too high: %d\n", int(loadAvg))
						}
						
						// 2. Memory usage
						if totalMem > 0 {
							memoryUsagePercent := float64(usedMem) / float64(totalMem) * 100
							if memoryUsagePercent > 80 {
								fmt.Printf("Memory usage too high: %d%%\n", int(memoryUsagePercent))
							}
						}
						
						// 3. Disk space
						if totalDisk > 0 {
							diskUsagePercent := float64(usedDisk) / float64(totalDisk) * 100
							if diskUsagePercent > 90 {
								freeDiskMB := int((totalDisk - usedDisk) / (1024 * 1024))
								fmt.Printf("Free disk space is too low: %d Mb left\n", freeDiskMB)
							}
						}
						
						// 4. Network bandwidth - КРИТИЧЕСКОЕ ИСПРАВЛЕНИЕ
						if totalNet > 0 {
							netUsagePercent := float64(usedNet) / float64(totalNet) * 100
							if netUsagePercent > 90 {
								availableNetBytes := totalNet - usedNet
								// Правильный расчёт: байты/сек → мегабиты/сек
								availableNetMbits := int(float64(availableNetBytes) / 131072) // 131072 = (1024 * 1024) / 8
								fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", availableNetMbits)
							}
						}
					}
				} else {
					errorCount++
					fmt.Printf("Failed to read response body: %v\n", scanner.Err())
				}
			}
			resp.Body.Close()
		}
		
		if errorCount >= 3 {
			fmt.Println("Unable to fetch server statistic")
			errorCount = 0
		}
		
		time.Sleep(10 * time.Second)
	}
}
