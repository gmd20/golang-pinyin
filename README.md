# golang-pinyin
把汉字转换为拼音golang库


# dict.go 拼音数据来源
主要是来自https://github.com/Lofanmi/pinyin-golang/blob/master/pinyin/dict.go，
但感觉他那个项目实现效率太低了，所以改写一下。
他的数据使用https://github.com/overtrue/pinyin/blob/master/data/ 这个php项目里面来的吧，
这种词组的方式感觉比https://github.com/mozillazg/pinyin-data单字的要好一些，
但pinyin-data会有多音字的数据，目前暂时用不上吧。



# 类似的项目
https://github.com/mozillazg/go-pinyin   
https://github.com/Lofanmi/pinyin-golang   


# 用法和例子

```go
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
```

