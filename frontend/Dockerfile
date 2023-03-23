FROM node:18.10.0 as builder
WORKDIR /app
COPY . .

ARG REACT_APP_API_URL
ARG REACT_APP_WS_URL

RUN npm i
RUN REACT_APP_API_URL=$REACT_APP_API_URL REACT_APP_WS_URL=$REACT_APP_WS_URL npm run build

FROM nginx:stable
WORKDIR /app
COPY ./docker/default.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/build/ /app/
EXPOSE 80 443