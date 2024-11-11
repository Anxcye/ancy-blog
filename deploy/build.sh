#!/bin/bash

# Create required directories
mkdir -p ./jars
mkdir -p ./dists/blog-frontend
mkdir -p ./dists/admin-frontend

# Build backend
mvn clean package -DskipTests -f ../ancy-backend

# Move JAR files
mv ../ancy-backend/ancy-admin/target/ancy-admin-1.0-SNAPSHOT.jar ./jars/admin.jar
mv ../ancy-backend/ancy-blog/target/ancy-blog-1.0-SNAPSHOT.jar ./jars/blog.jar

# Build frontend
cd ../ancy-frontend/blog
pnpm install
pnpm build

cd ../admin
pnpm install
pnpm build

# Copy frontend builds
cd ../../deploy
cp -r ../ancy-frontend/blog/dist/* ./dists/blog-frontend
cp -r ../ancy-frontend/admin/dist/* ./dists/admin-frontend