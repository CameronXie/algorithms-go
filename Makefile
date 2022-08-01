outputDir=${PWD}/_dist
testOutputDir=${outputDir}/tests

# Docker
up:
	@docker compose up --build -d

down:
	@docker compose down -v

test:
	@mkdir -p ${testOutputDir}
	@go clean -testcache
	@go test \
        -cover \
        -coverprofile=cp.out \
        -outputdir=${testOutputDir} \
        -race \
        -v \
        -failfast \
        ./...
	@go tool cover -html=${testOutputDir}/cp.out -o ${testOutputDir}/cp.html

lint:
	@golangci-lint run ./... -v
