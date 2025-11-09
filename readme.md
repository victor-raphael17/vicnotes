# VicNotes Remake

## The ideia of this project is to make a remake of VicNotes, a note taking app that I made on university.

## The tech stack will vary based on the version of the project. There will be a version for running it locally, a cloud version to run with low traffic and another cloud version to run with high traffic.

## Tech stack:

### AI

- Windsurf
- Quick AI programming technique from Felipe Forbeck

### Code

- Vue.js for the frontend
- Tailwind CSS for styling
- Go for the backend

### Infra

- Nginx for reverse proxy
- PostgreSQL for registering the users
- Docker for containerization 
- GitHub for version control
- AWS S3 for storing the notes context (cloud versions)
- AWS EC2 for the server (cloud versions)
- Redis for caching (cloud versions)
- GitHub Actions to CI/CD (cloud versions)
- Cloudflare for WAF (cloud versions)
- Porkbun for domain ownership (cloud versions)

## Configuration for devs
1. Make sure you have the following dependencies installed:
    - Docker and Docker Compose
2. Run the config.sh
