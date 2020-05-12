FROM node:12.16.3 as build
ENV REACT_APP_API_URL=https://api.notarybot.bots.pyr.sh
WORKDIR /usr/src/app
COPY package.json yarn.lock ./
RUN yarn
COPY . ./
RUN yarn build

FROM nginx:1.17.10-alpine
COPY --from=build /usr/src/app/build /usr/share/nginx/html
COPY default.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
