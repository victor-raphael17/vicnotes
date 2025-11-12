# Build stage
FROM node:20-slim as build-stage

WORKDIR /app

# Copy package files first for better layer caching
COPY package*.json ./

# Install all dependencies including devDependencies
RUN npm install --include=dev

# Copy the rest of the application
COPY . .

# Build the app with production flag in the build command
RUN npm run build

## Production stage
FROM debian:trixie-slim as production-stage

RUN apt-get update && apt-get install -y nginx && rm -rf /var/lib/apt/lists/*

# Copy built files from build-stage
COPY --from=build-stage /app/dist /usr/share/nginx/html

# Set proper permissions
RUN chown -R www-data:www-data /usr/share/nginx/html && \
    chmod -R 755 /usr/share/nginx/html

# Expose port 80
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
