.PHONY: local

build:
	sam build

local:
	@make build;\
	sam local start-api --env-vars ./local/local.envs.json --parameter-overrides $(cat ./local/local.sam-params)
