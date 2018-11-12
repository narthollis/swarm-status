FROM golang:1.11 as builder

RUN mkdir -p /go/src/swarm-status
COPY . /go/src/swarm-status
WORKDIR /go/src/swarm-status
RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main swarm-status

FROM scratch
COPY --from=builder /go/src/swarm-status/main /
COPY static /static
COPY templates /templates

#ENV METRIC_PAGE_SRC about:blank
#ENV DISPLAY_NAME_LABEL_KEY com.example.display.name
#ENV DISPLAY_GROUP_LABEL_KEY com.example.display.group
#ENV DOCKER_HOST

CMD ["/main"]
