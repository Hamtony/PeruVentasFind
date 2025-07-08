# frontend/Dockerfile

# Etapa de construcción
FROM node:20 AS builder
WORKDIR /app
COPY . .

# Asegúrate que node_modules incluye tailwind y postcss
RUN npm install
RUN npm run build

# Etapa de producción
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
