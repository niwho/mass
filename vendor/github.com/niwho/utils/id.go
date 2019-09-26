package utils

import (
	"github.com/bwmarrin/snowflake"
)

var IDFactory *snowflake.Node

func InitIDFactory(node int64) (err error) {
	IDFactory, err = snowflake.NewNode(node)
	return
}

func GetNewID() int64 {
	return int64(IDFactory.Generate())
}
