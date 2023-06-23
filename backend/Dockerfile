FROM node:18.16-slim

RUN apt-get update && \
    apt-get upgrade -y

WORKDIR /home/node/app

COPY *.json .

RUN npm ci

USER node

CMD [ "tail", "-f", "/dev/null" ]
