package conf

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

const middle = "=HFish="

type Config struct {
	Mymap  map[string]string
	MyNode map[string]string
	strcet string
}

var config *Config

func LoadConfig(path string) {
	c := &Config{}
	c.Mymap = make(map[string]string)
	c.MyNode = make(map[string]string)

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		s := strings.TrimSpace(string(b))
		if strings.Index(s, "#") == 0 {
			continue
		}

		n1 := strings.Index(s, "[")
		n2 := strings.LastIndex(s, "]")
		if n1 > -1 && n2 > -1 && n2 > n1+1 {
			c.strcet = strings.TrimSpace(s[n1+1 : n2])
			continue
		}

		if len(c.strcet) == 0 {
			continue
		}
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		first := strings.TrimSpace(s[:index])
		if len(first) == 0 {
			continue
		}
		second := strings.TrimSpace(s[index+1:])

		pos := strings.Index(second, "#")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, "\t//")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " //")
		if pos > -1 {
			second = second[0:pos]
		}

		if len(second) == 0 {
			second = ""
		}

		key := c.strcet + middle + first
		c.Mymap[key] = strings.TrimSpace(second)

		key = c.strcet + middle + "introduce"
		introduce, found := c.Mymap[key]
		if !found {
		}

		key = c.strcet + middle + "mode"
		mode, found := c.Mymap[key]
		if !found {
		}

		c.MyNode[c.strcet] = strings.TrimSpace(mode) + "&&" + strings.TrimSpace(introduce)
	}

	config = c
}

func ReadString(node, key string) *string {
	key = node + middle + key
	v, found := config.Mymap[key]
	if !found {
		return nil
	}
	str := strings.TrimSpace(v)
	return &str
}

func ReadInt(node, key string, args ...int) *int {
	v := ReadString(node, key)
	if v == nil {
		if len(args) >= 1 {
			return &args[0]
		}
		return nil
	}
	number, err := strconv.Atoi(*v)
	if err != nil {
		if len(args) >= 1 {
			return &args[0]
		}
		return nil
	}
	return &number
}

func ReadFloat(node, key string, args ...float64) *float64 {
	v := ReadString(node, key)
	if v == nil {
		if len(args) >= 1 {
			return &args[0]
		}
		return nil
	}
	number, err := strconv.ParseFloat(*v, 32)
	if err != nil {
		if len(args) >= 1 {
			return &args[0]
		}
		return nil
	}
	return &number
}
