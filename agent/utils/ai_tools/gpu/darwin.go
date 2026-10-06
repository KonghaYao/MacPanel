//go:build darwin

package gpu

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/1Panel-dev/1Panel/agent/utils/cmd"
)

type darwinGPU struct{}

func platformProviders() []provider {
	return []provider{darwinGPU{}}
}

func (d darwinGPU) LoadInfo(ctx context.Context) (*Info, error) {
	cmdMgr := cmd.NewCommandMgr(cmd.WithContext(ctx), cmd.WithTimeout(8*time.Second))
	ioregData, err := cmdMgr.RunWithStdout("ioreg", "-r", "-d", "1", "-w", "0", "-c", "IOAccelerator")
	if err != nil {
		return nil, fmt.Errorf("calling ioreg failed: %w", err)
	}

	accelerators := parseDarwinAccelerators(ioregData)
	if len(accelerators) == 0 {
		return nil, fmt.Errorf("no macOS GPU accelerators found")
	}

	displays, displayErr := loadDarwinDisplayInfo(ctx)
	if displayErr != nil {
		accelerators[0].warnings = append(accelerators[0].warnings, displayErr.Error())
	}

	info := &Info{Type: "apple"}
	for index, item := range accelerators {
		device := darwinDeviceFromAccelerator(item, index)
		if index < len(displays) {
			applyDarwinDisplayInfo(&device, displays[index])
		}
		info.Devices = append(info.Devices, device)
		if info.DriverVersion == "" && item.driverVersion != "" {
			info.DriverVersion = item.driverVersion
		}
	}
	if len(accelerators) > 0 {
		info.Warnings = append(info.Warnings, accelerators[0].warnings...)
	}
	if info.DriverVersion == "" && len(displays) > 0 && displays[0].OSVersion != "" {
		info.DriverVersion = displays[0].OSVersion
	}
	return info, nil
}

type darwinAccelerator struct {
	model           string
	coreCount       string
	metalPlugin     string
	driverVersion   string
	busID           string
	vramTotalMB     int64
	stats           map[string]float64
	warnings        []string
}

type darwinDisplayInfo struct {
	Model        string
	CoreCount    string
	MetalFamily  string
	Vendor       string
	OSVersion    string
}

func darwinDeviceFromAccelerator(item darwinAccelerator, index int) Device {
	device := Device{
		Type:          "apple",
		Index:         uint(index),
		ProductName:   firstNonEmpty(item.model, "Apple GPU"),
		Architecture:  formatMetalFamily(item.metalPlugin),
		DriverVersion: item.driverVersion,
		BusID:         item.busID,
		ProcessStatus: "unavailable",
	}

	if util, ok := item.stats["Device Utilization %"]; ok {
		device.GPUUtil = formatPercent(util)
	}
	if util, ok := item.stats["Renderer Utilization %"]; ok {
		device.EncoderUtil = formatPercent(util)
	}
	if util, ok := item.stats["Tiler Utilization %"]; ok {
		device.DecoderUtil = formatPercent(util)
	}

	if inUse, ok := item.stats["In use system memory"]; ok && inUse > 0 {
		device.MemUsed = formatBytesMiB(inUse)
	}
	if alloc, ok := item.stats["Alloc system memory"]; ok && alloc > 0 {
		device.MemTotal = formatBytesMiB(alloc)
		device.MemoryReserved = formatBytesMiB(alloc)
	}
	if item.vramTotalMB > 0 {
		device.MemTotal = fmt.Sprintf("%d MiB", item.vramTotalMB)
		if free, ok := item.stats["vramFreeBytes"]; ok && free >= 0 {
			totalBytes := item.vramTotalMB * 1024 * 1024
			usedBytes := float64(totalBytes) - free
			if usedBytes < 0 {
				usedBytes = 0
			}
			device.MemUsed = formatBytesMiB(usedBytes)
			device.MemoryFree = formatBytesMiB(free)
		}
	}
	if device.ProductName != "" {
		device.UUID = fmt.Sprintf("apple:%s:%d", strings.ReplaceAll(device.ProductName, " ", "-"), index)
	}
	if item.coreCount != "" {
		if device.Architecture != "" {
			device.Architecture = fmt.Sprintf("%s · %s cores", device.Architecture, item.coreCount)
		} else {
			device.Architecture = item.coreCount + " cores"
		}
	}
	return device
}

