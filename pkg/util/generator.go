package util

import (
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
)

// ReContextOrder 重新进行排除
func ReContextOrder(c *generator.Context, canonicalOrderName string) {
	for name, systemNamer := range c.Namers {
		if name == canonicalOrderName {
			orderer := namer.Orderer{Namer: systemNamer}
			c.Order = orderer.OrderUniverse(c.Universe)
		}
	}
}
