FROM node:16.15-alpine
RUN mkdir -p /usr/src/rga-demo
RUN mkdir -p /usr/src/rga-demo/frontend
WORKDIR /usr/src/rga-demo/frontend
COPY . .
EXPOSE 3000
RUN npm install
CMD ["npm", "run", "dev"]
