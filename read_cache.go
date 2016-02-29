package transit_go

import "github.com/nedap/transit-go/constants"

type ReadCache interface {
	CacheRead(str string, asMapKey bool, parser Parser) interface{}
	Init()
}

type readCache struct {
	cache []string
	index int
}

func NewReadCache() ReadCache {
	cache := &readCache{}
	cache.Init()
	return cache
}

func (c *readCache) Init() {
	c.cache = make([]string, constants.MaxCacheEntries, constants.MaxCacheEntries)
	c.index = 0
}

func (c *readCache) CacheRead(str string, asMapKey bool, parser Parser) interface{} {
	if len(str) != 0 {
		if cacheCode(str) {
			realVal := c.cache[codeToIndex(str)]
			return realVal
		} else if isCacheable(str, asMapKey) {
			if c.index == constants.MaxCacheEntries {
				c.Init()
			}
			value := str
			if parser != nil {
				parseResult, _ := parser.parseString(str)
				tag, isTag := parseResult.(Tag)
				if isTag {
					return tag
				}
			}
			c.cache[c.index] = value
			c.index++
		}
	}
	if parser != nil {
		parseResult, _ := parser.parseString(str)
		tag, isTag := parseResult.(Tag)
		if isTag {
			return tag
		} else {
			return parseResult
		}
	} else {
		return str
	}
}

func cacheCode(str string) bool {
	return str[0] == constants.SUB && str != constants.MAP_AS_ARRAY
}

func codeToIndex(code string) int {
	length := len(code)
	if length == 2 {
		val := int(code[1] - constants.BaseCharIndex)
		return val
	} else {
		val := (int(code[1]-constants.BaseCharIndex) * constants.CacheCodeDigits) +
			(int(code[2] - constants.BaseCharIndex))
		return val
	}
}
