package main

import (
	"fmt"
	"sort"
	"sync"
)

type UsageContainerInfo struct {
	ContainerId string
	CpuUsage    float64 // cpu利用率
	MemUsage    int64   //已使用内存
}

type UsageContainerInfoSorter []UsageContainerInfo

func (u UsageContainerInfoSorter) Len() int {
	return len(u)
}

func (u UsageContainerInfoSorter) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (u UsageContainerInfoSorter) Less(i, j int) bool {
	if u[i].CpuUsage > u[j].CpuUsage { //cpu利用率大的排在前面
		return true
	} else if u[i].CpuUsage == u[j].CpuUsage {
		if u[i].MemUsage > u[j].MemUsage {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func NewUsageContainerOrNodeList() *UsageContainerInfoList {
	return &UsageContainerInfoList{
		UsageContainerInfos: make([]UsageContainerInfo, 2000),
		UsageInfoMap:        make(map[string]UsageContainerInfo),
	}
}

type UsageContainerInfoList struct {
	sync.Mutex
	UsageContainerInfos []UsageContainerInfo
	UsageInfoMap        map[string]UsageContainerInfo
}

var usageNodeList UsageContainerInfoList

func main22() {
	defer func() {
		fmt.Printf("invoke recover begin\n")
		p := recover()
		if p != nil {
			fmt.Printf("AcquireContainer error，%v", p)
		}
	}()
	usageNodeList = *NewUsageContainerOrNodeList()
	keys := []string{"1", "2", "3"}
	usageNodeList.UsageInfoMap["1"] = UsageContainerInfo{
		ContainerId: "1",
		CpuUsage:    0.9,  // cpu利用率
		MemUsage:    1024, //已使用内存
	}
	usageNodeList.UsageInfoMap["2"] = UsageContainerInfo{
		ContainerId: "2",
		CpuUsage:    0.8,  // cpu利用率
		MemUsage:    2048, //已使用内存
	}
	usageNodeList.UsageInfoMap["3"] = UsageContainerInfo{
		ContainerId: "3",
		CpuUsage:    0.7, // cpu利用率
		MemUsage:    512, //已使用内存
	}
	fmt.Printf("initOver=%v\n", usageNodeList.UsageInfoMap["3"])
	usageContainerInfoSorter := getUsagePair(keys)
	fmt.Printf("result=%v\n", usageContainerInfoSorter)
	for idx := range usageContainerInfoSorter {
		key := usageContainerInfoSorter[idx].ContainerId
		fmt.Printf("sotred key, key=%d==%v\n", key, usageContainerInfoSorter[idx])
	}
}

func getUsagePair(keys []string) UsageContainerInfoSorter {
	sortKeys := make([]UsageContainerInfo, len(keys))
	for idx := range keys {
		fmt.Printf("START,idx=%d\n", idx)
		containerId := keys[idx]
		usageContainerInfo := usageNodeList.UsageInfoMap[containerId]
		fmt.Printf("containerInfo=%v\n", usageNodeList.UsageInfoMap["3"])
		sortKeys[idx] = usageContainerInfo
	}
	res := UsageContainerInfoSorter(sortKeys)
	sort.Sort(res)
	return res
}
