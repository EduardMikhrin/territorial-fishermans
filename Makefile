protogen:
	cd proto && \
	buf generate application --template=./templates/application.yaml --config=buf.yaml && \
	buf generate api --template=./templates/api.yaml --config=buf.yaml && \
	buf generate client --template=./templates/client.yaml --config=buf.yaml && \
    buf generate location --template=./templates/location.yaml --config=buf.yaml && \
    buf generate application_status --template=./templates/application_status.yaml --config=buf.yaml && \
    rm ../internal/api/types/api.pb.go


install:
	export CGO_ENABLED=1
	rm -f $(GOPATH)/bin/relayer-svc
	go build -o $(GOPATH)/bin

run: install
	kursach service run all