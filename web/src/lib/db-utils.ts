/**
 * Parses a database connection string (URI) into its components.
 * Supports PostgreSQL, MySQL, MongoDB, and Redis.
 */
export function parseConnectionString(uri: string) {
  try {
    // Basic format: scheme://[user[:password]@]host[:port][/database][?options]
    const url = new URL(uri);
    const scheme = url.protocol.replace(':', '').toLowerCase();

    let type = '';
    let defaultPort = '';

    switch (scheme) {
      case 'postgresql':
      case 'postgres':
        type = 'postgre';
        defaultPort = '5432';
        break;
      case 'mysql':
        type = 'mysql';
        defaultPort = '3306';
        break;
      case 'mongodb':
      case 'mongodb+srv':
        type = 'mongodb';
        defaultPort = '27017';
        break;
      case 'redis':
      case 'rediss':
        type = 'redis';
        defaultPort = '6379';
        break;
      default:
        return null;
    }

    const username = decodeURIComponent(url.username || '');
    const password = decodeURIComponent(url.password || '');
    const host = url.hostname || '';
    const port = url.port || defaultPort;

    // For MongoDB, the database name can be in the path, but sometimes it's empty
    // For others, it's usually the first segment of the path
    let database = url.pathname.replace(/^\//, '');

    // If there's a slash in the database name (e.g. for options), take only the first part
    if (database.includes('/')) {
      database = database.split('/')[0];
    }

    return {
      type,
      host,
      port,
      username,
      password,
      database
    };
  } catch (e) {
    console.error('Failed to parse connection string:', e);
    return null;
  }
}
