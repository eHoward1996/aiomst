FROM arm64v8/golang:alpine
RUN apk add --no-cache \
	ffmpeg \
	gcc \
	git \
	g++ \
	musl-dev \
	pkgconf \
	sqlite-dev \
	taglib-dev

RUN go get -u -v -d github.com/eHoward1996/aiomst
RUN go build -o /go/aiomst /go/src/github.com/eHoward1996/aiomst

EXPOSE 8090
ENTRYPOINT /go/aiomst --media=/media/Media/Music --sqlite=/media/Media/mediadb.db
