package util

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/beego/beego/v2/adapter/toolbox"
	"github.com/beego/beego/v2/core/logs"
	"github.com/bwmarrin/snowflake"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"strconv"
	"strings"
	"sync"
	"time"
)

var node *snowflake.Node
var fileTypeMap sync.Map

func init() {
	node, _ = snowflake.NewNode(1)
	initFileType()
	//initTaskShowSysLoad()
}

func GeneKey() string {
	id := node.Generate()
	logs.Trace("GeneKey() ", id.String())
	return id.String()
}
func GeneId() string {
	id := node.Generate()
	logs.Trace("GeneId()", id.String())
	return id.String()
}

func StrToIntArray(str string) ([]int, error) {
	res := make([]int, 0)
	stringSlice := strings.Split(str, ",")
	for _, v := range stringSlice {
		buf, err := strconv.Atoi(v)
		if err == nil && buf >= 1 {
			res = append(res, buf)
		}
	}
	return res, nil
}

func GetMemInfo() mem.VirtualMemoryStat {
	m, _ := mem.VirtualMemory()
	return *m
}
func GetCpuIdle() float64 {
	c, _ := cpu.Percent(time.Second, false)
	return c[0]
}
func ShowMemoryInfo() {
	m := GetMemInfo()
	sub := [...]string{" ", "▍", "▍", "▌", "▌", "▋", "▋", "▊", "▊", "▉"}
	width := 25
	Lcount := int(m.UsedPercent / 100.0 * float64(width))
	mid := sub[int(m.UsedPercent)%10/(100/width)+1]
	Rcount := width - Lcount - 1
	s := fmt.Sprintf("[%s%s%s]", strings.Repeat("█", Lcount), mid, strings.Repeat(" ", Rcount))
	fmt.Printf("[INFO]total memory used [%5s / %5s]MB %s %.1f%% \n", strconv.FormatUint(m.Used>>20, 10), strconv.FormatUint(m.Total>>20, 10), s, m.UsedPercent)
}
func ShowCpuInfo() {
	c := GetCpuIdle()
	sub := [...]string{" ", "▍", "▍", "▌", "▌", "▋", "▋", "▊", "▊", "▉"}
	width := 25
	Lcount := int(c / 100.0 * float64(width))
	mid := sub[int(c)%10/(100/width)]
	Rcount := width - Lcount - 1
	s := fmt.Sprintf("[%s%s%s]", strings.Repeat("█", Lcount), mid, strings.Repeat(" ", Rcount))
	fmt.Printf("[INFO]total cpu used    [%5s / %5s]%%  %s %.1f%% \n", strconv.FormatUint(uint64(uint(c)>>20), 10), strconv.FormatUint(uint64(c)>>20, 10), s, c)
}
func initTaskShowSysLoad() {
	tk := toolbox.NewTask("util.SysLoad", "1 * * * * *", func() error { ShowCpuInfo(); ShowMemoryInfo(); return nil })
	toolbox.AddTask("util.SysLoad", tk)
}

func initFileType() {
	fileTypeMap.Store("ffd8ffe000104a464946", "jpg") //JPEG (jpg)
	fileTypeMap.Store("89504e470d0a1a0a0000", "png") //PNG (png)
	fileTypeMap.Store("47494638396126026f01", "gif") //GIF (gif)
	fileTypeMap.Store("49492a00227105008037", "tif") //TIFF (tif)
	fileTypeMap.Store("424d228c010000000000", "bmp") //16色位图(bmp)
	fileTypeMap.Store("424d8240090000000000", "bmp") //24位位图(bmp)
	fileTypeMap.Store("424d8e1b030000000000", "bmp") //256色位图(bmp)
}

// 获取前面结果字节的二进制
func bytesToHexString(src []byte) string {
	res := bytes.Buffer{}
	if src == nil || len(src) <= 0 {
		return ""
	}
	temp := make([]byte, 0)
	for _, v := range src {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10))
		}
		res.WriteString(hv)
	}
	return res.String()
}

// GetFileType 判断文件类型
// fSrc: 文件字节流（就用前面几个字节）
func GetFileType(fSrc []byte) string {
	var fileType string
	fileCode := bytesToHexString(fSrc)

	fileTypeMap.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		if strings.HasPrefix(fileCode, strings.ToLower(k)) ||
			strings.HasPrefix(k, strings.ToLower(fileCode)) {
			fileType = v
			return false
		}
		return true
	})
	return fileType
}
