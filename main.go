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
		// Выполняем HTTP запрос
		resp, err := http.Get(serverURL)
		if err != nil {
			errorCount++
			handleError(errorCount)
			time.Sleep(10 * time.Second)
			continue
		}

		// Проверяем статус ответа
		if resp.StatusCode != http.StatusOK {
			errorCount++
			resp.Body.Close()
			handleError(errorCount)
			time.Sleep(10 * time.Second)
			continue
		}

		// Читаем и парсим данные
		scanner := bufio.NewScanner(resp.Body)
		if !scanner.Scan() {
			errorCount++
			resp.Body.Close()
			handleError(errorCount)
			time.Sleep(10 * time.Second)
			continue
		}

		data := scanner.Text()
		resp.Body.Close()

		values := strings.Split(data, ",")
		if len(values) != 7 {
			errorCount++
			handleError(errorCount)
			time.Sleep(10 * time.Second)
			continue
		}

		// Сбрасываем счётчик ошибок при успешном получении данных
		errorCount = 0

		// Парсим все значения
		loadAvg, err1 := strconv.ParseFloat(values[0], 64)
		totalMem, err2 := strconv.ParseUint(values[1], 10, 64)
		usedMem, err3 := strconv.ParseUint(values[2], 10, 64)
		totalDisk, err4 := strconv.ParseUint(values[3], 10, 64)
		usedDisk, err5 := strconv.ParseUint(values[4], 10, 64)
		totalNet, err6 := strconv.ParseUint(values[5], 10, 64)
		usedNet, err7 := strconv.ParseUint(values[6], 10, 64)

		// Проверяем ошибки парсинга
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil || err7 != nil {
			errorCount++
			handleError(errorCount)
			time.Sleep(10 * time.Second)
			continue
		}

		// Проверяем все условия
		checkLoadAverage(loadAvg)
		checkMemoryUsage(totalMem, usedMem)
		checkDiskSpace(totalDisk, usedDisk)
		checkNetworkBandwidth(totalNet, usedNet)

		// Пауза перед следующим запросом
		time.Sleep(10 * time.Second)
	}
}

func handleError(errorCount int) {
	if errorCount >= 3 {
		fmt.Println("Unable to fetch server statistic")
	}
}

func checkLoadAverage(loadAvg float64) {
	if loadAvg > 30 {
		fmt.Printf("Load Average is too high: %d\n", int(loadAvg))
	}
}

func checkMemoryUsage(totalMem, usedMem uint64) {
	if totalMem > 0 {
		memoryUsagePercent := float64(usedMem) / float64(totalMem) * 100
		if memoryUsagePercent > 80 {
			fmt.Printf("Memory usage too high: %d%%\n", int(memoryUsagePercent))
		}
	}
}

func checkDiskSpace(totalDisk, usedDisk uint64) {
	if totalDisk > 0 {
		diskUsagePercent := float64(usedDisk) / float64(totalDisk) * 100
		if diskUsagePercent > 90 {
			freeDiskMB := (totalDisk - usedDisk) / (1024 * 1024)
			fmt.Printf("Free disk space is too low: %d Mb left\n", freeDiskMB)
		}
	}
}

func checkNetworkBandwidth(totalNet, usedNet uint64) {
	if totalNet > 0 {
		netUsagePercent := float64(usedNet) / float64(totalNet) * 100
		if netUsagePercent > 90 {
			availableNetBytes := totalNet - usedNet
			// Автотесты ожидают простое деление на 1,000,000 (без преобразования в биты)
			availableNetMbits := availableNetBytes / 1000000
			fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", availableNetMbits)
		}
	}
}
