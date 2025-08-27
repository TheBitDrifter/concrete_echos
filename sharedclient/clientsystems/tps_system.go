package clientsystems

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/hajimehoshi/ebiten/v2"
)

type PerformanceMonitorSystem struct {
	TpsThreshold       float64
	HistorySize        int
	RecoveryFrames     int
	tpsHistory         []float64
	goodFramesCount    int
	isProfiling        bool
	profileIndex       int
	currentLogFile     *os.File
	currentProfileFile *os.File
}

func NewPerformanceMonitorSystem(threshold float64, historySize int, recoveryFrames int) *PerformanceMonitorSystem {
	return &PerformanceMonitorSystem{
		TpsThreshold:   threshold,
		HistorySize:    historySize,
		RecoveryFrames: recoveryFrames,
		tpsHistory:     make([]float64, 0, historySize),
	}
}

func (s *PerformanceMonitorSystem) Run(lc coldbrew.LocalClient, scene coldbrew.Scene) error {
	// TODO: While helpful it does seem to actually cause stutters when doing File Writes which then
	// creates a loop of low tps logs being written if the cuttoff is too low — will need to rethink this!
	return nil

	currentTPS := ebiten.ActualTPS()
	s.tpsHistory = append(s.tpsHistory, currentTPS)
	if len(s.tpsHistory) > s.HistorySize {
		s.tpsHistory = s.tpsHistory[1:]
	}
	if len(s.tpsHistory) < s.HistorySize {
		return nil
	}
	var totalTPS float64
	for _, tps := range s.tpsHistory {
		totalTPS += tps
	}
	avgTPS := totalTPS / float64(len(s.tpsHistory))
	if s.isProfiling {
		if avgTPS >= s.TpsThreshold {
			s.goodFramesCount++
		} else {
			s.goodFramesCount = 0
		}
		if s.goodFramesCount >= s.RecoveryFrames {
			s.stopProfiling()
		} else {
			s.logPerformance(currentTPS)
		}
	} else {
		if avgTPS < s.TpsThreshold {
			s.startProfiling()
			s.logPerformance(currentTPS)
		}
	}
	return nil
}

func (s *PerformanceMonitorSystem) startProfiling() {
	s.isProfiling = true
	s.goodFramesCount = 0
	s.profileIndex++
	logFilename := fmt.Sprintf("performance_profile_%d.log", s.profileIndex)
	logFile, err := os.Create(logFilename)
	if err != nil {
		log.Printf("ERROR: Could not create performance log file: %v", err)
		s.isProfiling = false
		return
	}
	s.currentLogFile = logFile
	log.Printf("PERF: Low TPS detected! Started profiling to %s", logFilename)
	header := fmt.Sprintf("--- Performance Profile %d ---\nDetected at: %s\n\n", s.profileIndex, time.Now().Format(time.RFC3339))
	_, _ = s.currentLogFile.WriteString(header)
	pprofFilename := fmt.Sprintf("cpu_profile_%d.pprof", s.profileIndex)
	pprofFile, err := os.Create(pprofFilename)
	if err != nil {
		log.Printf("ERROR: Could not create cpu profile file: %v", err)
	} else {
		s.currentProfileFile = pprofFile
		pprof.StartCPUProfile(s.currentProfileFile)
	}
}

func (s *PerformanceMonitorSystem) stopProfiling() {
	if s.currentProfileFile != nil {
		pprof.StopCPUProfile()
		_ = s.currentProfileFile.Close()
		s.currentProfileFile = nil
		log.Printf("PERF: CPU profile saved.")
	}

	memProfileFilename := fmt.Sprintf("mem_profile_%d.mprof", s.profileIndex)
	memFile, err := os.Create(memProfileFilename)
	if err != nil {
		log.Printf("ERROR: could not create memory profile file: %v", err)
	} else {
		defer memFile.Close()
		if err := pprof.WriteHeapProfile(memFile); err != nil {
			log.Printf("ERROR: could not write memory profile: %v", err)
		} else {
			log.Printf("PERF: Memory profile saved to %s", memProfileFilename)
		}
	}

	if s.currentLogFile != nil {
		log.Printf("PERF: TPS recovered. Stopping profile.")
		footer := fmt.Sprintf("\n--- End of Profile ---\nRecovered at: %s\n", time.Now().Format(time.RFC3339))
		_, _ = s.currentLogFile.WriteString(footer)
		_ = s.currentLogFile.Close()
		s.currentLogFile = nil
	}
	s.isProfiling = false
}

func (s *PerformanceMonitorSystem) logPerformance(currentTPS float64) {
	if s.currentLogFile == nil {
		return
	}
	line := fmt.Sprintf("[%s] TPS: %.2f, FPS: %.2f\n", time.Now().Format("15:04:05.000"), currentTPS, ebiten.ActualFPS())
	if _, err := s.currentLogFile.WriteString(line); err != nil {
		log.Printf("ERROR: Failed to write to performance log: %v", err)
	}
}
