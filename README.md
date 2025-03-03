# Restcol

**One Single RESTful API for Collaborating, Sharing, and Streaming Data**

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GitHub Issues](https://img.shields.io/github/issues/FootprintAI/restcol)](https://github.com/FootprintAI/restcol/issues)
[![GitHub Stars](https://img.shields.io/github/stars/FootprintAI/restcol)](https://github.com/FootprintAI/restcol/stargazers)

## Overview

`Restcol` is a RESTful document storage solution designed for collaboration, built to work with any kind of storage backend. It organizes data into **collections** and **documents**, providing a flexible and scalable way to manage and share data. Collections group documents with similar schemas, enabling schema evolution tracking over time, while documents store client data in formats like JSON, CSV, XML, or even media files. No predefined schema is required—schemas are dynamically created or updated with each document request.

This project aims to simplify data collaboration and streaming by offering a unified API that adapts to your application's needs.

## Features

- **Flexible Storage**: Works with any storage backend.
- **Collections**: Organizes documents with similar schemas for easy management and schema change detection.
- **Dynamic Schemas**: Automatically creates or modifies schemas based on document requests—no upfront schema definition needed.
- **Supported Formats**: Handles JSON, CSV, XML, and media data.
- **RESTful API**: Simple, intuitive endpoints for collaboration and data streaming.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/FootprintAI/restcol.git
   cd restcol
   ```

2. Install dependencies (assuming a Go-based project; adjust if different):
   ```bash
   go mod tidy
   ```

3. Build and run:
   ```bash
   go build
   ./restcol
   ```

*Note*: Specific setup instructions may vary depending on your environment and storage backend. Check the source code or configuration files for additional requirements.

## Usage

### Basic Example
To create a collection and add a document via the API:

```bash
# Create a new collection
curl -X POST http://localhost:8080/collections -d '{"name": "my-collection"}'

# Add a document to the collection
curl -X POST http://localhost:8080/collections/my-collection/documents -d '{"data": {"id": 1, "name": "example"}}'
```

For detailed API documentation, refer to the [API Reference](#api-reference) section (to be added).

## Configuration

- **Storage Backend**: Configure your preferred storage system (e.g., local filesystem, S3, etc.) in the config file or environment variables.
- **Port**: Default is `8080`. Override with the `PORT` environment variable.

Example configuration:
```bash
export STORAGE_TYPE="filesystem"
export STORAGE_PATH="/path/to/storage"
export PORT=8080
```

## Contributing

We welcome contributions! To get started:

1. Fork the repository.
2. Create a feature branch: `git checkout -b my-feature`.
3. Commit your changes: `git commit -m "Add my feature"`.
4. Push to your fork: `git push origin my-feature`.
5. Open a pull request.

Please read our [Contributing Guidelines](CONTRIBUTING.md) for more details (to be created if not present).

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## Roadmap

- Add support for additional storage backends.
- Implement real-time streaming capabilities.
- Enhance schema versioning and migration tools.

## Contact

For questions or support, open an issue on the [GitHub Issues page](https://github.com/FootprintAI/restcol/issues) or reach out to the [FootprintAI team](https://github.com/FootprintAI).

---

*Maintained by [FootprintAI](https://github.com/FootprintAI)*  
*Last updated: March 02, 2025*
