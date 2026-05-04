FROM node:24-alpine AS base

WORKDIR /pravah

COPY package*.json ./

FROM base AS development

RUN npm install
COPY . .

EXPOSE 3000

# USER node

CMD ["npm", "run", "dev"]

# FROM base AS builder

# RUN npm install

# COPY . .

# RUN npm run build

FROM base AS production

RUN npm install 
COPY . .

CMD ["npm","run","start"]
