//go:build darwin

package gpu

import "testing"

func TestParseDarwinPerformanceStatistics(t *testing.T) {
	stats := parseDarwinPerformanceStatistics(`"In use system memory (driver)"=0,"Alloc system memory"=4598005760,"Tiler Utilization %"=52,"Renderer Utilization %"=48,"Device Utilization %"=52,"In use system memory"=1085030400`)
	if stats["Device Utilization %"] != 52 {
		t.Fatalf("expected device util 52, got %v", stats["Device Utilization %"])
	}
	if stats["Renderer Utilization %"] != 48 {
		t.Fatalf("expected renderer util 48, got %v", stats["Renderer Utilization %"])
	}
	if stats["In use system memory"] != 1085030400 {
		t.Fatalf("expected in-use memory, got %v", stats["In use system memory"])
	}
}

func TestParseDarwinAccelerators(t *testing.T) {
	sample := `
+-o AGXAcceleratorG16G  <class AGXAcceleratorG16G, id 0x1000005d3, registered, matched, active, busy 0 (1632 ms), retain 59>
    {
      "PerformanceStatistics" = {"In use system memory (driver)"=0,"Alloc system memory"=4598005760,"Tiler Utilization %"=52,"Renderer Utilization %"=48,"Device Utilization %"=52,"In use system memory"=1085030400}
      "IOSourceVersion" = "353.14"
      "MetalPluginName" = "AGXMetalG16G_B0"
      "model" = "Apple M4"
      "gpu-core-count" = 8
      "IONameMatched" = "gpu,t8132"
    }
`
	items := parseDarwinAccelerators(sample)
	if len(items) != 1 {
		t.Fatalf("expected 1 accelerator, got %d", len(items))
	}
	device := darwinDeviceFromAccelerator(items[0], 0)
	if device.GPUUtil != "52 %" {
		t.Fatalf("unexpected gpu util: %q", device.GPUUtil)
	}
	if device.ProductName != "Apple M4" {
		t.Fatalf("unexpected product name: %q", device.ProductName)
	}
	if device.MemUsed == "" || device.MemTotal == "" {
		t.Fatalf("expected memory metrics, got used=%q total=%q", device.MemUsed, device.MemTotal)
	}
}

func TestDarwinGPULoadInfoIntegration(t *testing.T) {
	ok, client := New()
	if !ok {
		t.Skip("darwin GPU provider unavailable")
	}
	info, err := client.LoadInfoContext(t.Context())
	if err != nil {
		t.Fatalf("load darwin gpu info failed: %v", err)
	}
	if len(info.Devices) == 0 {
		t.Fatalf("expected at least one GPU device")
	}
	if info.Type != "apple" {
		t.Fatalf("expected apple type, got %q", info.Type)
	}
}

func TestParseDarwinDisplayInfo(t *testing.T) {
	payload := `{
  "SPDisplaysDataType": [
    {
      "_name": "Apple M4",
      "spdisplays_mtlgpufamilysupport": "spdisplays_metal4",
      "sppci_cores": "8",
      "sppci_model": "Apple M4",
      "spdisplays_vendor": "sppci_vendor_Apple"
    }
  ],
  "SPSoftwareDataType": [
    {
      "os_version": "macOS 26.6.2 (25G83)"
    }
  ]
}`
	displays, err := parseDarwinDisplayInfo(payload)
	if err != nil {
		t.Fatalf("parse display info failed: %v", err)
	}
	if len(displays) != 1 {
		t.Fatalf("expected 1 display entry, got %d", len(displays))
	}
	if displays[0].Model != "Apple M4" || displays[0].MetalFamily != "spdisplays_metal4" {
		t.Fatalf("unexpected display info: %+v", displays[0])
	}
}
