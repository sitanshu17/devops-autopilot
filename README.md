# DevOps Autopilot

An intelligent DevOps automation tool that generates Terraform infrastructure code using OpenAI's GPT API. Built with Go for high performance and reliability.

## 🚀 Features

- **AI-Powered Code Generation**: Uses OpenAI GPT to generate Terraform code based on natural language specifications
- **RESTful API**: Clean HTTP endpoints for integration with other tools
- **File Management**: Automatically saves generated code with unique filenames
- **Robust Error Handling**: Comprehensive error handling and validation
- **High Performance**: Built with Go for speed and efficiency

## 📋 Prerequisites

1. **Go 1.18 or higher** - [Download Go](https://golang.org/dl/)
2. **OpenAI API Key** - Get one from [OpenAI Platform](https://platform.openai.com/api-keys)

## 🛠️ Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/devops-autopilot.git
   cd devops-autopilot
   ```

2. **Create a `.env` file with your OpenAI API key:**
   ```env
   OPENAI_API_KEY=your_openai_api_key_here
   PORT=5000
   ```

3. **Install Go dependencies:**
   ```bash
   go mod tidy
   ```

## 🏃‍♂️ Running the Application

### Option 1: Run directly
```bash
go run main.go
```

### Option 2: Build and run
```bash
go build -o devops-autopilot
./devops-autopilot
```

The server will start on `http://localhost:5000`

## 📡 API Endpoints

### Health Check
```http
GET http://localhost:5000/api/provision/health
```

**Response:**
```json
{
  "status": "Service is healthy"
}
```

### Generate Terraform Code
```http
POST http://localhost:5000/api/provision/terraform
Content-Type: application/json

{
  "resource": "EC2 instance",
  "specs": "t3.micro instance in us-east-1 with Ubuntu 20.04"
}
```

**Response:**
```json
{
  "message": "Terraform code generated successfully",
  "terraformCode": "resource \"aws_instance\" \"example\" {\n  ami = \"ami-0c55b159cbfafe1d0\"\n  instance_type = \"t3.micro\"\n}"
}
```

## 📁 Project Structure

```
devops-autopilot/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── go.sum                  # Go dependencies checksum
├── handlers/
│   └── provision.go        # HTTP handlers (REST controllers)
├── utils/
│   └── openai.go          # OpenAI API integration
├── terraform/             # Generated Terraform files
├── .env                   # Environment variables (not committed)
├── .gitignore             # Git ignore rules
└── README.md              # This file
```

## 🔧 Configuration

Create a `.env` file in the project root:

```env
# Required
OPENAI_API_KEY=sk-your-openai-api-key-here

# Optional (defaults shown)
PORT=5000
```

## 🚀 Building for Production

```bash
# Build for current platform
go build -o devops-autopilot

# Cross-platform builds
GOOS=windows GOARCH=amd64 go build -o devops-autopilot.exe
GOOS=linux GOARCH=amd64 go build -o devops-autopilot
GOOS=darwin GOARCH=amd64 go build -o devops-autopilot
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⚠️ Security Note

Never commit your `.env` file or expose your OpenAI API key. The `.gitignore` file is configured to exclude sensitive files.
