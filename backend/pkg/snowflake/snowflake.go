package snowflake

import (
	"sync"

	sf "github.com/bwmarrin/snowflake"
)

var (
	node *sf.Node
	once sync.Once
)

// Init 初始化 Snowflake 节点
// nodeID: 节点ID，范围 0-1023
func Init(nodeID int64) error {
	var err error
	once.Do(func() {
		node, err = sf.NewNode(nodeID)
	})
	return err
}

// GenerateID 生成一个 Snowflake ID
func GenerateID() string {
	if node == nil {
		// 如果未初始化，默认使用节点 1
		_ = Init(1)
	}
	return node.Generate().String()
}
