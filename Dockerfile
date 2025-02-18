FROM gobase:1.0.0

COPY . .
RUN go mod tidy
