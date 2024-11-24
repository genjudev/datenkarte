# **Datenkarte**

Datenkarte is a self-hosted software solution for transforming CSV data into structured JSON payloads and integrating with APIs. It offers dynamic data mapping, validation, enrichment, and now supports custom handlers for advanced data processing. Designed for scalability, Datenkarte runs in Docker or Kubernetes and supports HTTP integrations out of the box.

---

## **Features**

- **Dynamic Data Mapping**:
  Convert CSV headers to JSON keys, including support for nested and dynamic fields.

- **Validation**:
  Ensure data integrity by validating fields with rules like `string`, `number`, `email`, or custom regex patterns.

- **Data Enrichment**:
  Populate missing fields using static values, row-specific values, or predefined arrays.

- **Custom Handlers**:
  Process fields with external scripts or binaries for advanced logic and transformations.

- **Insert Into Keys**:
  Append values from one field into an existing array field in the payload.

- **Dry Run Mode**:
  Preview the payloads that would be sent to APIs using the `?dry=true` query parameter.

- **HTTP Integration**:
  Send processed JSON payloads to external APIs with customizable headers and authentication.

- **Scalable and Containerized**:
  Runs in Kubernetes or Docker for easy deployment and high availability.

---

## **How It Works**

1. **Define Rules**:
   Use a YAML configuration file to define mappings, validations, handlers, and HTTP integration.

2. **Place Handlers**:
   Add your custom handler scripts into the `handlers` directory.

3. **Upload CSV Files**:
   Send your CSV files via an HTTP POST endpoint.

4. **Dry Run or Process**:
   Use the `?dry=true` parameter to simulate processing and preview payloads, or run normally to process data and send it to APIs.

5. **Integrate**:
   Processed payloads are sent to APIs, complete with custom headers and authentication.

---

## **Handlers**

Handlers allow you to process fields with custom logic using external scripts or binaries. This feature is ideal for tasks like hashing passwords, data normalization, or any advanced processing that goes beyond standard mapping and validation.

### **How to Use Handlers**

1. **Create a Handler Script**:
   Write your handler script and place it in the `handlers` directory. Ensure it is executable.

2. **Configure the Handler**:
   Update your YAML configuration to include the handler under the `Handlers` section.

3. **Assign Handler to Fields**:
   In your mapping rules, assign the handler to the fields you want to process.

### **Example Handler Configuration**

```yaml
Handlers:
  - name: "hasher"
    persistent: true
```

- **name**: The name of your handler script (should match the script filename in the `handlers` directory).
- **persistent**: If set to `true`, the handler will remain running between calls, improving performance for resource-intensive scripts.

### **Example Mapping with Handler**

```yaml
each_line:
  - map:
      - name: "Password"
        handlers:
          - "hasher"
```

---

## **Example Configuration**

### **YAML Configuration**

```yaml
Rules:
  - id: "string-for-url"
    delimiter: ";"
    type: http
    http:
      url: "http://localhost:3000/api/test"
      method: POST
      headers:
        - name: "content-type"
          value: ""
        - name: "mykey"
          value: "myvalue"
      auth:
        type: bearer
        value: CHANGEME
      payload_key: "data"
    each_line:
      - map:
          - name: "Last name"
            to: "lastName"
            required: true
          - name: "First name"
            to: "firstName"
          - name: "id"
            fill:
              type: "string"
              prefix: "changeme-"
              value: "row_number"
          - name: "groups"
            fill:
              type: "array"
              value:
                - "user"
                - "company default group"
          - name: "Location"
            insert_into: "groups"
          - name: "Password"
            handlers:
              - "hasher"
        validation:
          - field: "Last name"
            type: "string"
          - field: "First name"
            type: "string"
Handlers:
  - name: "hasher"
    persistent: true
```

---

## **Example Input and Output**

### **Input CSV**

```csv
Last name,First name,Location,Password
Doe,John,New York,Pass123
Smith,Jane,California,Secret456
```

### **Processed Payloads**

#### **Row 1**

```json
{
  "lastName": "Doe",
  "firstName": "John",
  "id": "changeme-1",
  "groups": ["user", "company default group", "New York"],
  "Password": "hashed_value_of_Pass123"
}
```

#### **Row 2**

