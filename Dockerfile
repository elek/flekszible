FROM ubuntu
RUN apt-get update && apt-get install -y git
ADD flekszible /usr/bin/flekszible

