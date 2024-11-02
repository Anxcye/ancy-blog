mvn clean package -DskipTests -f ../ancy-backend
mv ../ancy-backend/ancy-admin/target/ancy-admin-1.0-SNAPSHOT.jar ./jars/admin.jar
mv ../ancy-backend/ancy-blog/target/ancy-blog-1.0-SNAPSHOT.jar ./jars/blog.jar

cd ../ancy-frontend/blog
pnpm build
cd ../admin
pnpm build
cd ../../deploy
cp -r ../ancy-frontend/blog/dist ./dists/blog-frontend
cp -r ../ancy-frontend/admin/dist ./dists/admin-frontend
