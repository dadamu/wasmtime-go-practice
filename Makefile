build-wasm:
	@echo "Building WebAssembly..."
	@cd script && yarn && yarn asbuild:release
	@cp script/build/release.wasm ./release.wasm

start-docker-test: stop-docker-test
	@echo "Starting Docker container..."
	@docker run --name test -e POSTGRES_USER=test -e POSTGRES_PASSWORD=test -e POSTGRES_DB=test -d -p 5432:5432 postgres
.PHONY: start-docker-test

stop-docker-test:
	@echo "Stopping Docker container..."
	@docker stop test || true && docker rm test || true
.PHONY: stop-docker-test