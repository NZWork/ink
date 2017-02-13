![ink](https://pbs.twimg.com/profile_images/1216524203/Ink_twitter.jpg)

ink 是 tiki的核心底层服务，负责用户 markdown 文件转换（ markdown => HTML）

# 大体构思
## markdown源文件结构
单层文件结构，通过 Hash 进行索引查找，每个仓库均为一个独立的 wiki 项目，源文件存储在仓库的 `markdown` 目录下

## 存储
转换后存储在仓库内的 `parsed` 目录下，并且会热缓存在 redis 中以 `<repo>_<hash>` 的形式存储

## 完整流程
读取文件 => 转换 => 写缓存 => 写文件


## 当前性能
测试流程：读取=>转换=>存储

测试markdown源大小为12kb，330行左右，总共38830行。转换后的html 大概为 12kb
```
finished parse 110 files
real	0m0.188s
user	0m0.169s
sys 	0m0.020s
```

当前单个文件处理耗时

```
read  file 56085   nanosecond
parse file 2279902 nanosecond
write file 126012  nanosecond
```
可见更多的时间消耗在 `parse` 这一步上，读写文件耗时之和／转换 = .079870538。在大文件的场景下，比例会更小。
因为 io 是阻塞资源，因此可以将转换与文件读写异步进行。一次文件转换的时间足够进行*10*次文件读写。
理想状况下，10个 `parse worker` 与 一个 `io worker` 异步，可以消去文件转换的峰值，可将一次完整流程耗时压至*2个文件读写时间*，性能提升约为10倍（前提是文件数 > 10）

## markdown 转换
逐个读取 `markdown` 目录下的所有文件，转换后，逐个进行存储。其次，需要关注文件的更新状况，此处计划根据git 来进行版本管理

## 优化计划

使用 goroutine 进行转换，开放 chan 给 io 使用。因为性能比为1/10，因此通过起 10 个 channel 可以间接的转换为 10 个 goroutine


# 优化记录

## 第一次优化
读取整个文件夹后，根据文件数暴力的开 goroutine 进行转换和存储

优化后的效果
```shell
110 files in total
time 56.526073m
```

```
3.14s of 3.30s total (95.15%)
Dropped 64 nodes (cum <= 0.02s)
Showing top 80 nodes out of 172 (cum >= 0.17s)
----------------------------------------------------------|-------------
      flat  flat%   sum%        cum   cum%   calls calls% + context 	 	 
----------------------------------------------------------|-------------
                                             0.94s   100% |   runtime.systemstack
     0.94s 28.48% 28.48%      0.94s 28.48%                | runtime.mach_semaphore_wait
----------------------------------------------------------|-------------
                                             0.15s   100% |   runtime.gcFlushBgCredit
     0.75s 22.73% 51.21%      0.75s 22.73%                | runtime.usleep
----------------------------------------------------------|-------------
                                             0.33s 94.29% |   os.OpenFile
                                             0.02s  5.71% |   os.(*File).Stat
     0.40s 12.12% 63.33%      0.40s 12.12%                | syscall.Syscall
----------------------------------------------------------|-------------
                                             0.19s   100% |   runtime.systemstack
     0.19s  5.76% 69.09%      0.19s  5.76%                | runtime.mach_semaphore_timedwait
----------------------------------------------------------|-------------
                                             0.03s   100% |   runtime.mallocgc
     0.09s  2.73% 71.82%      0.09s  2.73%                | runtime.memclr
----------------------------------------------------------|-------------
     0.08s  2.42% 74.24%      0.08s  2.42%                | runtime.mach_semaphore_signal
----------------------------------------------------------|-------------
                                             0.15s   100% |   bytes.makeSlice
     0.08s  2.42% 76.67%      0.55s 16.67%                | runtime.mallocgc
                                             0.08s 53.33% |   runtime.gcAssistAlloc
                                             0.04s 26.67% |   runtime.stkbucket
                                             0.03s 20.00% |   runtime.memclr
----------------------------------------------------------|-------------
                                             0.08s   100% |   strings.IndexAny
     0.08s  2.42% 79.09%      0.09s  2.73%                | runtime.stringiter2
----------------------------------------------------------|-------------
                                             0.10s 66.67% |   golang.org/x/net/html.EscapeString
                                             0.05s 33.33% |   golang.org/x/net/html.escape
     0.07s  2.12% 81.21%      0.15s  4.55%                | strings.IndexAny
                                             0.08s   100% |   runtime.stringiter2
----------------------------------------------------------|-------------
                                             0.01s   100% |   runtime.concatstrings
     0.06s  1.82% 83.03%      0.06s  1.82%                | runtime.memmove
----------------------------------------------------------|-------------
                                             0.04s   100% |   runtime.mallocgc
     0.04s  1.21% 84.24%      0.04s  1.21%                | runtime.stkbucket
----------------------------------------------------------|-------------
                                             0.03s   100% |   github.com/russross/blackfriday.(*Html).Smartypants
     0.03s  0.91% 85.15%      0.04s  1.21%                | github.com/russross/blackfriday.attrEscape
                                             0.01s   100% |   bytes.(*Buffer).Write
----------------------------------------------------------|-------------
                                             0.03s   100% |   runtime.systemstack
     0.03s  0.91% 86.06%      0.03s  0.91%                | nanotime
----------------------------------------------------------|-------------
                                             0.01s 33.33% |   github.com/microcosm-cc/bluemonday.(*Policy).sanitize
                                             0.01s 33.33% |   golang.org/x/net/html.(*Tokenizer).Token
                                             0.01s 33.33% |   runtime.pcvalue
     0.03s  0.91% 86.97%      0.03s  0.91%                | runtime.duffcopy
----------------------------------------------------------|-------------
                                             0.02s 50.00% |   github.com/microcosm-cc/bluemonday.(*Policy).AllowElements
                                             0.02s 50.00% |   github.com/microcosm-cc/bluemonday.(*attrPolicyBuilder).OnElements
     0.03s  0.91% 87.88%      0.05s  1.52%                | runtime.mapassign1
----------------------------------------------------------|-------------
                                             0.01s   100% |   runtime.systemstack
     0.03s  0.91% 88.79%      0.04s  1.21%                | runtime.scanobject
----------------------------------------------------------|-------------
                                             0.09s   100% |   runtime.systemstack
     0.03s  0.91% 89.70%      0.11s  3.33%                | runtime.scanstack
                                             0.07s   100% |   runtime.gentraceback
----------------------------------------------------------|-------------
                                             0.12s 70.59% |   bytes.(*Buffer).Write
                                             0.05s 29.41% |   bytes.(*Buffer).WriteString
     0.02s  0.61% 90.30%      0.17s  5.15%                | bytes.(*Buffer).grow
                                             0.15s   100% |   bytes.makeSlice
----------------------------------------------------------|-------------
                                             0.02s   100% |   regexp.Compile
     0.02s  0.61% 90.91%      0.02s  0.61%                | regexp/syntax.ranges.Less
----------------------------------------------------------|-------------
                                                 0   100% |   runtime.systemstack
     0.02s  0.61% 91.52%      0.18s  5.45%                | runtime.gcFlushBgCredit
                                             0.15s   100% |   runtime.usleep
----------------------------------------------------------|-------------
                                             0.02s   100% |   runtime.pcvalue
     0.02s  0.61% 92.12%      0.02s  0.61%                | runtime.step
----------------------------------------------------------|-------------
                                             0.46s   100% |   github.com/microcosm-cc/bluemonday.(*Policy).SanitizeBytes
     0.01s   0.3% 92.42%      0.46s 13.94%                | github.com/microcosm-cc/bluemonday.(*Policy).sanitize
                                             0.22s 51.16% |   golang.org/x/net/html.Token.String
                                             0.10s 23.26% |   golang.org/x/net/html.(*Tokenizer).Token
                                             0.07s 16.28% |   github.com/microcosm-cc/bluemonday.(*Policy).sanitizeAttrs
                                             0.03s  6.98% |   bytes.(*Buffer).WriteString
                                             0.01s  2.33% |   runtime.duffcopy
----------------------------------------------------------|-------------
                                             0.02s   100% |   github.com/russross/blackfriday.(*parser).table
     0.01s   0.3% 92.73%      0.02s  0.61%                | github.com/russross/blackfriday.(*parser).tableHeader
                                             0.01s   100% |   github.com/russross/blackfriday.(*parser).tableRow
----------------------------------------------------------|-------------
                                             0.10s   100% |   github.com/microcosm-cc/bluemonday.(*Policy).sanitize
     0.01s   0.3% 93.03%      0.10s  3.03%                | golang.org/x/net/html.(*Tokenizer).Token
                                             0.01s   100% |   runtime.duffcopy
----------------------------------------------------------|-------------
                                             0.03s   100% |   golang.org/x/net/html.Token.String
     0.01s   0.3% 93.33%      0.03s  0.91%                | runtime.concatstrings
                                             0.01s   100% |   runtime.memmove
----------------------------------------------------------|-------------
                                             0.08s   100% |   runtime.mallocgc
     0.01s   0.3% 93.64%      0.08s  2.42%                | runtime.gcAssistAlloc
                                             0.06s   100% |   runtime.systemstack
----------------------------------------------------------|-------------
                                             0.07s   100% |   runtime.scanstack
     0.01s   0.3% 93.94%      0.09s  2.73%                | runtime.gentraceback
                                             0.04s 66.67% |   runtime.pcvalue
                                             0.02s 33.33% |   runtime.scanblock
----------------------------------------------------------|-------------
                                             0.04s   100% |   runtime.gentraceback
     0.01s   0.3% 94.24%      0.04s  1.21%                | runtime.pcvalue
                                             0.02s 66.67% |   runtime.step
                                             0.01s 33.33% |   runtime.duffcopy
----------------------------------------------------------|-------------
                                             0.02s   100% |   runtime.gentraceback
     0.01s   0.3% 94.55%      0.02s  0.61%                | runtime.scanblock
----------------------------------------------------------|-------------
     0.01s   0.3% 94.85%      0.95s 28.79%                | runtime.semasleep
                                             0.94s   100% |   runtime.systemstack
----------------------------------------------------------|-------------
                                             0.94s 94.00% |   runtime.semasleep
                                             0.06s  6.00% |   runtime.gcAssistAlloc
     0.01s   0.3% 95.15%      1.48s 44.85%                | runtime.systemstack
                                             0.94s 74.31% |   runtime.mach_semaphore_wait
                                             0.19s 15.02% |   runtime.mach_semaphore_timedwait
                                             0.09s  7.11% |   runtime.scanstack
                                             0.03s  2.37% |   nanotime
                                             0.01s  0.85% |   runtime.scanobject
                                                 0  0.34% |   runtime.gcFlushBgCredit
----------------------------------------------------------|-------------
                                             0.04s 36.36% |   github.com/russross/blackfriday.(*Html).Smartypants
                                             0.04s 36.36% |   github.com/russross/blackfriday.expandTabs
                                             0.02s 18.18% |   github.com/russross/blackfriday.(*parser).listItem
                                             0.01s  9.09% |   github.com/russross/blackfriday.attrEscape
         0     0% 95.15%      0.12s  3.64%                | bytes.(*Buffer).Write
                                             0.12s   100% |   bytes.(*Buffer).grow
----------------------------------------------------------|-------------
                                             0.03s 60.00% |   github.com/microcosm-cc/bluemonday.(*Policy).sanitize
                                             0.01s 20.00% |   golang.org/x/net/html.Token.tagString
                                             0.01s 20.00% |   golang.org/x/net/html.escape
         0     0% 95.15%      0.05s  1.52%                | bytes.(*Buffer).WriteString
                                             0.05s   100% |   bytes.(*Buffer).grow
----------------------------------------------------------|-------------
                                             0.15s   100% |   bytes.(*Buffer).grow
         0     0% 95.15%      0.15s  4.55%                | bytes.makeSlice
                                             0.15s   100% |   runtime.mallocgc
----------------------------------------------------------|-------------
                                             0.02s   100% |   github.com/microcosm-cc/bluemonday.UGCPolicy
         0     0% 95.15%      0.02s  0.61%                | github.com/microcosm-cc/bluemonday.(*Policy).AllowElements
                                             0.02s   100% |   runtime.mapassign1
----------------------------------------------------------|-------------
                                             0.02s   100% |   github.com/microcosm-cc/bluemonday.UGCPolicy
         0     0% 95.15%      0.02s  0.61%                | github.com/microcosm-cc/bluemonday.(*Policy).AllowTables
                                             0.01s   100% |   github.com/microcosm-cc/bluemonday.(*attrPolicyBuilder).OnElements
----------------------------------------------------------|-------------
                                             0.46s   100% |   ink.mdParseStream
         0     0% 95.15%      0.46s 13.94%                | github.com/microcosm-cc/bluemonday.(*Policy).SanitizeBytes
                                             0.46s   100% |   github.com/microcosm-cc/bluemonday.(*Policy).sanitize
----------------------------------------------------------|-------------
                                             0.07s   100% |   github.com/microcosm-cc/bluemonday.(*Policy).sanitize
         0     0% 95.15%      0.07s  2.12%                | github.com/microcosm-cc/bluemonday.(*Policy).sanitizeAttrs
                                             0.02s 66.67% |   github.com/microcosm-cc/bluemonday.(*Policy).validURL
                                             0.01s 33.33% |   net/url.Parse
----------------------------------------------------------|-------------
                                             0.02s   100% |   github.com/microcosm-cc/bluemonday.(*Policy).sanitizeAttrs
         0     0% 95.15%      0.02s  0.61%                | github.com/microcosm-cc/bluemonday.(*Policy).validURL
                                             0.02s   100% |   net/url.Parse
----------------------------------------------------------|-------------
                                             0.01s 50.00% |   github.com/microcosm-cc/bluemonday.(*Policy).AllowTables
                                             0.01s 50.00% |   github.com/microcosm-cc/bluemonday.UGCPolicy
         0     0% 95.15%      0.02s  0.61%                | github.com/microcosm-cc/bluemonday.(*attrPolicyBuilder).OnElements
                                             0.02s   100% |   runtime.mapassign1
----------------------------------------------------------|-------------
                                             0.22s   100% |   ink.mdParseStream
         0     0% 95.15%      0.22s  6.67%                | github.com/microcosm-cc/bluemonday.UGCPolicy
                                             0.16s 76.19% |   regexp.MustCompile
                                             0.02s  9.52% |   github.com/microcosm-cc/bluemonday.(*Policy).AllowElements
                                             0.02s  9.52% |   github.com/microcosm-cc/bluemonday.(*Policy).AllowTables
                                             0.01s  4.76% |   github.com/microcosm-cc/bluemonday.(*attrPolicyBuilder).OnElements
----------------------------------------------------------|-------------
                                             0.03s   100% |   github.com/russross/blackfriday.(*parser).prefixHeader
         0     0% 95.15%      0.03s  0.91%                | github.com/russross/blackfriday.(*Html).Header
                                             0.02s   100% |   github.com/russross/blackfriday.(*parser).prefixHeader.func1
----------------------------------------------------------|-------------
                                             0.09s   100% |   github.com/russross/blackfriday.(*parser).list
         0     0% 95.15%      0.09s  2.73%                | github.com/russross/blackfriday.(*Html).List
                                             0.09s   100% |   github.com/russross/blackfriday.(*parser).list.func1
----------------------------------------------------------|-------------
                                             0.13s   100% |   github.com/russross/blackfriday.(*parser).inline
         0     0% 95.15%      0.13s  3.94%                | github.com/russross/blackfriday.(*Html).NormalText
                                             0.13s   100% |   github.com/russross/blackfriday.(*Html).Smartypants
----------------------------------------------------------|-------------
                                             0.05s   100% |   github.com/russross/blackfriday.(*parser).renderParagraph
         0     0% 95.15%      0.05s  1.52%                | github.com/russross/blackfriday.(*Html).Paragraph
                                             0.05s   100% |   github.com/russross/blackfriday.(*parser).renderParagraph.func1
----------------------------------------------------------|-------------
                                             0.13s   100% |   github.com/russross/blackfriday.(*Html).NormalText
         0     0% 95.15%      0.13s  3.94%                | github.com/russross/blackfriday.(*Html).Smartypants
                                             0.04s 57.14% |   bytes.(*Buffer).Write
                                             0.03s 42.86% |   github.com/russross/blackfriday.attrEscape
----------------------------------------------------------|-------------
                                             0.24s 88.89% |   github.com/russross/blackfriday.secondPass
                                             0.03s 11.11% |   github.com/russross/blackfriday.(*parser).quote
         0     0% 95.15%      0.24s  7.27%                | github.com/russross/blackfriday.(*parser).block
                                             0.09s 39.13% |   github.com/russross/blackfriday.(*parser).list
                                             0.06s 26.09% |   github.com/russross/blackfriday.(*parser).paragraph
                                             0.05s 21.74% |   github.com/russross/blackfriday.(*parser).table
                                             0.03s 13.04% |   github.com/russross/blackfriday.(*parser).prefixHeader
----------------------------------------------------------|-------------
                                             0.07s 46.67% |   github.com/russross/blackfriday.(*parser).listItem
                                             0.05s 33.33% |   github.com/russross/blackfriday.(*parser).renderParagraph.func1
                                             0.02s 13.33% |   github.com/russross/blackfriday.(*parser).prefixHeader.func1
                                             0.01s  6.67% |   github.com/russross/blackfriday.(*parser).tableRow
         0     0% 95.15%      0.15s  4.55%                | github.com/russross/blackfriday.(*parser).inline
                                             0.13s   100% |   github.com/russross/blackfriday.(*Html).NormalText
----------------------------------------------------------|-------------
                                             0.09s   100% |   github.com/russross/blackfriday.(*parser).block
         0     0% 95.15%      0.09s  2.73%                | github.com/russross/blackfriday.(*parser).list
                                             0.09s   100% |   github.com/russross/blackfriday.(*Html).List
----------------------------------------------------------|-------------
                                             0.09s   100% |   github.com/russross/blackfriday.(*Html).List
         0     0% 95.15%      0.09s  2.73%                | github.com/russross/blackfriday.(*parser).list.func1
                                             0.09s   100% |   github.com/russross/blackfriday.(*parser).listItem
----------------------------------------------------------|-------------
                                             0.09s   100% |   github.com/russross/blackfriday.(*parser).list.func1
         0     0% 95.15%      0.09s  2.73%                | github.com/russross/blackfriday.(*parser).listItem
                                             0.07s 77.78% |   github.com/russross/blackfriday.(*parser).inline
                                             0.02s 22.22% |   bytes.(*Buffer).Write
----------------------------------------------------------|-------------
                                             0.06s   100% |   github.com/russross/blackfriday.(*parser).block
         0     0% 95.15%      0.06s  1.82%                | github.com/russross/blackfriday.(*parser).paragraph
                                             0.06s   100% |   github.com/russross/blackfriday.(*parser).renderParagraph
----------------------------------------------------------|-------------
                                             0.03s   100% |   github.com/russross/blackfriday.(*parser).block
         0     0% 95.15%      0.03s  0.91%                | github.com/russross/blackfriday.(*parser).prefixHeader
                                             0.03s   100% |   github.com/russross/blackfriday.(*Html).Header
----------------------------------------------------------|-------------
                                             0.02s   100% |   github.com/russross/blackfriday.(*Html).Header
         0     0% 95.15%      0.02s  0.61%                | github.com/russross/blackfriday.(*parser).prefixHeader.func1
                                             0.02s   100% |   github.com/russross/blackfriday.(*parser).inline
----------------------------------------------------------|-------------
         0     0% 95.15%      0.03s  0.91%                | github.com/russross/blackfriday.(*parser).quote
                                             0.03s   100% |   github.com/russross/blackfriday.(*parser).block
----------------------------------------------------------|-------------
                                             0.06s   100% |   github.com/russross/blackfriday.(*parser).paragraph
         0     0% 95.15%      0.06s  1.82%                | github.com/russross/blackfriday.(*parser).renderParagraph
                                             0.05s   100% |   github.com/russross/blackfriday.(*Html).Paragraph
----------------------------------------------------------|-------------
                                             0.05s   100% |   github.com/russross/blackfriday.(*Html).Paragraph
         0     0% 95.15%      0.05s  1.52%                | github.com/russross/blackfriday.(*parser).renderParagraph.func1
                                             0.05s   100% |   github.com/russross/blackfriday.(*parser).inline
----------------------------------------------------------|-------------
                                             0.05s   100% |   github.com/russross/blackfriday.(*parser).block
         0     0% 95.15%      0.05s  1.52%                | github.com/russross/blackfriday.(*parser).table
                                             0.02s 50.00% |   github.com/russross/blackfriday.(*parser).tableHeader
                                             0.02s 50.00% |   github.com/russross/blackfriday.(*parser).tableRow
----------------------------------------------------------|-------------
                                             0.02s 66.67% |   github.com/russross/blackfriday.(*parser).table
                                             0.01s 33.33% |   github.com/russross/blackfriday.(*parser).tableHeader
         0     0% 95.15%      0.03s  0.91%                | github.com/russross/blackfriday.(*parser).tableRow
                                             0.01s   100% |   github.com/russross/blackfriday.(*parser).inline
----------------------------------------------------------|-------------
                                             0.29s   100% |   ink.mdParseStream
         0     0% 95.15%      0.29s  8.79%                | github.com/russross/blackfriday.MarkdownCommon
                                             0.29s   100% |   github.com/russross/blackfriday.MarkdownOptions
----------------------------------------------------------|-------------
                                             0.29s   100% |   github.com/russross/blackfriday.MarkdownCommon
         0     0% 95.15%      0.29s  8.79%                | github.com/russross/blackfriday.MarkdownOptions
                                             0.24s 85.71% |   github.com/russross/blackfriday.secondPass
                                             0.04s 14.29% |   github.com/russross/blackfriday.firstPass
----------------------------------------------------------|-------------
                                             0.04s   100% |   github.com/russross/blackfriday.firstPass
         0     0% 95.15%      0.04s  1.21%                | github.com/russross/blackfriday.expandTabs
                                             0.04s   100% |   bytes.(*Buffer).Write
----------------------------------------------------------|-------------
                                             0.04s   100% |   github.com/russross/blackfriday.MarkdownOptions
         0     0% 95.15%      0.04s  1.21%                | github.com/russross/blackfriday.firstPass
                                             0.04s   100% |   github.com/russross/blackfriday.expandTabs
----------------------------------------------------------|-------------
                                             0.24s   100% |   github.com/russross/blackfriday.MarkdownOptions
         0     0% 95.15%      0.24s  7.27%                | github.com/russross/blackfriday.secondPass
                                             0.24s   100% |   github.com/russross/blackfriday.(*parser).block
----------------------------------------------------------|-------------
                                             0.15s   100% |   golang.org/x/net/html.Token.String
         0     0% 95.15%      0.15s  4.55%                | golang.org/x/net/html.EscapeString
                                             0.10s 71.43% |   strings.IndexAny
                                             0.04s 28.57% |   golang.org/x/net/html.escape
----------------------------------------------------------|-------------
                                             0.22s   100% |   github.com/microcosm-cc/bluemonday.(*Policy).sanitize
         0     0% 95.15%      0.22s  6.67%                | golang.org/x/net/html.Token.String
                                             0.15s 68.18% |   golang.org/x/net/html.EscapeString
                                             0.04s 18.18% |   golang.org/x/net/html.Token.tagString
                                             0.03s 13.64% |   runtime.concatstrings
----------------------------------------------------------|-------------
                                             0.04s   100% |   golang.org/x/net/html.Token.String
         0     0% 95.15%      0.04s  1.21%                | golang.org/x/net/html.Token.tagString
                                             0.02s 66.67% |   golang.org/x/net/html.escape
                                             0.01s 33.33% |   bytes.(*Buffer).WriteString
----------------------------------------------------------|-------------
                                             0.04s 66.67% |   golang.org/x/net/html.EscapeString
                                             0.02s 33.33% |   golang.org/x/net/html.Token.tagString
         0     0% 95.15%      0.06s  1.82%                | golang.org/x/net/html.escape
                                             0.05s 83.33% |   strings.IndexAny
                                             0.01s 16.67% |   bytes.(*Buffer).WriteString
----------------------------------------------------------|-------------
         0     0% 95.15%      1.33s 40.30%                | ink.mdParseStream
                                             0.46s 34.59% |   github.com/microcosm-cc/bluemonday.(*Policy).SanitizeBytes
                                             0.36s 27.07% |   io/ioutil.WriteFile
                                             0.29s 21.80% |   github.com/russross/blackfriday.MarkdownCommon
                                             0.22s 16.54% |   github.com/microcosm-cc/bluemonday.UGCPolicy
----------------------------------------------------------|-------------
         0     0% 95.15%      0.03s  0.91%                | ink.mdReadStream
                                             0.03s   100% |   io/ioutil.ReadFile
----------------------------------------------------------|-------------
                                             0.03s   100% |   ink.mdReadStream
         0     0% 95.15%      0.03s  0.91%                | io/ioutil.ReadFile
                                             0.02s   100% |   os.(*File).Stat
----------------------------------------------------------|-------------
                                             0.36s   100% |   ink.mdParseStream
         0     0% 95.15%      0.36s 10.91%                | io/ioutil.WriteFile
                                             0.33s 91.67% |   os.OpenFile
                                             0.03s  8.33% |   os.(*File).Write
----------------------------------------------------------|-------------
                                             0.02s 66.67% |   github.com/microcosm-cc/bluemonday.(*Policy).validURL
                                             0.01s 33.33% |   github.com/microcosm-cc/bluemonday.(*Policy).sanitizeAttrs
         0     0% 95.15%      0.03s  0.91%                | net/url.Parse
                                             0.03s   100% |   net/url.parse
----------------------------------------------------------|-------------
                                             0.03s   100% |   net/url.Parse
         0     0% 95.15%      0.03s  0.91%                | net/url.parse
----------------------------------------------------------|-------------
                                             0.02s   100% |   io/ioutil.ReadFile
         0     0% 95.15%      0.02s  0.61%                | os.(*File).Stat
                                             0.02s   100% |   syscall.Syscall
----------------------------------------------------------|-------------
                                             0.03s   100% |   io/ioutil.WriteFile
         0     0% 95.15%      0.03s  0.91%                | os.(*File).Write
                                             0.03s   100% |   os.(*File).write
----------------------------------------------------------|-------------
                                             0.03s   100% |   os.(*File).Write
         0     0% 95.15%      0.03s  0.91%                | os.(*File).write
----------------------------------------------------------|-------------
                                             0.33s   100% |   io/ioutil.WriteFile
         0     0% 95.15%      0.33s 10.00%                | os.OpenFile
                                             0.33s   100% |   syscall.Syscall
----------------------------------------------------------|-------------
                                             0.17s   100% |   regexp.MustCompile
         0     0% 95.15%      0.17s  5.15%                | regexp.Compile
                                             0.02s   100% |   regexp/syntax.ranges.Less
----------------------------------------------------------|-------------
                                             0.16s   100% |   github.com/microcosm-cc/bluemonday.UGCPolicy
         0     0% 95.15%      0.17s  5.15%                | regexp.MustCompile
                                             0.17s   100% |   regexp.Compile
----------------------------------------------------------|-------------
```

### 结论
喂不饱，过多的等待

## 第二次优化

* 使用 sync.Mutex 代替 channel
* 减少 GC

### 优化结果

随着文件数增长，性能越好

```

parsing 123 files
time 54.499284ms

parsing 123 files
time 68.931115ms

parsing 123 files
time 55.680802ms

parsing 492 files
time 221.63780400000002ms

parsing 492 files
time 233.301539ms

parsing 492 files
time 197.314185ms

parsing 492 files
time 196.232665ms
```

### 492 个文件转换分析

```
Showing top 80 nodes out of 161 (cum >= 0.90s)
----------------------------------------------------------|-------------
      flat  flat%   sum%        cum   cum%   calls calls% + context 	 	 
----------------------------------------------------------|-------------
                                             0.09s   100% |   bytes.(*Buffer).ReadFrom
     3.10s 32.98% 32.98%      3.10s 32.98%                | syscall.Syscall
----------------------------------------------------------|-------------
                                             1.41s   100% |   runtime.systemstack
     1.41s 15.00% 47.98%      1.41s 15.00%                | runtime.mach_semaphore_wait
----------------------------------------------------------|-------------
                                             0.90s   100% |   runtime.lock
     1.33s 14.15% 62.13%      1.33s 14.15%                | runtime.usleep
----------------------------------------------------------|-------------
                                             0.26s 96.30% |   strings.IndexAny
                                             0.01s  3.70% |   strings.Map
     0.24s  2.55% 64.68%      0.27s  2.87%                | runtime.stringiter2
----------------------------------------------------------|-------------
                                             0.43s   100% |   golang.org/x/net/html.EscapeString
     0.24s  2.55% 67.23%      0.50s  5.32%                | strings.IndexAny
                                             0.26s   100% |   runtime.stringiter2
----------------------------------------------------------|-------------
                                             0.19s   100% |   runtime.systemstack
     0.19s  2.02% 69.26%      0.19s  2.02%                | runtime.mach_semaphore_timedwait
----------------------------------------------------------|-------------
                                             0.31s 32.98% |   runtime.newobject
                                             0.30s 31.91% |   runtime.makeslice
                                             0.22s 23.40% |   runtime.rawstring
                                             0.11s 11.70% |   runtime.growslice
     0.19s  2.02% 71.28%      0.94s 10.00%                | runtime.mallocgc
                                             0.07s   100% |   runtime.memclr
----------------------------------------------------------|-------------
                                             0.14s   100% |   golang.org/x/net/html.(*Tokenizer).Next
     0.18s  1.91% 73.19%      0.18s  1.91%                | golang.org/x/net/html.(*Tokenizer).readByte
----------------------------------------------------------|-------------
     0.17s  1.81% 75.00%      0.17s  1.81%                | runtime.mach_semaphore_signal
----------------------------------------------------------|-------------
                                             0.03s 37.50% |   bytes.(*Buffer).Write
                                             0.03s 37.50% |   bytes.(*Buffer).WriteString
                                             0.01s 12.50% |   fmt.(*pp).doPrintf
                                             0.01s 12.50% |   runtime.concatstrings
     0.12s  1.28% 76.28%      0.12s  1.28%                | runtime.memmove
----------------------------------------------------------|-------------
                                             0.07s   100% |   runtime.mallocgc
     0.11s  1.17% 77.45%      0.11s  1.17%                | runtime.memclr
----------------------------------------------------------|-------------
                                             0.07s   100% |   runtime.systemstack
     0.10s  1.06% 78.51%      0.10s  1.06%                | nanotime
----------------------------------------------------------|-------------
                                             0.16s 81.11% |   github.com/microcosm-cc/bluemonday.(*Policy).sanitize
                                             0.04s 18.89% |   golang.org/x/net/html.Token.tagString
     0.08s  0.85% 79.36%      0.26s  2.77%                | bytes.(*Buffer).WriteString
                                             0.15s 83.33% |   bytes.(*Buffer).grow
                                             0.03s 16.67% |   runtime.memmove
----------------------------------------------------------|-------------
                                             0.04s 57.14% |   runtime.scanobject
                                             0.03s 42.86% |   runtime.shade
     0.08s  0.85% 80.21%      0.08s  0.85%                | runtime.heapBitsForObject
----------------------------------------------------------|-------------
                                             0.27s   100% |   github.com/microcosm-cc/bluemonday.(*Policy).sanitize
     0.07s  0.74% 80.96%      0.27s  2.87%                | golang.org/x/net/html.(*Tokenizer).Next
                                             0.14s   100% |   golang.org/x/net/html.(*Tokenizer).readByte
----------------------------------------------------------|-------------
                                             0.07s 87.50% |   github.com/microcosm-cc/bluemonday.(*Policy).sanitize
                                             0.01s 12.50% |   github.com/microcosm-cc/bluemonday.(*Policy).sanitizeAttrs
     0.07s  0.74% 81.70%      0.09s  0.96%                | runtime.mapaccess2_faststr
----------------------------------------------------------|-------------
                                             0.37s   100% |   github.com/microcosm-cc/bluemonday.(*Policy).sanitize
     0.06s  0.64% 82.34%      0.37s  3.94%                | golang.org/x/net/html.(*Tokenizer).Token
                                             0.06s 35.29% |   golang.org/x/net/html.(*Tokenizer).TagName
                                             0.06s 35.29% |   golang.org/x/net/html/atom.Lookup
                                             0.03s 17.65% |   runtime.growslice
                                             0.02s 11.76% |   runtime.duffcopy
----------------------------------------------------------|-------------
                                             0.06s   100% |   golang.org/x/net/html.(*Tokenizer).Token
     0.06s  0.64% 82.98%      0.06s  0.64%                | golang.org/x/net/html/atom.Lookup
----------------------------------------------------------|-------------
                                             0.02s 66.67% |   golang.org/x/net/html.(*Tokenizer).Token
                                             0.01s 33.33% |   github.com/microcosm-cc/bluemonday.(*Policy).sanitizeAttrs
     0.06s  0.64% 83.62%      0.06s  0.64%                | runtime.duffcopy
----------------------------------------------------------|-------------
     0.06s  0.64% 84.26%      0.11s  1.17%                | runtime.scanobject
                                             0.04s   100% |   runtime.heapBitsForObject
----------------------------------------------------------|-------------
                                             0.13s   100% |   golang.org/x/net/html.Token.String
     0.05s  0.53% 84.79%      0.13s  1.38%                | runtime.concatstrings
                                             0.07s 87.50% |   runtime.rawstringtmp
                                             0.01s 12.50% |   runtime.memmove
----------------------------------------------------------|-------------

```

## 第三次优化

* 减少 channel 带来的锁

### 优化结果
```
parsing 492 files
time 210.924972ms
neodeMacBook-Pro:ink neo$ ./ink
parsing 492 files
time 188.60070000000002ms
neodeMacBook-Pro:ink neo$ ./ink
parsing 492 files
time 186.06235999999998ms
neodeMacBook-Pro:ink neo$ ./ink
parsing 492 files
time 183.858968ms
```

## 第四次优化

* 减少 GC

### 优化结果
```
neodeMacBook-Pro:ink neo$ ./ink 
parsing 500 files
time 171.346057ms
neodeMacBook-Pro:ink neo$ ./ink 
parsing 500 files
time 178.415718ms
neodeMacBook-Pro:ink neo$ ./ink 
parsing 500 files
time 158.84637ms
neodeMacBook-Pro:ink neo$ ./ink 
parsing 500 files
time 173.826369ms
neodeMacBook-Pro:ink neo$ ./ink 
parsing 500 files
time 181.290882ms
neodeMacBook-Pro:ink neo$ ./ink 
parsing 500 files
time 159.058857ms
```

### 优化结果
```
neodeMacBook-Pro:ink neo$ ./ink 
parsing 500 files
time 166.6263ms
neodeMacBook-Pro:ink neo$ ./ink 
parsing 500 files
time 179.101538ms
neodeMacBook-Pro:ink neo$ ./ink 
parsing 500 files
time 170.17852299999998ms
neodeMacBook-Pro:ink neo$ ./ink 
parsing 500 files
time 169.913514ms
neodeMacBook-Pro:ink neo$ ./ink 
parsing 500 files
time 190.60826ms
neodeMacBook-Pro:ink neo$ ./ink 
parsing 500 files
time 168.667486ms
neodeMacBook-Pro:ink neo$ 

```
