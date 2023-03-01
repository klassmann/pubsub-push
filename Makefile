
all: clean win_64 win_32 linux_64 linux_32 darwin_64

win_64:
	env GOOS=windows GOARCH=amd64 go build -o "dist/pubsub-push.exe" ./cmd/pubsub-push
	zip -T -j -9 "dist/pubsub-push_$(shell cat VERSION)_$@.zip" dist/pubsub-push.exe
	rm -f dist/pubsub-push.exe

win_32:
	env GOOS=windows GOARCH=386 go build -o "dist/pubsub-push.exe" ./cmd/pubsub-push
	zip -T -j -9 "dist/pubsub-push_$(shell cat VERSION)_$@.zip" dist/pubsub-push.exe
	rm -f dist/pubsub-push.exe

linux_64:
	env GOOS=linux GOARCH=amd64 go build -o dist/pubsub-push ./cmd/pubsub-push
	gzip dist/pubsub-push -c > "dist/pubsub-push_$(shell cat VERSION)_$@.gz"
	rm -f dist/pubsub-push

linux_32:
	env GOOS=linux GOARCH=386 go build -o dist/pubsub-push ./cmd/pubsub-push
	gzip dist/pubsub-push -c > "dist/pubsub-push_$(shell cat VERSION)_$@.gz"
	rm -f dist/pubsub-push

darwin_64:
	env GOOS=darwin GOARCH=amd64 go build -o dist/pubsub-push ./cmd/pubsub-push
	gzip dist/pubsub-push -c > "dist/pubsub-push_$(shell cat VERSION)_$@.gz"
	rm -f dist/pubsub-push

clean:
	rm -rf dist/*

