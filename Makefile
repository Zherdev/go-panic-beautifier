docker-build:
	docker build --tag go-panic-beautifier:1.0.0 .

docker-run: docker-build
	docker run -p 80:80 --name go-panic-beautifier go-panic-beautifier:1.0.0

docker-clean:
	docker container kill go-panic-beautifier | true
	docker container rm go-panic-beautifier | true
	docker images | grep go-panic-beautifier | tr -s ' ' | cut -d ' ' -f 2 | xargs -I {} docker rmi go-panic-beautifier:{} | true

docker-rm:
	docker container kill go-panic-beautifier | true
	docker container rm go-panic-beautifier | true
