FROM golang
LABEL maintainer="jessezhang007007 <jessezhang007007@gmail.com>"

RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - \
    && echo 'deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main' | tee /etc/apt/sources.list.d/google-chrome.list \
    && apt-get update \
    && apt-get install -y google-chrome-stable \
    && go get -u github.com/chromedp/chromedp \
    && go get -u github.com/chromedp/cdproto/emulation \
    && go get -u github.com/chromedp/cdproto/page \
    && apt-get install -y locales ttf-wqy-microhei ttf-wqy-zenhei xfonts-wqy xfonts-intl-chinese fonts-arphic-uming fonts-noto
WORKDIR /go/src/app

ADD app.go ./
RUN go build app.go
EXPOSE 8082
CMD ["/go/src/app/app"]