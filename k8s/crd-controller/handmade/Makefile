

run: crd build dockerbuild apply

run_xiaoming: crd build dockerbuild apply_xiaoming

crd:
	kubectl apply -f deploy/crd.definition.yaml
	kubectl apply -f deploy/crd.instance.yaml

build:
	CGO_ENABLED=0 go build -o controller ./

dockerbuild:
	docker build . -t github.com/trainyao/playground/crd-controller/handmade:latest

apply: apply_xiaoming apply_xiaowang apply_xiaojie

apply_xiaoming:
	kubectl apply -f deploy/controller.xiaoming.yaml

restart_xiaoming:
	kubectl rollout restart deployment xiaoming

apply_xiaowang:
	kubectl apply -f deploy/controller.xiaowang.yaml

restart_xiaowang:
	kubectl rollout restart deployment xiaowang

apply_xiaojie:
	kubectl apply -f deploy/controller.xiaojie.yaml

restart_xiaojie:
	kubectl rollout restart deployment xiaojie

clear:
	kubectl delete -f deploy/
	rm -f controller
	
.PHONY: run_xiaoming build dockerbuild apply apply_xiaoming apply_xiaowang apply_xiaojie clear
