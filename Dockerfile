# Use the official Ubuntu base image from Docker Hub
FROM ubuntu:latest

# Set environment variables, if needed
ENV DEBIAN_FRONTEND=noninteractive

# Update the package repository and install necessary packages
RUN apt-get update && \
    apt-get install -y \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /VelociStore

# Copy your application files into the container
COPY . /VelociStore

# Specify the command to run on container startup
CMD ["bash"]
