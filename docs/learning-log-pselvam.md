- We're using a monorepo where we'll host all the code relating to our app. Things can be deployed separately but the code all lives together, which is nice.
- The backend directory is where go code lives. There's 2 builds we have here:
    - The api (at cmd/api): This is the backend API for the bookmarks service
    - The worker (at cmd/worker): This is for running async workers
- The code in internal and migrations is where the guts really live, but the cmd is basically the entry point. This is where we have our "main" packages.
- We use Docker to build and containerize our app.
- We started with a Dockerfile in the backend directory. This uses a multi-stage process to 1. Build the app binary and 2. Run the app binary in a container using Alpine linux. Alpine is very small, so this ensures things build fast.
- The Dockerfile contains a list of instructions (think of it like a recipe). docker-compose.yml is also instructions but this is about how to string together multiple different services/containers and add networking between them all. A docker process will read the instructions in the dockerfile and then create an image/blueprint. This image/blueprint can then be run in a container (you can spin up multiple containers from the same blueprint).
- What's super cool here is someone else could download my code, do a "docker compose up --build" and bam they have the same local development setup as me. Really fucking cool honestly.


- We used fly.io for deploying the backend. This was extremely simple and basically takes your dockerfile and uses it to build up the same code in a container on fly.io
- We used Cloudflare pages to deploy the frontend. This was also super simple. You just do npm run build to first get your dist directory, and then you just deploy that as a static site on cloudflare pages
- The frontend was built using React + Vite. This is a lot simpler than nextjs - I always thought of nextjs as being a wrapper around React to setup and manage React apps but it's basically it's own framework that's an extension of React. So it's a lot more complex than a basic React app.