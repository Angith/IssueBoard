# IssueBoard Frontend

This is the frontend application for IssueBoard, built with Next.js 16 and React 19. It provides an interactive interface to manage your GitHub repository inventory and view issues grouped by labels.

## 📋 Prerequisites

- **Node.js**: v20.x or later
- **npm**: v10.x or later
- **Supabase Project**: For authentication and database access

## ⚙️ Configuration

Create a `.env.local` file in this directory with the following variables:

```env
NEXT_PUBLIC_SUPABASE_URL=your_supabase_url
NEXT_PUBLIC_SUPABASE_ANON_KEY=your_supabase_anon_key
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## 🚀 Setup

### Frontend Setup
1. Install dependencies:
   ```bash
   npm install
   ```
2. Start the development server:
   ```bash
   npm run dev
   ```
   The application will be available at `http://localhost:3000`.

## 🧪 Testing

Run frontend linting and type checking:
```bash
# Linting
npm run lint

# Type checking (if applicable)
npx tsc --noEmit
```

*(Note: Automated unit and E2E tests are currently in the roadmap.)*

## 🎨 Tech Stack

- **Framework**: [Next.js 16 (App Router)](https://nextjs.org/)
- **Library**: [React 19](https://react.dev/)
- **Styling**: [Tailwind CSS 4](https://tailwindcss.com/)
- **Auth & DB**: [Supabase SDK](https://supabase.com/docs/reference/javascript/introduction)
- **Icons**: [Lucide React](https://lucide.dev/) (if used)

## 📁 Project Structure

- `src/app/`: Next.js App Router pages and layouts.
  - `login/`: GitHub OAuth login page.
  - `inventory/`: Management of tracked repositories.
  - `repos/[id]/`: Detailed issue board view for a specific repository.
- `src/components/`: Reusable UI components.
- `src/hooks/`: Custom React hooks for data fetching and state.
- `src/lib/`: Shared utilities and configurations (e.g., Supabase client).
- `src/services/`: API client services for interacting with the backend.
