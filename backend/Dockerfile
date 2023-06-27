FROM node:18.16-slim

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y openssl

USER node

WORKDIR /home/node/app

COPY package*.json ./
RUN npm config set prefix '/home/node/.npm-global'
RUN npm ci --silent

CMD [ "tail", "-f", "/dev/null" ]
