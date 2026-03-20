# Runtime IP Data

Place offline IP database files here during local or production deployment.

- Expected default paths:
  - `backend/runtime-data/ip/ip2region_v4.xdb`
  - `backend/runtime-data/ip/ip2region_v6.xdb`
- Configure the backend with:
  - `IP2REGION_V4_XDB_PATH`
  - `IP2REGION_V6_XDB_PATH`
- Use `deploy/sync-ip2region.sh` to download or refresh the files
- Do not commit `.xdb` files into git
