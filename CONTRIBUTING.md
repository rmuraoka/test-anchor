# Contribution Guidelines - TestAnchor

Thank you for your interest in contributing to the TestAnchor project! Our project is open-source and we welcome contributions from the community. Please read the following guidelines to contribute effectively to this wonderful community.

## Before Contributing
- Please review the [README](/README.md) and [LICENSE](/LICENSE) of this project.
- Check for already reported Issues and Pull Requests to avoid duplication.

## Ways to Contribute
1. **Reporting Issues**: If you find bugs or have suggestions, please create an Issue first.
2. **Submitting Pull Requests**: If you want to contribute directly to the code, create a new branch, commit your changes, and then submit a Pull Request.

## Coding Standards
- Strive for clean and readable code.
- Commit messages should be clear and specific.

## Testing
- When adding new features, also add relevant tests.
- Ensure that existing tests are not broken.

### Backend
- The instructions to run the backend tests are as follows:
```
cd backend
go test ./...
```

### E2E Tests
- The instructions to run the E2E tests are as follows:
- To avoid impacting the development environment, the containers are separate.
```
docker-compose --env-file .env.test -f docker-compose.yml -f docker-compose.test.yml up --build -d
npm test
docker-compose --env-file .env.test -f docker-compose.yml -f docker-compose.test.yml down
```
Please create a `.env.test` file. Its basic contents are the same as `.env`, but with the following differences:

- `MAIL_HOST`: Set it to `mailhog`.
- `MAIL_PORT`: Use the mailhog port number (1025).
- `USE_TLS`: Set to `false`.
- `TEST_LOGIN_PASSWORD`: If an initial user has already been created, specify the password. This will skip the password verification process.

## Documentation
- Update the documentation as necessary for new features and changes.

## Code Reviews
- Pull Requests will be reviewed by community members.
- Engage in constructive discussion in response to feedback from reviews.

If you have any questions about contributing, please contact us through the [Issue Tracker](https://github.com/rmuraoka/test-anchor/issues).

We appreciate your contributions to the TestAnchor project!
