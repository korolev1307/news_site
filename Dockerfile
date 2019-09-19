FROM golang

WORKDIR /go/src/github.com/korolev1307/news_site
COPY . .

RUN go get -d -v ./...
RUN go build

CMD ["news_site"]
