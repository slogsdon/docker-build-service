FROM slogsdon/racket

WORKDIR /src/app
ADD . /src/app
RUN raco exe main.rkt

CMD []
ENTRYPOINT ["/src/app/main"]