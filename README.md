# TestAnchor - Test Case Management Tool

## Overview
TestAnchor is a simple and modern tool for managing test cases. It is designed to be user-friendly and intuitive, enabling quick creation and editing of test cases, as well as efficient test management and tracking.

**Note**: This project is currently under development, and key features are still being implemented.

## Installation

To install TestAnchor, follow these steps:

1. Clone the repository using the following command:
```
git clone git@github.com:rmuraoka/test-anchor.git
cd test-anchor
```
2. Set up the environment variables as described in the [Environment Variables Setup](#environment-variables-setup) section.
3. Ensure Docker is installed on your system.
4. Run the following command to start the application:
```
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

This will build and start the necessary services for TestAnchor using Docker.

## Environment Variables Setup
This project utilizes the following environment variables. Please set these variables in your `.env` file.

- `DB_HOST`: The host address of the database server. In docker-compose, this is the DB service name. Example: `db`
- `DB_USER`: Username for accessing the database. Example: `user`
- `DB_PASSWORD`: Password for the database user. Example: `password`
- `DB_NAME`: Name of the database to use. Example: `testanchorDB`
- `MYSQL_ROOT_PASSWORD`: Password for the MySQL database's root user. Example: `anotherpassword`
- `BACKEND_URL`: URL of the backend server. Example: `http://localhost:8000`
- `FRONTEND_ORIGIN`: Origin URL of the frontend server. Example: `http://localhost:3000`
- `JWT_SECRET_KEY`: Secret key used for JSON Web Token (JWT). Generate and use a random string.
- `INITIAL_USER_EMAIL`: Email address of the initial user. Example: `user@example.com`
- `INITIAL_USER_NAME`: Name of the initial user. Example: `Username`
- `MAIL_HOST`: Host for the mail service. `mailhog` is often used for development.
- `MAIL_PORT`: Port number for the mail service. Example: `1025`
- `MAIL_USERNAME`: Username for the mail service. Example: `user@mail.com`
- `MAIL_PASSWORD`: Password for the mail service. Example: `mailpassword`
- `FROM_EMAIL`: Sender email address for outgoing emails. Example: `noreply@example.com`
- `USE_TLS`: Whether to use TLS for email sending. `true` or `false`.

**Note**: The `.env` file contains sensitive information, so do not upload it to public repositories.

## Usage
Usage instructions will be updated in accordance with the progress of the project.

### API Documentation

Our project's API documentation is provided through Swagger. To test the API locally and view the documentation, please visit the following URL:

- Swagger UI: [http://localhost:8080](http://localhost:8080)

Details about how to use the API and information on the endpoints are available on the Swagger UI.

## How to Contribute
TestAnchor is an open-source project and welcomes contributions from the community. For details on how to contribute, please refer to [CONTRIBUTING.md](/CONTRIBUTING.md).

## License
This project is published under the [MIT License](/LICENSE).

## Version Information
Version information will be updated here as the project progresses.

## Contact
If you need feedback or support, please use the [Issue Tracker](https://github.com/rmuraoka/test-anchor/issues).
