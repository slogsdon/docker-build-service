FROM dockerfile/python

WORKDIR /src/app
ADD . /src/app/

CMD ["main.py"]
ENTRYPOINT ["python"]
