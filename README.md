# **Datenkarte**

Datenkarte is a self-hosted software solution for transforming CSV data into structured JSON payloads and integrating with APIs. It offers dynamic data mapping, validation, and enrichment features. Designed for scalability, Datenkarte runs in Docker or Kubernetes and supports HTTP integrations out of the box.

---

## **Features**

- **Dynamic Data Mapping**:
  Convert CSV headers to JSON keys, including support for nested and dynamic fields.

- **Validation**:
  Ensure data integrity by validating fields with rules like `string`, `number`, `email`, or custom regex patterns.

- **Dry Run Mode**:
  Preview the payloads that would be sent to APIs using the `?dry=true` query parameter.

- **Data Enrichment**:
  Populate missing fields using static values, row-specific values, or predefined arrays.

- **Insert Into Keys**:
  Append values from one field into an existing array field in the payload.

- **HTTP Integration**:
  Send processed JSON payloads to external APIs with customizable headers and authentication.

- **Scalable and Containerized**:
  Runs in Kubernetes or Docker for easy deployment and high availability.

---

## **How It Works**

1. **Define Rules**:
   Use a YAML configuration file to define mappings, validations, and HTTP integration.

2. **Upload CSV Files**:
   Send your CSV files via an HTTP POST endpoint.

3. **Dry Run or Process**:
   Use the `?dry=true` parameter to simulate processing and preview payloads, or run normally to process data and send it to APIs.

4. **Integrate**:
   Processed payloads are sent to APIs, complete with custom headers and authentication.

---

## **Dry Run Mode**

When using the `?dry=true` query parameter, Datenkarte processes the CSV and generates the payloads without sending them to the configured API. This feature is useful for debugging and validating your configuration.

### Example

#### Request:

```bash
curl -X POST -F "file=@data.csv" "http://localhost:8080/dk/upload/string-for-url?dry=true"
```

#### Response:

```json
{
  "status": "dry-run",
  "payloads": [
    {
      "lastName": "Doe",
      "firstName": "John",
      "id": "changeme-1",
      "groups": ["user", "company default group", "New York"]
    },
    {
      "lastName": "Smith",
      "firstName": "Jane",
      "id": "changeme-2",
      "groups": ["user", "company default group", "California"]
    }
  ]
}
```

---

## **Example Configuration**

### YAML Configuration

see config.example.yaml

---

## **Handlers**

**NOTE**: Handlers are currently marked as **TODO** and are not functional in the current version. Future releases will enable handlers to process rows with external binaries for advanced logic.

---

## **Example Input and Output**

### Input CSV

```csv
Last name,First name,Location
Doe,John,New York
Smith,Jane,California
```

### Processed Payloads

#### Row 1

```json
{
  "lastName": "Doe",
  "firstName": "John",
  "id": "changeme-1",
  "groups": ["user", "company default group", "New York"]
}
```

#### Row 2

```json
{
  "lastName": "Smith",
  "firstName": "Jane",
  "id": "changeme-2",
  "groups": ["user", "company default group", "California"]
}
```

---

## **Endpoints**

### Upload CSV

```http
POST /dk/upload/{ruleID}?dry=true
```

- **Description**: Uploads a CSV file for processing based on a specific rule. Use `?dry=true` to preview payloads without sending them to the API.
- **Request**:
  - Content-Type: `multipart/form-data`
  - Body: `file=@data.csv`
- **Example**:

  ```bash
  curl -X POST -F "file=@data.csv" http://localhost:8080/dk/upload/string-for-url
  ```

- **Response (Standard Run)**:

  ```json
  {
    "status": "success",
    "processed_rows": 2
  }
  ```

- **Response (Dry Run)**:

  ```json
[
  {
    "lastName": "Doe",
    "firstName": "John",
    "id": "changeme-1",
    "groups": ["user", "company default group", "New York"]
  },
  {
    "lastName": "Smith",
    "firstName": "Jane",
    "id": "changeme-2",
    "groups": ["user", "company default group", "California"]
  }
]
  ```

---

## **Deployment**

### Docker

1. Build the Docker image:

   ```bash
   docker build -t datenkarte .
   ```

2. Run the container:

   ```bash
   docker run -p 8080:8080 -e AUTH_TOKEN=mytoken -v $(pwd)/config:/app/config datenkarte
   ```

### Kubernetes

1. Deploy the app using the provided Kubernetes configuration.

2. Expose the service:

   ```bash
   kubectl expose deployment datenkarte --type=LoadBalancer --name=datenkarte
   ```

---

## **Error Handling**

- **Validation Errors**:
  - Skips rows with validation issues and logs errors.

- **Dry Run**:
  - Logs payloads that would be sent, ensuring no external API calls.

- **HTTP Errors**:
  - Logs API request failures, including status codes and error messages.

---

## **Key Features in Detail**

### Dry Run

- Use `?dry=true` to preview the payloads without sending them to the API.
- Ideal for debugging configurations.

### Data Validation

- Validate fields like `Last name` and `First name` as strings before processing.

### Dynamic Field Mapping

- Map CSV headers to JSON keys dynamically.
- Use `insert_into` to append values into existing arrays.
