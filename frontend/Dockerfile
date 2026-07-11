FROM node:24-alpine AS build
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm install
COPY . .
RUN npm run build

FROM node:24-alpine
WORKDIR /app
COPY --from=build /app .
EXPOSE 4173
CMD ["npm", "run", "preview"]
