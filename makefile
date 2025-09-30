
dev:
	git checkout development

prod:
	git checkout master && git merge development

bench:
	go test -run=xxxx -bench=. -benchtime=5s -cpuprofile=cpu.pprof -benchmem

seebench:
	go tool pprof cpu.pprof