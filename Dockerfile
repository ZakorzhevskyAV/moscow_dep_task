FROM golang:latest
WORKDIR /moscow_dep_task
COPY . .
RUN go build -o moscow_dep_task
CMD ["./moscow_dep_task"]