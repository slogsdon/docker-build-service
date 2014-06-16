FROM ruby

WORKDIR /src/app
ADD . /src/app/

CMD ["main.rb"]
ENTRYPOINT ["ruby"]
