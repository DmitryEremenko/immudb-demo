# Accounting Information Application

This application is a simple accounting information system built with a Go backend using immudb, and a React frontend with Material-UI. It allows users to add and view account information.

## Prerequisites

Before you begin, ensure you have the following installed on your system:
- Docker
- Docker Compose

## Quick Start Guide

If you just want to run the application and play around with it, follow these simple steps:

1. Clone the repository:
   ```
   git clone <repository-url>
   cd <repository-folder>
   ```

2. Start the application using Docker Compose:
   ```
   docker-compose up
   ```

3. Wait for the containers to start up. You should see logs indicating that the services are running.

4. Open your web browser and go to `http://localhost:3000`

5. You can now use the application:
   - Use the form to add new account information
   - View the list of accounts in the table below the form

6. When you're done, stop the application by pressing `Ctrl+C` in the terminal where Docker Compose is running

That's it! You can now explore the application and its features.

## Getting Started (for Development)

If you want to set up the application for development, follow these more detailed steps:

1. Clone the repository:
   ```
   git clone <repository-url>
   cd <repository-folder>
   ```

2. Create a `.env` file in the `frontend` directory with the following content:
   ```
   VITE_API_URL=http://localhost:8080
   ```

3. Build and run the Docker containers:
   ```
   docker-compose up --build
   ```

4. Wait for the containers to start up. You should see logs indicating that the services are running.

5. Access the application:
   - Frontend: Open your browser and go to `http://localhost:3000`
   - Backend API: Available at `http://localhost:8080`

## Project Structure

The project is organized into three main components:

- `backend/`: Contains the Go backend code
- `frontend/`: Contains the React frontend code
- `docker-compose.yml`: Defines the multi-container Docker application

### Backend

The backend is built with Go and uses the following main dependencies:
- Gin: Web framework
- immudb: Immutable database for storing account information

### Frontend

The frontend is built with React and Vite, and uses the following main dependencies:
- Material-UI: For styling and UI components
- Axios: For making HTTP requests to the backend

## Usage

Once the application is running:

1. Open your web browser and navigate to `http://localhost:3000`.
2. Use the form to add new account information:
   - Enter the account number, name, IBAN, address, and amount
   - Select the account type (sending or receiving)
   - Click the "Add Account" button to submit the information
3. View the list of accounts in the table below the form:
   - The table displays all entered accounts with their details
   - The table is automatically updated when new accounts are added

## Development

To make changes to the application:

1. Stop the running containers with `Ctrl+C` or `docker-compose down`.
2. Make your changes to the backend or frontend code:
   - Backend: Modify the Go files in the `backend/` directory
   - Frontend: Update the React components in the `frontend/src/` directory
3. Rebuild and run the containers with `docker-compose up --build`.
4. Test your changes in the browser at `http://localhost:3000`.

