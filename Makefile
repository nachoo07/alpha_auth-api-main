build:
	@docker-compose build --no-cache --force-rm --pull

down:
	@docker-compose down

run: down
	@docker-compose up --build

run_local:
	@eval $$(egrep -v '^#' variables.env | xargs) APP_PATH=$$PWD go run cmd/api/*.go

dev:
	@docker-compose up -d testlocal

fs:
	@make down
	@make dev
	@sleep 30
	@make migrate

lint:
	@eval $$(egrep -v '^#' variables.env | xargs) APP_PATH=$$PWD go list -e -compiled -test=true -export=false -deps=true -find=false ./... > /dev/null
	@eval $$(egrep -v '^#' variables.env | xargs) APP_PATH=$$PWD golangci-lint run --config=.code_quality/.golangci.yml --new-from-rev=HEAD~1 --fix

test: fs
	@docker-compose build fury_shp-travel-management-api --no-cache --force-rm
	@docker-compose run --rm fury_shp-travel-management-api /commands/test.sh


migrate:
	@eval $$(egrep -v '^#' variables.env | xargs) APP_PATH=$$PWD go run cmd/migration/main.go UP

cleanup:
	@find . -type d -name mocks -exec rm -rf {} \;

mocks: cleanup
	@eval $$(egrep -v '^#' variables.env | xargs) APP_PATH=$$PWD go generate ./...

test_local:
	@eval $$(egrep -v '^#' variables.env | xargs) APP_PATH=$$PWD go test ./... -count 1 -tags=integration -cover -p=1

test_up: down
	@docker-compose up -d

test_run:
	@docker run -v $$PWD:/app fury_shp-travel-management-api /core/test.sh
