FROM debian:stretch-slim
RUN apt-get update && apt-get install -y ca-certificates &&  apt-get clean all
COPY ./app /app
ENTRYPOINT /app