
# Voice-Enabled Task Tracker

This project is a fullstack application with a Go backend and React frontend. It allows users to create, update, and delete tasks manually or via voice input, with real-time board updates and a clear, accessible UI.

## Project Structure

```
voice-task-tracker/
├── backend/
│   ├── main.go
│   ├── handlers/
│   │   ├── handlers.go
│   │   └── voice_parser.go
│   ├── middleware/
│   │   └── middleware.go
│   ├── routes/
│   │   └── routes.go
│   ├── go.mod
├── frontend/
│   ├── src/
│   │   ├── App.js
│   │   ├── TaskList.js
│   │   ├── TaskModal.js
│   │   ├── EditTaskModal.js
│   │   ├── VoicePreviewModal.js
│   │   └── ...
│   ├── public/
│   │   └── index.html
│   ├── package.json
├── README.md
```

## Setup Instructions

### Backend
1. Navigate to the `backend` directory.
2. Run `go mod tidy` to install dependencies.
3. Start the server with `go run main.go`.

### Frontend
1. Navigate to the `frontend` directory.
2. Run `npm install` to install dependencies.
3. Start the development server with `npm start`.

## Usage

- Backend API: `http://localhost:8080`
- Frontend: `http://localhost:3000`

## Features

- Create, update, and delete tasks
- Voice input for task creation (with preview and edit)
- Edit tasks via modal
- Real-time board updates
- Accessible UI (ARIA labels, keyboard navigation)

## API Design

- RESTful endpoints: `/tasks` (GET, POST), `/tasks/{id}` (PUT, DELETE), `/parse-voice` (POST)
- Proper HTTP status codes and error responses

## Assumptions & Design Decisions

- In-memory storage for tasks (no database)
- Voice parsing uses simple keyword extraction
- All modals are accessible and keyboard-navigable
- Error messages are shown to users via alerts
- Default values: status = "To Do", priority = "Medium"

## Accessibility

- All interactive elements have ARIA labels
- Modals are focus-trapped and keyboard accessible

## Contributing

Feel free to submit issues or pull requests for improvements or bug fixes.