func applyDarwinDisplayInfo(device *Device, display darwinDisplayInfo) {
	if device.ProductName == "" || device.ProductName == "Apple GPU" {
		device.ProductName = firstNonEmpty(display.Model, device.ProductName)
	}
	if display.MetalFamily != "" {
		device.Architecture = formatMetalFamily(display.MetalFamily)
		if display.CoreCount != "" {
			device.Architecture = fmt.Sprintf("%s · %s cores", device.Architecture, display.CoreCount)
		}
	} else if display.CoreCount != "" && !strings.Contains(device.Architecture, "cores") {
		if device.Architecture != "" {
			device.Architecture = fmt.Sprintf("%s · %s cores", device.Architecture, display.CoreCount)
		} else {
			device.Architecture = display.CoreCount + " cores"
		}
	}
	if device.BusID == "" && display.Vendor != "" {
		device.BusID = display.Vendor
	}
}

func loadDarwinDisplayInfo(ctx context.Context) ([]darwinDisplayInfo, error) {
	cmdMgr := cmd.NewCommandMgr(cmd.WithContext(ctx), cmd.WithTimeout(8*time.Second))
	data, err := cmdMgr.RunWithStdout("system_profiler", "SPDisplaysDataType", "SPSoftwareDataType", "-json")
	if err != nil {
		return nil, fmt.Errorf("calling system_profiler failed: %w", err)
	}
	return parseDarwinDisplayInfo(data)
}

func parseDarwinDisplayInfo(data string) ([]darwinDisplayInfo, error) {
	var payload map[string]json.RawMessage
	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		return nil, err
	}

	osVersion := ""
	if raw, ok := payload["SPSoftwareDataType"]; ok {
		var software []map[string]any
		if err := json.Unmarshal(raw, &software); err == nil && len(software) > 0 {
			if value, ok := software[0]["os_version"].(string); ok {
				osVersion = value
			}
		}
	}

	var displays []darwinDisplayInfo
	raw, ok := payload["SPDisplaysDataType"]
	if !ok {
		return displays, nil
	}
	var entries []map[string]any
	if err := json.Unmarshal(raw, &entries); err != nil {
		return nil, err
	}
	for _, entry := range entries {
		info := darwinDisplayInfo{OSVersion: osVersion}
		if value, ok := entry["sppci_model"].(string); ok {
			info.Model = value
		}
		if value, ok := entry["_name"].(string); ok && info.Model == "" {
			info.Model = value
		}
		if value, ok := entry["sppci_cores"].(string); ok {
			info.CoreCount = value
		}
		if value, ok := entry["spdisplays_mtlgpufamilysupport"].(string); ok {
			info.MetalFamily = value
		}
		if value, ok := entry["spdisplays_vendor"].(string); ok {
			info.Vendor = strings.TrimPrefix(value, "sppci_vendor_")
		}
		displays = append(displays, info)
	}
	return displays, nil
}

