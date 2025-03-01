# LETISGO

LETISGO is a robust framework designed to implement the Command Query Responsibility Segregation (CQRS)
pattern with built-in support for both REST.

This framework aims to simplify the development of scalable and maintainable applications
by separating the read and write operations, thus optimizing performance and scalability.

## Development Setup

To set up the development environment for LETISGO.

### Prerequisites:

TailwindCSS CLI requires watchmand installed, and this may vary depending on your operating system.
check the link below for more information on how to install watchman on your operating system.
https://facebook.github.io/watchman/docs/install

### Steps:

1. **Clone the repository:**
   ```sh
   git clone https://github.com/yourusername/letisgo.git
   cd letisgo
   ```
2. **Configure Environment**
   ```sh
   cp .env.example .env
   ```

3. **Define AWS Profile for local developmnent**
   ```sh
   export AWS_PROFILE=localstack
   echo " \n\
[localstack] \n\
aws_access_key_id = localstack \n\
aws_secret_access_key = localstack \n\
region = us-east-1 \n\
endpoint_url = http://localhost:4566 \n\
output = json \n\
    " >> ~/.aws/credentials
   ```

4. **Install dependencies:**
   ```sh
   make prepare
   ```

5. **Set up environment variables:**
   Create a `.env` file in the root directory and add the necessary environment variables.

6. **Run the development server:**
   ```sh
   make dev
   ```

## Usage

LETISGO provides a seamless way to handle both REST and GraphQL APIs. Here’s how you can use it:

### REST API

To define a REST endpoint on `/backend/routes`,
create a controller in the `controllers` directory
and define your routes in the `routes` file.

### Frontend

This framework includes a frontend application built with HTMX, Templ, and Tailwind CSS.
To start the frontend server, run the following command:

```sh
make dev
```

### TODO

- [x] Rest API
    - [x] Routes
    - [x] Views
    - [x] Static Files
- [ ] Frontend:
    - [x] HTMX
    - [x] Tailwind CSS
    - [x] Templ
    - [x] Hot Reloading
- [ ] Add Authentication/Authorization
- [ ] Add Cache support
- [ ] Add DynamoDB support
- [ ] Add EventBridge support
- [ ] Add S3 support
- [ ] Implement a testing framework
- [ ] Add support for WebSockets
