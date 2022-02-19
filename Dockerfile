FROM golang:1.18rc1-buster

RUN apt update -y && \
	apt install -y build-essential cmake && \
	git clone --recursive https://github.com/WebAssembly/wabt /wabt && \
	cd /wabt && \
	git submodule update --init && \
	mkdir build && \
 	cd build && \
 	cmake .. && \
 	cmake --build . && \
	cp ../bin/* /usr/local/bin
	
# devcontainer golang settings
# RUN go install -v \
# 	github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest \
# 	github.com/ramya-rao-a/go-outline@latest \
# 	github.com/cweill/gotests/gotests@latest \
# 	github.com/fatih/gomodifytags@latest \
# 	github.com/josharian/impl@latest \
# 	github.com/haya14busa/goplay/cmd/goplay@latest \
# 	github.com/go-delve/delve/cmd/dlv@latest \
# 	honnef.co/go/tools/cmd/staticcheck@latest \
# 	golang.org/x/tools/gopls@latest

WORKDIR /

CMD ["/bin/bash"]
