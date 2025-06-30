FROM node:20 AS builder
WORKDIR /app

COPY . .

RUN npm install --legacy-peer-deps --ignore-scripts
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
