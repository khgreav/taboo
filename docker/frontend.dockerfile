FROM node:lts AS builder

ENV CI=true

WORKDIR /app

RUN npm install --global npm pnpm

COPY frontend .
RUN pnpm install
RUN pnpm run build

# Run

FROM nginx:stable-alpine

RUN rm -rf /usr/share/nginx/html/*
COPY --from=builder /app/dist /usr/share/nginx/html

# TODO: copy nginx config

EXPOSE 80 443

CMD ["nginx", "-g", "daemon off;"]
