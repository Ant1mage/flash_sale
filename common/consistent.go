package common

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type units []uint32

// 返回长度
func (x units) Len() int {
	return len(x)
}

// 对比大小
func (x units) Less(i, j int) bool {
	return x[i] < x[j]
}

// 交换
func (x units) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

var errEmpty = errors.New("Hash 环没有数据")

type Consistent struct {
	// hash环, key为哈希值, value 存放节点信息
	circle       map[uint32]string
	sortedHashes units
	// 虚拟节点个数
	VirtualNode int
	// 读写锁
	sync.RWMutex
}

func NewConsistent() *Consistent {
	return &Consistent{
		circle:       make(map[uint32]string),
		sortedHashes: nil,
		VirtualNode:  20,
		RWMutex:      sync.RWMutex{},
	}
}

// 自动生成 key 值
func (c *Consistent) generateKey(element string, index int) string {
	return element + strconv.Itoa(index)
}

// 获取 hash 位置
func (c *Consistent) hashkey(key string) uint32 {
	if len(key) < 64 {
		var srcatch [64]byte
		copy(srcatch[:], key)
		// IEEE 多项式返回数据的 CRC-32 校验和
		return crc32.ChecksumIEEE(srcatch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

func (c *Consistent) updateSortedHashes() {
	hashes := c.sortedHashes[:0]
	// 判断切片容量
	if cap(c.sortedHashes)/(c.VirtualNode*4) > len(c.circle) {
		hashes = nil
	}
	for k := range c.circle {
		hashes = append(hashes, k)
	}
	// 对所有节点排序, 方便二分查找
	sort.Sort(hashes)
	// 重新赋值
	c.sortedHashes = hashes
}

func (c *Consistent) Add(element string) {
	c.Lock()
	defer c.Unlock()
	c.add(element)
}

// 添加节点
func (c *Consistent) add(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		c.circle[c.hashkey(c.generateKey(element, i))] = element
	}
	// 更新排序
	c.updateSortedHashes()
}

func (c *Consistent) Remove(element string) {
	c.Lock()
	defer c.Unlock()
	c.remove(element)
}

// 删除节点
func (c *Consistent) remove(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		delete(c.circle, c.hashkey(c.generateKey(element, i)))
	}
	c.updateSortedHashes()
}

// 顺时针查找最近的节点
func (c *Consistent) search(key uint32) int {
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}
	// 二分查找搜索指定节点
	i := sort.Search(len(c.sortedHashes), f)
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return i
}

func (c *Consistent) Get(name string) (string, error) {
	c.RLock()
	defer c.Unlock()
	if len(c.circle) == 0 {
		return "", errEmpty
	}
	// 计算哈希值
	key := c.hashkey(name)
	i := c.search(key)
	return c.circle[c.sortedHashes[i]], nil
}
