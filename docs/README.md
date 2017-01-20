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
