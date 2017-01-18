package utils

import "fmt"

type HostResponse struct {
	ArpCache       []*ArpCache       `json:"arp_cache,omitempty"`
	CurrentRam     *CurrentRam       `json:"current_ram,omitempty"`
	DiskPartitions []*DiskPartitions `json:"disk_partitions,omitempty"`
	Host           *ConfigHost       `json:"host"`
}

func (r *HostResponse) Print() {
	if r == nil {
		return
	}

	if len(r.ArpCache) > 0 {
		printArpCache(r.ArpCache)
		fmt.Println("")
	}
	if r.CurrentRam != nil {
		printCurrentRam(r.CurrentRam)
		fmt.Println("")
	}
	if len(r.DiskPartitions) > 0 {
		printDiskPartitions(r.DiskPartitions)
		fmt.Println("")
	}
}

type ArpCache struct {
	Addr   string `json:"addr"`
	HwType string `json:"hw_type"`
	HwAddr string `json:"hw_addr"`
	Mask   string `json:"mask"`
}

func printArpCache(a []*ArpCache) {
	fmt.Println("Arp Cache:")
	for _, arp := range a {
		fmt.Printf("Address: %s\n", arp.Addr)
		fmt.Printf("HwType: %s\n", arp.HwType)
		fmt.Printf("HwAddr: %s\n", arp.HwAddr)
		fmt.Printf("Mask: %s\n\n", arp.Mask)
	}
}

type CurrentRam struct {
	Total     float64 `json:"total"`
	Used      float64 `json:"used"`
	Available float64 `json:"available"`
}

func printCurrentRam(c *CurrentRam) {
	fmt.Println("Current Ram:")
	fmt.Printf("Total: %.2fMB\n", c.Total)
	fmt.Printf("Used: %.2fMB\n", c.Used)
	fmt.Printf("Available: %.2fMB\n", c.Available)
}

type DiskPartitions struct {
	FileSystem  string `json:"file_system"`
	Size        string `json:"size"`
	Used        string `json:"used"`
	Avail       string `json:"avail"`
	UsedPercent string `json:"used%"`
	Mounted     string `json:"mounted"`
}

func printDiskPartitions(d []*DiskPartitions) {
	fmt.Println("Disk Partitions:")
	for _, partition := range d {
		fmt.Printf("Filesystem: %s\n", partition.FileSystem)
		fmt.Printf("Size: %s\n", partition.Size)
		fmt.Printf("Used: %s\n", partition.Used)
		fmt.Printf("Available: %s\n", partition.Avail)
		fmt.Printf("Used %%: %s\n", partition.UsedPercent)
		fmt.Printf("Mount Point: %s\n\n", partition.Mounted)
	}
}
