# Stage 1: Build the React application
FROM node:20-alpine AS build

# Set the working directory
WORKDIR /app

# Copy package.json and yarn.lock
COPY package.json yarn.lock ./

# Install dependencies using Yarn
# Using --frozen-lockfile is a best practice for CI/CD to ensure exact dependency versions
RUN yarn install --frozen-lockfile

# Copy the rest of the application source code
COPY . .

# Run the build script using Yarn
# This command is taken from your package.json
RUN yarn build

# Stage 2: Serve the application using Nginx
FROM nginx:stable-alpine

# Copy the build output from the 'build' stage to the Nginx html directory
# Vite's default build output directory is 'dist'
COPY --from=build /app/dist /usr/share/nginx/html

# Copy a custom Nginx configuration file
# This is important for single-page applications (SPAs) like React
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
