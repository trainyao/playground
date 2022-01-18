

run: build dockerbuild apply

run_xiaoming: build dockerbuild apply_xiaoming

build:
	CGO_ENABLED=0 go build -o controller ./

dockerbuild:
	docker build . -t github.com/trainyao/playground/crd-controller/handmade:latest

apply: apply_xiaoming apply_xiaowang apply_xiaojie

apply_xiaoming:
	kubectl apply -f deploy/controller.xiaoming.yaml
	kubectl rollout restart deployment xiaoming

apply_xiaowang:
	kubectl apply -f deploy/controller.xiaowang.yaml
	kubectl rollout restart deployment xiaowang

apply_xiaojie:
	kubectl apply -f deploy/controller.xiaojie.yaml
	kubectl rollout restart deployment xiaojie
	
.PHONY: run_xiaoming build dockerbuild apply apply_xiaoming apply_xiaowang apply_xiaojie
