# Use an official Node.js runtime as a parent image
FROM node:18.17.1 as build

# Set the working directory
WORKDIR /usr/local/app

# Add the source code to app
COPY ./ /usr/local/app/

# Install all the dependencies
RUN npm install

# Generate the build of the application
RUN npm run build


# Stage 2: Serve app with nginx server

# Use official nginx image as the base image
FROM nginx:latest

# Copy the build output to replace the default nginx contents.
RUN rm -rf /usr/share/nginx/html/*
COPY nginx.conf /usr/share/nginx/html
COPY --from=build /usr/local/app/dist/webapp/browser /usr/share/nginx/html
RUN chown nginx:nginx /usr/share/nginx/html/*

EXPOSE 80