func parseDarwinAccelerators(data string) []darwinAccelerator {
	blocks := splitDarwinIORegBlocks(data)
	items := make([]darwinAccelerator, 0, len(blocks))
	for _, block := range blocks {
		item := darwinAccelerator{stats: map[string]float64{}}
		if value := extractDarwinQuotedValue(block, "model"); value != "" {
			item.model = value
		}
		if value := extractDarwinQuotedValue(block, "gpu-core-count"); value != "" {
			item.coreCount = value
		}
		if value := extractDarwinQuotedValue(block, "MetalPluginName"); value != "" {
			item.metalPlugin = value
		}
		if value := extractDarwinQuotedValue(block, "IOSourceVersion"); value != "" {
			item.driverVersion = value
		}
		if value := extractDarwinQuotedValue(block, "IONameMatched"); value != "" {
			item.busID = strings.Trim(value, "\"")
		}
		if value := extractDarwinNumericValue(block, "VRAM,totalMB"); value > 0 {
			item.vramTotalMB = int64(value)
		}
		if statsLine := extractDarwinPerformanceStatistics(block); statsLine != "" {
			item.stats = parseDarwinPerformanceStatistics(statsLine)
		}
		if len(item.stats) > 0 || item.model != "" || item.vramTotalMB > 0 {
			items = append(items, item)
		}
	}
	return items
}

var darwinIORegBlockPattern = regexp.MustCompile(`(?m)^\+-o `)

func splitDarwinIORegBlocks(data string) []string {
	indices := darwinIORegBlockPattern.FindAllStringIndex(data, -1)
	if len(indices) == 0 {
		return nil
	}
	blocks := make([]string, 0, len(indices))
	for i, index := range indices {
		start := index[0]
		end := len(data)
		if i+1 < len(indices) {
			end = indices[i+1][0]
		}
		blocks = append(blocks, data[start:end])
	}
	return blocks
}

func extractDarwinQuotedValue(block, key string) string {
	pattern := regexp.MustCompile(`"` + regexp.QuoteMeta(key) + `" = "([^"]*)"`)
	match := pattern.FindStringSubmatch(block)
	if len(match) == 2 {
		return match[1]
	}
	return ""
}

func extractDarwinNumericValue(block, key string) float64 {
	pattern := regexp.MustCompile(`"` + regexp.QuoteMeta(key) + `" = (\d+)`)
	match := pattern.FindStringSubmatch(block)
	if len(match) != 2 {
		return 0
	}
	value, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return 0
	}
	return value
}

func extractDarwinPerformanceStatistics(block string) string {
	const marker = `"PerformanceStatistics" = {`
	start := strings.Index(block, marker)
	if start < 0 {
		return ""
	}
	start += len(marker)
 depth := 1
	for i := start; i < len(block); i++ {
		switch block[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return block[start:i]
			}
		}
	}
	return ""
}

func parseDarwinPerformanceStatistics(data string) map[string]float64 {
	stats := map[string]float64{}
	for _, part := range splitDarwinDictEntries(data) {
		eq := strings.Index(part, "=")
		if eq <= 0 {
			continue
		}
		key := strings.TrimSpace(part[:eq])
		key = strings.Trim(key, `"`)
		valueText := strings.TrimSpace(part[eq+1:])
		value, err := strconv.ParseFloat(valueText, 64)
		if err != nil {
			continue
		}
		stats[key] = value
	}
	return stats
}

func splitDarwinDictEntries(data string) []string {
	var (
		parts   []string
		current strings.Builder
		inQuote bool
	)
	for i := 0; i < len(data); i++ {
		ch := data[i]
		switch ch {
		case '"':
			inQuote = !inQuote
			current.WriteByte(ch)
		case ',':
			if inQuote {
				current.WriteByte(ch)
				continue
			}
			part := strings.TrimSpace(current.String())
			if part != "" {
				parts = append(parts, part)
			}
			current.Reset()
		default:
			current.WriteByte(ch)
		}
	}
	if tail := strings.TrimSpace(current.String()); tail != "" {
		parts = append(parts, tail)
	}
	return parts
}

func formatPercent(value float64) string {
	return fmt.Sprintf("%.0f %%", value)
}

func formatBytesMiB(value float64) string {
	return fmt.Sprintf("%.0f MiB", value/1024/1024)
}

func formatMetalFamily(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if strings.HasPrefix(value, "spdisplays_metal") {
		return "Metal " + strings.TrimPrefix(value, "spdisplays_metal")
	}
	if strings.HasPrefix(value, "AGXMetal") {
		return value
	}
	return value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
