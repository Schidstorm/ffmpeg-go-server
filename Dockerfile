FROM golang as build

WORKDIR /go/src/github.com/schidstorm/ffmpeg-go-server/
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

FROM alpine
COPY --from=mwader/static-ffmpeg /ffmpeg /usr/bin/ffmpeg
COPY --from=mwader/static-ffmpeg /ffprobe /usr/bin/ffprobe
COPY --from=build /go/bin/fmpeg-go-server /usr/bin/fmpeg-go-server

ENTRYPOINT [ "ffmpeg-go-server" ]