# Use a single-stage build with Node
FROM node:22-alpine

# This Dockerfile is designed to work with the docker-compose.yml.
# The 'args' in docker-compose.yml will pass variables from your root .env file
# to these ARG instructions at build time.

# --- Environment Variables Section ---
# Declare build arguments. Add one for each variable from your .env file that the frontend needs.
# Vite requires variables to be prefixed with VITE_.
ARG VITE_API_URL

# Set environment variables from the build arguments.
# Vite will read these ENV variables during the 'yarn build' step.
ENV VITE_API_URL=$VITE_API_URL
# --- End of Section ---

# Set the working directory
WORKDIR /app

# Copy package.json and yarn.lock
COPY package.json yarn.lock ./

# Install all dependencies needed for the build
RUN yarn install --frozen-lockfile

# Copy the rest of the application source code
COPY . .

# Run the build script. Vite will embed the environment variables
# from the 'ENV' instructions above into the static files.
RUN yarn build

# Expose the port the app will run on.
# This should match the port used in the CMD instruction.
EXPOSE 3000

# Start the app using Vite's preview server.
# This serves the files from the 'dist' directory.
# '--host 0.0.0.0' is crucial to make the server accessible from outside the container.
CMD ["yarn", "serve", "--host", "0.0.0.0", "--port", "3000"]
