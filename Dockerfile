
# STEP 1 build executable binary
FROM golang:alpine as builder
COPY . $GOPATH/src/github.com/r4ulcl/bntoolkit

WORKDIR $GOPATH/src/github.com/r4ulcl/bntoolkit
#get dependancies
RUN apk -U add alpine-sdk
RUN go get -d -v
RUN go get github.com/anacrolix/utp
#build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-w -s' -o /go/bin/bntoolkit

# STEP 2 build a small image
# start from scratch
FROM scratch
#GOPATH doesn-t exists in scratch
ENV GOPATH='/go' 

# Copy our static executable
COPY --from=builder /$GOPATH/bin/bntoolkit /$GOPATH/bin/bntoolkit
#Copy all code (sql, conf files, etc)
COPY --from=builder /$GOPATH/src/github.com/r4ulcl/bntoolkit/sql.sql /$GOPATH/src/github.com/r4ulcl/bntoolkit/sql.sql
COPY --from=builder /$GOPATH/src/github.com/r4ulcl/bntoolkit/configFile.toml /$GOPATH/src/github.com/r4ulcl/bntoolkit/configFile.toml
COPY --from=builder /$GOPATH/src/github.com/r4ulcl/bntoolkit/doc /$GOPATH/src/github.com/r4ulcl/bntoolkit/doc

ENTRYPOINT ["/go/bin/bntoolkit"]
