FROM golang:1.24

WORKDIR /app 

COPY . .

RUN go build -o meal-api

EXPOSE 8080

CMD ["./meal-api"]