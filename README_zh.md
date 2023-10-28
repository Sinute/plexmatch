[English](README.md) 

# 介绍

plexmatch 提供了结构化解析和生成 [.plexmatch](https://support.plex.tv/articles/plexmatch/) 文件的功能

# 如何使用

## 从 .plexmatch 文件解析内容

```go
package main

import (
	"fmt"

	"github.com/sinute/plexmatch"
)

func main() {
	data := `title: 【推しの子】
		tmdbid: 203737
		season: 1
		
		ep: SP01: Oshi no Ko - sp.mkv
		ep: 01: Oshi no Ko - 01.mkv
		ep: 02: Oshi no Ko - 02.mkv`
	p, err := plexmatch.Unmarshal([]byte(data))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", p)
	// {
	// 	Title:【推しの子】
	// 	Show:
	// 	TMDBID:203737
	// 	TVDBID:0
	// 	IMDBID:0
	// 	Year:0
	// 	Season:1
	// 	Episodes:[
	// 		{
	// 			ID:1
	// 			Name:Oshi no Ko - sp.mkv
	// 			Special:true
	// 		}
	// 		{
	// 			ID:1
	// 			Name:Oshi no Ko - 01.mkv
	// 			Special:false
	// 		}
	// 		{
	// 			ID:2
	// 			Name:Oshi no Ko - 02.mkv
	// 			Special:false
	// 		}
	// 	]
	// }
}

```

## 生成 .plexmatch 文件

```go
package main

import (
	"fmt"

	"github.com/sinute/plexmatch"
)

func main() {
	p := plexmatch.PlexMatch{
		Title:  "【推しの子】",
		TMDBID: 203737,
		Season: 1,
		Episodes: []plexmatch.Episode{
			{ID: 1, Name: "Oshi no Ko - sp.mkv", Special: true},
			{ID: 1, Name: "Oshi no Ko - 01.mkv", Special: false},
			{ID: 2, Name: "Oshi no Ko - 02.mkv", Special: false},
		},
	}
	fmt.Println(string(plexmatch.Marshal(p)))
	// title: 【推しの子】
	//
	// tmdbid: 203737
	//
	// season: 1
	//
	// ep: SP01: Oshi no Ko - sp.mkv
	// ep: 01: Oshi no Ko - 01.mkv
	// ep: 02: Oshi no Ko - 02.mkv
}

```
