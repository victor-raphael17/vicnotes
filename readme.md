# VicNotes Remake

## The ideia of this project is to make a remake of VicNotes, a note taking app that I made on university.

## Technological stack that I will use are:
- Quick AI programming technique from Felipe Forbeck
- PostgreSQL for registering the users
- AWS S3 for storing the notes context
- Node.js for the backend
- Vue.js for the frontend
- Tailwind CSS for styling
- Redis for caching
- Docker for containerization
- Nginx for reverse proxy
- GitHub for version control
- GitHub Actions to CI/CD
- AWS EC2 for the server
- Cloudflare for WAF
- Porkbun for domain ownership

### The utilization of these technologies will vary based on the version of the project, there will be a version for running locally, a cloud version to run with low traffic and a version to run with high traffic.

## Configuration for devs
1. Make sure you have the following dependencies installed:
    - Docker and Docker Compose
    - Node.js
2. Run the config.sh
3. Make sure your backend package.json has the following scripts:
    "test": "echo \"Error: no test specified\" && exit 1",
    "start": "node src/index.js",
    "dev": "nodemon src/index.js"

