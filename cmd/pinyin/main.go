package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/gmd20/golang-pinyin/pinyin"
)

func main() {
	var spaceDelimiter bool
	var keepNonChinese bool
	var style int
	flag.BoolVar(&spaceDelimiter, "d", true, "space delimiter")
	flag.BoolVar(&keepNonChinese, "k", true, "keep non-Chinese characters")
	flag.IntVar(&style, "s", 0, "output style")
	flag.Parse()

	d := pinyin.NewDict(spaceDelimiter, keepNonChinese, style)
	reader := bufio.NewReader(os.Stdin)
	for {
		s, _ := reader.ReadString('\n')
		if len(s) == 0 {
			break
		}
		s = d.RenMing(s)
		fmt.Println(s)
	}
}
