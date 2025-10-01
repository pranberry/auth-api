# auth-api

Creates and returns a JWT. Allows a user to access protected endpoints with a valid JWT

## Endpoints

- `/register` register a user
- `/login` login and retrieve a JWT
- `/secret` validates a legit JWT and sends the client some guarded assets

## Build

- This project is containerized. Build with: `docker-compose up --build`
- If to run the API with the debugger, use the `.env.local` file in the vscode's launch.json