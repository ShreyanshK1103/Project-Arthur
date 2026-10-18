# Development image: Vite dev server with hot reload.
# (glibc-based image on purpose: Vite's native bindings are happiest there.)
FROM node:24-slim
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
EXPOSE 5173
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
