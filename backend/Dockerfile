FROM golang:latest
RUN mkdir -p /usr/src/rga-demo
RUN mkdir -p /usr/src/rga-demo/backend
WORKDIR /usr/bin/rga-demo/backend
COPY . .
EXPOSE 3001
CMD ["go", "run", "."]
