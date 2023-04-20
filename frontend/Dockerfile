FROM node:19-alpine AS build

WORKDIR /workspace

COPY . .

RUN rm -rf ./dist

RUN npm install

RUN npm run build

FROM nginx:alpine

COPY --from=build /workspace/env.sh /docker-entrypoint.d/env.sh
COPY --from=build /workspace/env-template.js env-template.js

COPY --from=build /workspace/dist /usr/share/nginx/html