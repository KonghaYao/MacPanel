//go:build darwin

package service

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/1Panel-dev/1Panel/agent/app/dto/response"
	"github.com/1Panel-dev/1Panel/agent/utils/cmd"
	gdisk "github.com/shirou/gopsutil/v4/disk"
)

type diskutilListJSON struct {
	AllDisksAndPartitions []diskutilDiskJSON `json:"AllDisksAndPartitions"`
}

type diskutilDiskJSON struct {
	DeviceIdentifier   string                  `json:"DeviceIdentifier"`
	Size               int64                   `json:"Size"`
	Content            string                  `json:"Content"`
	OSInternal         bool                    `json:"OSInternal"`
	Partitions         []diskutilPartitionJSON `json:"Partitions"`
	APFSVolumes        []diskutilAPFSVolumeJSON `json:"APFSVolumes"`
	APFSPhysicalStores []struct {
		DeviceIdentifier string `json:"DeviceIdentifier"`
	} `json:"APFSPhysicalStores"`
}

type diskutilPartitionJSON struct {
	DeviceIdentifier string `json:"DeviceIdentifier"`
	Size             int64  `json:"Size"`
	Content          string `json:"Content"`
	MountPoint       string `json:"MountPoint"`
	VolumeName       string `json:"VolumeName"`
}

type diskutilAPFSVolumeJSON struct {
	DeviceIdentifier string `json:"DeviceIdentifier"`
	Size             int64  `json:"Size"`
	MountPoint       string `json:"MountPoint"`
	VolumeName       string `json:"VolumeName"`
	CapacityInUse    int64  `json:"CapacityInUse"`
	OSInternal       bool   `json:"OSInternal"`
	MountedSnapshots []struct {
		SnapshotMountPoint string `json:"SnapshotMountPoint"`
		SnapshotBSD        string `json:"SnapshotBSD"`
	} `json:"MountedSnapshots"`
}

type diskutilInfoJSON struct {
	DeviceIdentifier string `json:"DeviceIdentifier"`
	MediaName        string `json:"MediaName"`
	SolidState       bool   `json:"SolidState"`
	Removable        bool   `json:"Removable"`
	Internal         bool   `json:"Internal"`
	VolumeName       string `json:"VolumeName"`
}

func (s *DiskService) GetCompleteDiskInfo() (*response.CompleteDiskInfo, error) {
	disks, err := loadDiskutilList()
	if err != nil {
		return nil, err
	}

	diskByID := make(map[string]diskutilDiskJSON, len(disks))
	partitionToPhysical := make(map[string]string)
	for _, disk := range disks {
		diskByID[disk.DeviceIdentifier] = disk
		if disk.Content != "GUID_partition_scheme" {
			continue
		}
		for _, part := range disk.Partitions {
			partitionToPhysical[part.DeviceIdentifier] = disk.DeviceIdentifier
		}
	}

	var result response.CompleteDiskInfo
	seenVolumes := make(map[string]struct{})

	for _, disk := range disks {
		if disk.Content != "GUID_partition_scheme" {
			continue
		}

		info, err := loadDiskutilInfo(disk.DeviceIdentifier)
		if err != nil {
			info = diskutilInfoJSON{DeviceIdentifier: disk.DeviceIdentifier}
		}

		diskInfo := response.DiskInfo{
			DiskBasicInfo: response.DiskBasicInfo{
				Device:      "/dev/" + disk.DeviceIdentifier,
				Size:        formatDarwinBytes(disk.Size),
				Model:       info.MediaName,
				DiskType:    darwinDiskType(info),
				IsRemovable: info.Removable,
				IsSystem:    info.Internal && !info.Removable,
				Filesystem:  disk.Content,
			},
			Partitions: collectDarwinVolumes(disk, diskByID, partitionToPhysical, seenVolumes),
		}

		if diskInfo.IsSystem {
			result.SystemDisks = append(result.SystemDisks, diskInfo)
		} else if len(diskInfo.Partitions) == 0 {
			result.UnpartitionedDisks = append(result.UnpartitionedDisks, diskInfo.DiskBasicInfo)
		} else {
			result.Disks = append(result.Disks, diskInfo)
		}
		result.TotalDisks++
		result.TotalCapacity += disk.Size
	}

	return &result, nil
}

func loadDiskutilList() ([]diskutilDiskJSON, error) {
	cmdMgr := cmd.NewCommandMgr(cmd.WithTimeout(20 * time.Second))
	plistOutput, err := cmdMgr.RunWithStdout("diskutil", "list", "-plist")
	if err != nil {
		return nil, fmt.Errorf("failed to run diskutil list: %w", err)
	}

	plutilCmd := exec.Command("plutil", "-convert", "json", "-o", "-", "-")
	plutilCmd.Stdin = strings.NewReader(plistOutput)
	jsonOutput, err := plutilCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to convert diskutil plist: %w", err)
	}

	var payload diskutilListJSON
	if err := json.Unmarshal(jsonOutput, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse diskutil list: %w", err)
	}
	return payload.AllDisksAndPartitions, nil
}

