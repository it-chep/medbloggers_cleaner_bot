FROM golang:alpine
RUN apk add ffmpeg
WORKDIR /app
COPY . .
RUN go mod download
#RUN apt-get -y install make
EXPOSE 8000
CMD ["go", "run", "cmd/medbloggers_cleaner_bot/local/main.go"]
