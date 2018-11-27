FROM golang as builder
RUN apt-get update && apt-get install -y libpcap-dev
WORKDIR /interlockmon
COPY ./go.mod .
RUN go mod download
COPY . .
RUN go install 

FROM ubuntu
RUN apt-get update && apt-get install -y libpcap-dev
WORKDIR /interlockmon
COPY --from=builder /go/bin/interlockmon .
ENTRYPOINT [ "./interlockmon" ]