projectName:=$(shell gcloud config get-value project)
clusterName:=greeting-cluster
zone:=asia-northeast1-c
services:=service-a service-b

.PHONY: init
init:
	gcloud container clusters create $(clusterName) --zone=$(zone) --num-nodes=1 --preemptible
	gcloud container clusters get-credentials $(clusterName) --zone=$(zone)
	helm init

.PHONY: build
build:
	@for service in $(services); do \
		docker build ./$$service -f $$service/Dockerfile -t $$service:stable;\
		docker tag $$service:stable gcr.io/$(projectName)/$$service:stable;\
		docker push gcr.io/$(projectName)/$$service:stable;\
	done

.PHONY: deploy
deploy:
	@for service in $(services); do \
		helm install --name=$$service ./$$service/$$service --set image.project=$(projectName)\
	done