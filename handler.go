package stats_api

import (
	"encoding/json"
	"io"
	"net/http"
	"runtime"
	"strconv"
)

var newLineTerm bool = false

type resultRuntime struct {
	GoVersion    string `json:"go_version"`
	GoOs         string `json:"go_os"`
	GoArch       string `json:"go_arch"`
	CpuNum       int    `json:"cpu_num"`
	GoroutineNum int    `json:"goroutine_num"`
	Gomaxprocs   int    `json:"gomaxprocs"`
	CgoCallNum   int64  `json:"cgo_call_num"`
}

type resultMemory struct {
	MemoryAlloc      uint64 `json:"alloc"`
	MemoryTotalAlloc uint64 `json:"total_alloc"`
	MemorySys        uint64 `json:"sys"`
	MemoryLookups    uint64 `json:"lookups"`
	MemoryMallocs    uint64 `json:"mallocs"`
	MemoryFrees      uint64 `json:"frees"`
}

type resultHeap struct {
	HeapAlloc    uint64 `json:"alloc"`
	HeapSys      uint64 `json:"sys"`
	HeapIdle     uint64 `json:"idle"`
	HeapInuse    uint64 `json:"inuse"`
	HeapReleased uint64 `json:"released"`
	HeapObjects  uint64 `json:"objects"`
}

type resultGc struct {
	GcNext uint64 `json:"next"`
	GcLast uint64 `json:"last"`
	GcNum  uint32 `json:"num"`
}

type result struct {
	Rr resultRuntime `json:"runtime"`
	Rm resultMemory  `json:"memory"`
	Rh resultHeap    `json:"heap"`
	Rg resultGc      `json:"gc"`
}

func NewLineTermEnabled() {
	newLineTerm = true
}

func NewLineTermDisabled() {
	newLineTerm = false
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	result := &result{
		Rr: resultRuntime{
			GoVersion:    runtime.Version(),
			GoOs:         runtime.GOOS,
			GoArch:       runtime.GOARCH,
			CpuNum:       runtime.NumCPU(),
			GoroutineNum: runtime.NumGoroutine(),
			Gomaxprocs:   runtime.GOMAXPROCS(0),
			CgoCallNum:   runtime.NumCgoCall(),
		},
		Rm: resultMemory{
			MemoryAlloc:      mem.Alloc,
			MemoryTotalAlloc: mem.TotalAlloc,
			MemorySys:        mem.Sys,
			MemoryLookups:    mem.Lookups,
			MemoryMallocs:    mem.Mallocs,
			MemoryFrees:      mem.Frees,
		},
		Rh: resultHeap{
			HeapAlloc:    mem.HeapAlloc,
			HeapSys:      mem.HeapSys,
			HeapIdle:     mem.HeapIdle,
			HeapInuse:    mem.HeapInuse,
			HeapReleased: mem.HeapReleased,
			HeapObjects:  mem.HeapObjects,
		},
		Rg: resultGc{
			GcNext: mem.NextGC,
			GcLast: mem.LastGC,
			GcNum:  mem.NumGC,
		},
	}

	jsonBytes, jsonErr := json.Marshal(result)
	var body string
	if jsonErr != nil {
		body = jsonErr.Error()
	} else {
		body = string(jsonBytes)
	}

	if newLineTerm {
		body += "\n"
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Content-Length"] = strconv.Itoa(len(body))
	for name, value := range headers {
		w.Header().Set(name, value)
	}

	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	io.WriteString(w, body)
}
