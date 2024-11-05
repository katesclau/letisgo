# LETISGO

LETISGO is a robust framework designed to implement the Command Query Responsibility Segregation (CQRS) pattern with built-in support for both REST and GraphQL APIs. This framework aims to simplify the development of scalable and maintainable applications by separating the read and write operations, thus optimizing performance and scalability.

## Development Setup

To set up the development environment for LETISGO, follow these steps:

1. **Clone the repository:**
   ```sh
   git clone https://github.com/yourusername/letisgo.git
   cd letisgo
   ```

2. **Install dependencies:**
   ```sh
   make install
   ```

3. **Set up environment variables:**
   Create a `.env` file in the root directory and add the necessary environment variables.

4. **Run the development server:**
   ```sh
   make dev
   ```

## Usage

LETISGO provides a seamless way to handle both REST and GraphQL APIs. Here’s how you can use it:

### REST API

To define a REST endpoint, create a controller in the `controllers` directory and define your routes in the `routes` file.

### GraphQL API

To define a GraphQL endpoint, create a resolver in the `resolvers` directory and define your schema in the `schema` file.

## Extension

LETISGO is designed to be easily extendable. You can add new features or modify existing ones by following these steps:

1. **Add a new command or query:**
   Create a new file in the `commands` or `queries` directory and implement your logic.

2. **Extend the API:**
   Add new routes or resolvers as needed in the `controllers` or `resolvers` directory.

3. **Update the Makefile:**
   If you add new scripts or commands, update the `Makefile` to include these changes.

## Makefile

The `Makefile` is used to automate common tasks. Here are some of the key commands:

- `make install`: Install all dependencies.
- `make dev`: Run the development server.
- `make test`: Run the test suite.
- `make build`: Build the project for production.

By following these guidelines, you can effectively use and extend LETISGO to suit your project's needs.
