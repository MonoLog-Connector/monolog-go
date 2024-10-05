package client

import (
	"bufio"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/sirupsen/logrus"
)

type GinSDK struct {
	Logger     *logrus.Logger
	logFile    *os.File
	bufWriter  *bufio.Writer
	pool       sync.Pool
	shutdownCh chan struct{}
	wg         sync.WaitGroup
}

type requestInfo struct {
	DateTime      string  `json:"dateTime"`
	RequestMethod string  `json:"requestMethod"`
	RequestURL    string  `json:"requestURL"`
	Status        int     `json:"status"`
	Latency       string  `json:"latency"`
	CPUDelta      float64 `json:"cpuDelta"`
	MemoryDelta   float64 `json:"memoryDelta"`
}

func NewSDK() *GinSDK {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatal(err)
	}

	bufWriter := bufio.NewWriter(logFile)
	logger := logrus.New()
	logger.SetOutput(bufWriter)
	logger.SetFormatter(&logrus.JSONFormatter{})
	sdk := &GinSDK{
		Logger:     logger,
		logFile:    logFile,
		bufWriter:  bufWriter,
		shutdownCh: make(chan struct{}),
	}
	sdk.pool.New = func() interface{} {
		return &requestInfo{}
	}
	sdk.wg.Add(1)
	go sdk.periodicFlush()
	sdk.SetupCleanup()
	return sdk
}

func (sdk *GinSDK) periodicFlush() {
	defer sdk.wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			sdk.bufWriter.Flush()
		case <-sdk.shutdownCh:
			return
		}
	}
}

func (sdk *GinSDK) SetupCleanup() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		sdk.Logger.Info("Shutdown signal received, cleaning up...")
		close(sdk.shutdownCh)
		sdk.wg.Wait()
		sdk.bufWriter.Flush()
		sdk.logFile.Close()
		sdk.Logger.Info("Cleanup complete")
		os.Exit(0)
	}()
}

func (sdk *GinSDK) Close() {
	close(sdk.shutdownCh)
	sdk.bufWriter.Flush()
	sdk.logFile.Close()
}

func (sdk *GinSDK) GinTrackerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		info := sdk.pool.Get().(*requestInfo)
		defer sdk.pool.Put(info)
		start := time.Now()
		initialCPU, _ := cpu.Percent(0, false)
		initialMem, _ := mem.VirtualMemory()
		ctx.Next() // Process the request
		finalCPU, _ := cpu.Percent(0, false)
		finalMem, _ := mem.VirtualMemory()
		info.DateTime = start.Format(time.RFC3339)
		info.RequestMethod = ctx.Request.Method
		info.RequestURL = ctx.Request.URL.Path
		info.Status = ctx.Writer.Status()
		info.Latency = time.Since(start).String()
		info.CPUDelta = finalCPU[0] - initialCPU[0]
		info.MemoryDelta = finalMem.UsedPercent - initialMem.UsedPercent
		sdk.Logger.WithFields(logrus.Fields{
			"DateTime":      info.DateTime,
			"RequestMethod": info.RequestMethod,
			"RequestURL":    info.RequestURL,
			"Status":        info.Status,
			"Latency":       info.Latency,
			"CPU Delta":     info.CPUDelta,
			"Memory Delta":  info.MemoryDelta,
		}).Info("Request details logged")
	}
}