```json
{
  "lastName": "Smith",
  "firstName": "Jane",
  "id": "changeme-2",
  "groups": ["user", "company default group", "California"],
  "Password": "hashed_value_of_Secret456"
}
```

*Note: The `Password` field is processed by the `hasher` handler.*

---

## **Dry Run Mode**

When using the `?dry=true` query parameter, Datenkarte processes the CSV and generates the payloads without sending them to the configured API. This feature is useful for debugging and validating your configuration.

### **Example**

#### **Request:**

```bash
curl -X POST -F "file=@example.csv" -H "Authorization: Bearer token" "http://localhost:8080/dk/upload/string-for-url?dry=true"
```

#### **Response:**

```json
{
  "status": "dry-run",
  "payloads": [
    {
      "lastName": "Doe",
      "firstName": "John",
      "id": "changeme-1",
      "groups": ["user", "company default group", "New York"],
      "Password": "hashed_value_of_Pass123"
    },
    {
      "lastName": "Smith",
      "firstName": "Jane",
      "id": "changeme-2",
      "groups": ["user", "company default group", "California"],
      "Password": "hashed_value_of_Secret456"
    }
  ]
}
```

---

## **Endpoints**

### **Upload CSV**

```http
POST /dk/upload/{ruleID}?dry=true
```

- **Description**: Uploads a CSV file for processing based on a specific rule. Use `?dry=true` to preview payloads without sending them to the API.
- **Request**:
  - Content-Type: `multipart/form-data`
  - Body: `file=@data.csv`
- **Example**:

  ```bash
  curl -X POST -F "file=@example.csv" -H "Authorization: Bearer token" "http://localhost:8080/dk/upload/string-for-url"
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
  {
    "status": "dry-run",
    "payloads": [
      {
        "lastName": "Doe",
        "firstName": "John",
        "id": "changeme-1",
        "groups": ["user", "company default group", "New York"],
        "Password": "hashed_value_of_Pass123"
      },
      {
        "lastName": "Smith",
        "firstName": "Jane",
        "id": "changeme-2",
        "groups": ["user", "company default group", "California"],
        "Password": "hashed_value_of_Secret456"
      }
    ]
  }
  ```

---

## **Deployment**

### **Docker**

1. **Build the Docker image:**

   ```bash
   docker build -t datenkarte .
   ```

2. **Run the container:**

   ```bash
   docker run -p 8080:8080 -e AUTH_TOKEN=mytoken -v $(pwd)/config:/app/config -v $(pwd)/handlers:/app/handlers datenkarte
   ```

   - **Note**: Mount the `handlers` directory to include your custom handler scripts.

### **Kubernetes**

1. **Deploy the app using the provided Kubernetes configuration.**

2. **Expose the service:**

   ```bash
   kubectl expose deployment datenkarte --type=LoadBalancer --name=datenkarte
   ```

3. **Add Handlers:**

   - Use a ConfigMap or PersistentVolume to include your handler scripts in the deployment.

---

## **Error Handling**

- **Validation Errors**:
  - Skips rows with validation issues and logs errors.

- **Handler Errors**:
  - Logs any errors encountered during handler execution.

- **Dry Run**:
  - Logs payloads that would be sent, ensuring no external API calls.

- **HTTP Errors**:
  - Logs API request failures, including status codes and error messages.

---

## **Key Features in Detail**

### **Custom Handlers**

- **Purpose**: Extend Datenkarte's functionality by processing fields with custom scripts.
- **Usage**:
  - Place your handler scripts in the `handlers` directory.
  - Make sure the scripts are executable.
  - Configure the handlers in your YAML file.
- **Persistent Handlers**:
  - Set `persistent: true` if your handler should remain running between calls.
  - Useful for handlers that have a startup cost or need to maintain state.

### **Dry Run**

- Use `?dry=true` to preview the payloads without sending them to the API.
- Ideal for debugging configurations.

### **Data Validation**

- Validate fields like `Last name` and `First name` as strings before processing.
- Skipped rows are logged for review.

### **Dynamic Field Mapping**

- Map CSV headers to JSON keys dynamically.
- Use `insert_into` to append values into existing arrays.
- Supports nested JSON structures.

---

## **Contact and Support**

For issues, feature requests, or contributions, please open an issue on our [GitHub repository](https://github.com/your-repo/datenkarte).

---

## **License**

Datenkarte is released under the [MIT License](LICENSE).
