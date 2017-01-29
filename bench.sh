go test -bench="." -cpuprofile=prof.cpu
go tool pprof ink.test prof.cpu
