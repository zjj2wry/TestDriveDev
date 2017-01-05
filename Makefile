SHELL:=/bin/bash

test:
	go test -cover ./...

#采集测试覆盖率，是否覆盖
collect_bool:
	rm all.out coverage.out
	for pkg in $$(go list ./...);do\
				echo "1";\
				go test -coverprofile=coverage.out $${pkg} || exit $$?;\
				if [ -f coverage.out ] ; then \
				sed -i '1d' coverage.out ;\
				cat coverage.out >> all.out ;\
				fi ; \
	done;\
	sed -i '1 i\mode: set' all.out 

#采集测试覆盖率，显示次数，可以得到代码块的调用次数
collect_count:
	rm all.out coverage.out
	for pkg in $$(go list ./...);do\
				echo "1";\
				go test -coverprofile=coverage.out -covermode=count $${pkg} || exit $$?;\
				if [ -f coverage.out ] ; then \
				sed -i '1d' coverage.out ;\
				cat coverage.out >> all.out ;\
				fi ; \
	done;\
	sed -i '1 i\mode: count' all.out

#将数据用html显示，持久化成文件形式 
html:
	go tool cover -html=all.out -o coverage.html

#将数据用html显示，持久化成文件形式 
html_open:
	go tool cover -html=all.out	

#函数形式显示
func:
	go tool cover -func=all.out

.PHONY: collect
