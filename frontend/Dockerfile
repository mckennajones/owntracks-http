FROM node:19-alpine AS build

WORKDIR /workspace

COPY . .

RUN rm -rf ./dist

RUN npm install

RUN npm run build

FROM nginx:alpine

# Copy the static website
# Use the .dockerignore file to control what ends up inside the image!
COPY --from=build /workspace/dist /usr/share/nginx/html
