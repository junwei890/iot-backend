# iot-backend for EE3180
This is a server I'm writing for a solar powered IOT sensor network.

## Dev logs
### 30/8/2025
Done:
- Set up repo.
- Created a database client that endpoints will share.
- Created 3 handlers: `GET /health` to check if the server is up, `POST /data` to insert readings from the STM into InfluxDB and `GET /data` to retrieve data from InfluxDB on an interval.
- Created a server using the **net/http** package.
- Tested endpoints locally using **Postman** to ensure mock data could be inserted and retrieved.

Learning outcomes:
- InfluxDB's Go client documentation was easy to read.

### 1/9/2025
Done:
- Created a **Dockerfile** to test containerization locally.
- Integrated continuous deployment to build and push the Docker image to **Google Artifact Registry**, then deploy on **Google Cloud Run** on merge to main branch.
- Tested endpoints again using **Postman**, now on the URL of the deployed container to ensure data could be inserted and retrieved.

Learning outcomes:
- `/healthz` is actually a reserved route on **GCP**, had to change it to `/health` endpoint.
- You actually need to unignore the binary in `.gcloudignore` so that the Docker image has access to it.
