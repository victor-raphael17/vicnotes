# VicNotes Frontend

A modern, responsive Vue.js frontend for the VicNotes note-taking application.

## Features

- **User Authentication**: Register and login with JWT tokens
- **Note Management**: Create, read, update, and delete notes
- **Responsive Design**: Beautiful UI with Tailwind CSS
- **Real-time Updates**: Instant note synchronization
- **Dark-friendly**: Modern gradient design with smooth animations

## Tech Stack

- **Vue.js 3**: Progressive JavaScript framework
- **Vite**: Next-generation frontend build tool
- **Tailwind CSS**: Utility-first CSS framework
- **Pinia**: State management
- **Axios**: HTTP client
- **Lucide Vue**: Beautiful icon library

## Project Structure

```
frontend/
├── src/
│   ├── pages/           # Page components (Login, Register, Notes)
│   ├── stores/          # Pinia stores (auth, notes)
│   ├── api/             # API client configuration
│   ├── router/          # Vue Router configuration
│   ├── App.vue          # Root component
│   ├── main.js          # Application entry point
│   └── style.css        # Global styles
├── index.html           # HTML template
├── vite.config.js       # Vite configuration
├── tailwind.config.js   # Tailwind configuration
└── package.json         # Dependencies
```

## Setup

### Prerequisites

- Node.js 16+ and npm/yarn/pnpm

### Installation

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

The development server runs on `http://localhost:5173` and proxies API requests to `http://localhost:8080`.

## Environment Variables

Create a `.env.local` file:

```
VITE_API_URL=http://localhost:8080
```

## API Integration

The frontend communicates with the Go backend via REST API:

- **Authentication**: `/api/v1/auth/register`, `/api/v1/auth/login`
- **Notes**: `/api/v1/notes` (CRUD operations)

Authentication tokens are stored in localStorage and automatically included in requests.

## Features

### Authentication
- User registration with email and password
- Secure login with JWT tokens
- Automatic logout on token expiration
- Protected routes

### Notes Management
- Create notes with title and content
- Edit existing notes
- Delete notes with confirmation
- View all notes in a responsive grid
- Search and filter capabilities

### UI/UX
- Clean, modern interface
- Smooth animations and transitions
- Responsive design (mobile, tablet, desktop)
- Loading states and error handling
- Empty state messaging

## Best Practices Implemented

- **Component-based Architecture**: Reusable, maintainable components
- **State Management**: Centralized state with Pinia
- **Error Handling**: Comprehensive error messages and user feedback
- **Security**: Secure token storage and API authentication
- **Performance**: Lazy loading, code splitting with Vite
- **Accessibility**: Semantic HTML and ARIA labels
- **Responsive Design**: Mobile-first approach

## Development

### Code Style

```bash
# Lint and fix code
npm run lint
```

### Building

```bash
# Production build
npm run build

# Output is in the `dist` directory
```

## Deployment

The built frontend can be deployed to any static hosting service:
- Netlify
- Vercel
- AWS S3 + CloudFront
- GitHub Pages

## License

See LICENSE file in the root directory.
