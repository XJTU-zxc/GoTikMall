.PHONY: all
all: help

default: help

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: gen-frontend
gen-frontend: ## gen frontend code
	@cd app/frontend && cwgo server -I ../../idl --type HTTP --module github.com/XJTU-zxc/MyCloudWeGo/app/frontend --service frontend  --idl ../../idl/frontend/home.proto

.PHONY: gen-client
gen-client: ## gen client code of {svc}. example: make gen-client svc=product
	@cd rpc_gen && cwgo client --type RPC --service ${svc} --module github.com/XJTU-zxc/GoTikMall/rpc_gen --I ../idl  --idl ../idl/${svc}.proto

.PHONY: gen-server
gen-server: ## gen service code of {svc}. example: make gen-server svc=product
	@cd app/${svc} && cwgo server --type RPC --service ${svc} --module github.com/XJTU-zxc/GoTikMall/app/${svc} --pass "-use github.com/XJTU-zxc/GoTikMall/rpc_gen/kitex_gen" --I ../../idl  --idl ../../idl/${svc}.proto



