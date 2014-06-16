FROM alco/ubuntu-elixir

WORKDIR /src/app
ADD . /src/app/

CMD ["main.ex"]
ENTRYPOINT ["elixirc"]