func loadDiskutilInfo(deviceID string) (diskutilInfoJSON, error) {
	cmdMgr := cmd.NewCommandMgr(cmd.WithTimeout(10 * time.Second))
	plistOutput, err := cmdMgr.RunWithStdout("diskutil", "info", "-plist", deviceID)
	if err != nil {
		return diskutilInfoJSON{}, err
	}

	plutilCmd := exec.Command("plutil", "-convert", "json", "-o", "-", "-")
	plutilCmd.Stdin = strings.NewReader(plistOutput)
	jsonOutput, err := plutilCmd.Output()
	if err != nil {
		return diskutilInfoJSON{}, err
	}

	var info diskutilInfoJSON
	if err := json.Unmarshal(jsonOutput, &info); err != nil {
		return diskutilInfoJSON{}, err
	}
	return info, nil
}

func collectDarwinVolumes(
	physical diskutilDiskJSON,
	diskByID map[string]diskutilDiskJSON,
	partitionToPhysical map[string]string,
	seenVolumes map[string]struct{},
) []response.DiskBasicInfo {
	var volumes []response.DiskBasicInfo

	for _, part := range physical.Partitions {
		if part.Content == "Apple_APFS" || part.Content == "Apple_APFS_ISC" || part.Content == "Apple_APFS_Recovery" {
			for _, disk := range diskByID {
				for _, store := range disk.APFSPhysicalStores {
					if store.DeviceIdentifier != part.DeviceIdentifier {
						continue
					}
					volumes = append(volumes, collectAPFSVolumes(disk, seenVolumes)...)
				}
			}
			continue
		}

		if part.MountPoint == "" && part.VolumeName == "" {
			continue
		}
		volumes = append(volumes, buildDarwinVolume(
			part.DeviceIdentifier,
			part.VolumeName,
			part.Content,
			part.MountPoint,
			part.Size,
			0,
			false,
		))
	}

	return volumes
}

func collectAPFSVolumes(container diskutilDiskJSON, seenVolumes map[string]struct{}) []response.DiskBasicInfo {
	var volumes []response.DiskBasicInfo
	for _, volume := range container.APFSVolumes {
		deviceID := volume.DeviceIdentifier
		if _, exists := seenVolumes[deviceID]; exists {
			continue
		}

		mountPoint := volume.MountPoint
		if mountPoint == "" {
			for _, snapshot := range volume.MountedSnapshots {
				if snapshot.SnapshotMountPoint != "" {
					mountPoint = snapshot.SnapshotMountPoint
					deviceID = snapshot.SnapshotBSD
					break
				}
			}
		}
		if !shouldShowDarwinAPFSVolume(volume, mountPoint) {
			continue
		}

		seenVolumes[deviceID] = struct{}{}
		volumes = append(volumes, buildDarwinVolume(
			deviceID,
			volume.VolumeName,
			"apfs",
			mountPoint,
			volume.Size,
			volume.CapacityInUse,
			volume.OSInternal || isDarwinSystemMount(mountPoint),
		))
	}
	return volumes
}

func buildDarwinVolume(deviceID, volumeName, filesystem, mountPoint string, size, capacityInUse int64, isSystem bool) response.DiskBasicInfo {
	info := response.DiskBasicInfo{
		Device:     "/dev/" + deviceID,
		Size:       formatDarwinBytes(size),
		Model:      volumeName,
		DiskType:   "APFS",
		Filesystem: filesystem,
		MountPoint: mountPoint,
		IsMounted:  mountPoint != "",
		IsSystem:   isSystem,
	}

	if mountPoint != "" {
		totalSize, used, avail, usePercent := darwinUsage(mountPoint)
		if totalSize != "" {
			info.Size = totalSize
		}
		if used != "" {
			info.Used = used
		}
		if avail != "" {
			info.Avail = avail
		}
		info.UsePercent = usePercent
	} else if capacityInUse > 0 {
		info.Used = formatDarwinBytes(capacityInUse)
	}

	return info
}

func darwinUsage(mountPoint string) (size, used, avail string, usePercent int) {
	usage, err := gdisk.Usage(mountPoint)
	if err != nil || usage == nil {
		return "", "", "", 0
	}
	return formatDarwinBytes(int64(usage.Total)),
		formatDarwinBytes(int64(usage.Used)),
		formatDarwinBytes(int64(usage.Free)),
		int(usage.UsedPercent)
}

func shouldShowDarwinAPFSVolume(volume diskutilAPFSVolumeJSON, mountPoint string) bool {
	if mountPoint == "" {
		return false
	}
	if mountPoint == "/" || mountPoint == "/System/Volumes/Data" {
		return true
	}
	return strings.HasPrefix(mountPoint, "/Volumes/")
}

func isDarwinSystemMount(mountPoint string) bool {
	if mountPoint == "" {
		return false
	}
	switch mountPoint {
	case "/", "/System/Volumes/Data":
		return true
	default:
		return strings.HasPrefix(mountPoint, "/System/Volumes/")
	}
}

func darwinDiskType(info diskutilInfoJSON) string {
	if info.SolidState {
		return "SSD"
	}
	return "HDD"
}

func formatDarwinBytes(bytes int64) string {
	if bytes <= 0 {
		return "0B"
	}
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	value := float64(bytes) / float64(div)
	if value >= 100 {
		return fmt.Sprintf("%.0f%ci", value, "KMGTPE"[exp])
	}
	return fmt.Sprintf("%.1f%ci", value, "KMGTPE"[exp])
}
