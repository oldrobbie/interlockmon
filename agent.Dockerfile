FROM golang as builder
WORKDIR /interlockmon
COPY ./go.mod .
RUN go mod download
COPY . .
RUN go install 

FROM ubuntu
WORKDIR /interlockmon
COPY --from=builder /go/bin/interlockmon .
ENTRYPOINT [ "./interlockmon" ]