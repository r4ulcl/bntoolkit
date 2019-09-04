
# STEP 1 build executable binary
FROM golang:alpine as builder
COPY . $GOPATH/src/github.com/RaulCalvoLaorden/bntoolkit
COPY ../torrent $GOPATH/src/github.com/RaulCalvoLaorden/torrent

WORKDIR $GOPATH/src/github.com/RaulCalvoLaorden/bntoolkit
#get dependancies
#you can also use dep
RUN apk -U add alpine-sdk
RUN go get -d -v
RUN go get github.com/anacrolix/utp
#build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-w -s' -o /go/bin/bntoolkit
# STEP 2 build a small image
# start from scratch
FROM scratch
# Copy our static executable
COPY --from=builder /go/bin/bntoolkit /go/bin/bntoolkit
ENTRYPOINT ["/go/bin/bntoolkit"]


#FROM scratch
#COPY findTorrent /app/
#WORKDIR /app
#ENTRYPOINT ["./findTorrent"]

#docker stop $(docker ps -aq) ; docker rm $(docker ps -aq) ;docker rmi $(docker images -q)
#sudo docker run -d -p 5432:5432 --mount type=bind,source=/media/user/9CE2F7F5E2F7D20E/postgres/,target=/var/lib/postgresql/data --name hashpostgres -e POSTGRES_PASSWORD=postgres99 postgres
#sudo docker run --name phppgadmin -ti -d -p 8080:80 keepitcool/phppgadmin
#docker pull dockage/phppgadmin:latest
## READ https://blog.cloud66.com/how-to-create-the-smallest-possible-docker-image-for-your-golang-application/
