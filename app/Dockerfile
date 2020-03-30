FROM golang:1.13-alpine3.11

ENV APP_NAME argocd-playground
ENV PORT 8888

COPY . /go/src/${APP_NAME}
WORKDIR /go/src/${APP_NAME}

RUN go get ./
RUN go build -o ${APP_NAME}

CMD ./${APP_NAME}

EXPOSE ${PORT}