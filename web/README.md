# Db Backup Web Interface

A modern, user-friendly web interface for managing your database backups. Built with SvelteKit and Shadcn UI.

## Screenshots

### Dashboard Statistics
![Stats](../assets/%5BDBB%5D%20Stats.png)

### Database Management
![My Databases](../assets/%5BDBB%5D%20My%20Databases.png)

### Backup History
![Backups](../assets/%5BDBB%5D%20Backups.png)

### create Backup Form
![Backup Form](../assets/%5BDBB%5D%20Backups%20Form.png)

## Features

- **Multi-Database Support**: Create backups for PostgreSQL, MySQL, MongoDB, and Redis.
- **My Databases**: Manage saved database configurations for quick backup triggering and automation.
- **Automated Backups**: Configure cron schedules for each database configuration.
- **Flexible Connection Methods**: Support for both standard connection strings and manual host/user/pass forms.
- **Dashboard**: Monitor all your backups in one place with real-time status updates.
- **Backup Management**:
  - Create new backups on-demand using saved configurations or manual details.
  - Download backup files directly (via presigned URLs).
  - Delete old backups (removes from both metadata and object storage).
- **Statistics**: View reports and statistics about your backup history.
- **Object Storage Integration**: Seamlessly works with S3-compatible storage (Cloudflare R2, AWS S3, etc.).

## Tech Stack

- **Framework**: [SvelteKit](https://kit.svelte.dev/)
- **UI Components**: [Shadcn UI](https://www.shadcn-svelte.com/)
- **Styling**: [Tailwind CSS](https://tailwindcss.com/)
- **Icons**: [Lucide Svelte](https://lucide.dev/guide/packages/lucide-svelte)
- **API Client**: Generated from OpenAPI/Swagger definition.

## Getting Started

### Prerequisites

- Node.js (v18+)
- Bun (optional, but recommended)

### Installation

```bash
# Install dependencies
npm install
# or
bun install
```

### Developing

Start the development server:

```bash
npm run dev
# or
bun run dev
```

The application will be available at `http://localhost:5173`.

### Building

To create a production version of your app:

```sh
npm run build
```

You can preview the production build with `npm run preview`.

## Docker Support

You can run the web application using Docker Compose from the root directory:

```bash
docker compose up -d
```

## Authentication

**Important:** This application currently does not implement its own authentication mechanism. It is intended to be run in a secure environment.

If you are exposing this application to the public internet or an untrusted network, it is highly recommended to use a reverse proxy with authentication support.

For example, if you are using [Caddy](https://caddyserver.com/), you can easily add Basic Authentication to your Caddyfile:

```caddyfile
example.com {
    reverse_proxy localhost:3000
    basicauth {
        # Username "admin", Password "password" (hashed)
        # Use `caddy hash-password` to generate the hash
        admin $2a$14$1...
    }
}
```
