import { DatabaseService } from '$lib/api';
import type { PageLoad } from './$types';

export const ssr = false;

export const load: PageLoad = async ({ url, depends }) => {
  depends('app:databases');

  const page = Number(url.searchParams.get('page')) || 1;
  const limit = Number(url.searchParams.get('limit')) || 10;

  const response = await DatabaseService.getDatabases(page, limit);

  return {
    databases: response.databases || [],
    pagination: {
      page: response.page,
      limit: response.limit,
      total: response.total
    }
  };
};
