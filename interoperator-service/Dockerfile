# Use the official Node.js image as the base image
#FROM node:18-alpine

# Set the working directory inside the container
#WORKDIR /app

# Copy package.json and package-lock.json to the working directory
#COPY package.json package-lock.json ./

# Install dependencies
#RUN npm install

# Copy the rest of the application code to the working directory
#COPY . .

# Set the environment variable to point to the .env file
#ENV NODE_ENV=production
#ENV ENV_FILE_PATH=./src/config/dev.env

# Build the application
#RUN npm run build

# Expose the port the app runs on
#EXPOSE 3000

# Define the command to run the application
#CMD ["npm", "run", "start:prod"]

# Use the official Node.js image
FROM node:18

# Set the working directory inside the container
WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the entire project into the container
COPY . .

# Set environment variables for testing
ENV NODE_ENV=test
ENV API_BASE_URL=http://localhost:3000
ENV AUTH_SERVICE_URL=http://localhost:4000
ENV DOCUMENT_SERVICE_URL=http://localhost:5000
ENV RABBITMQ_URL=amqp://localhost:5672

# Expose the port for debugging (optional)
EXPOSE 9229

# Run all tests
CMD ["npm", "run", "test"]
