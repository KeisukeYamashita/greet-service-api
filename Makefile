projectName:=$(shell gcloud config get-value project)
clusterName:=greeting-cluster
zone:=asia-northeast1-c
services:=serviceA serviceB

.PHONY: init
init:
	gcloud container clusters create $(clusterName) --zone=$(zone) --num-nodes=1 --preemptible
	gcloud container clusters get-credentials $(clusterName) --zone=$(zone)

.PHONY: build
build:
	@for service in $(services); do \
		docker build ./$$service -t $$service;\ 
		docker tag $$service gcr.io/$$projectName/$$service ;\
		docker push gcr.io/$$projectName/$$service;\
	done

.PHONY: deploy
deploy:
	@for service in $(services); do \
		helm install --name=$$service ./$$service;\
	done