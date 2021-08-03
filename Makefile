.PHONY:
.SILENT:

build:
	go build -o ./runtime/.bin/bot cmd/bot/main.go

run: build
	./runtime/.bin/bot

build-image:
	docker build -t pocketer-telegram-bot:v0.1 .

start-container:
	docker run --name pocketer-bot -p 80:80 --env-file .env pocketer-telegram-bot:v0.1
