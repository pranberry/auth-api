# auth-api

Creates and returns a JWT. Allows a user to access protected endpoints with a valid JWT

## Endpoints

- `/register` register a user
- `/login` login and retrieve a JWT
- `/secret` validates a legit JWT and sends the client some guarded assets

## Build

- This project is containerized. Build with: `docker-compose up --build`
- Run the API with the vscode debugger, use the `.env.local` as `envFile` in `launch.json`
- (The `.env` files and `launch.json` live out of this repo)