import { BackupService } from '$lib/api';
import type { PageLoad } from './$types';

export const ssr = false;

export const load: PageLoad = async ({ url, depends }) => {
  depends('app:backups');

  const page = Number(url.searchParams.get('page')) || 1;
  const limit = Number(url.searchParams.get('limit')) || 10;
  const statuses = url.searchParams.get('statuses') || undefined;
  const types = url.searchParams.get('types') || undefined;
  const search = url.searchParams.get('search') || undefined;
  const orderBy = url.searchParams.get('orderBy') || 'createdAt';
  const orderDir = url.searchParams.get('orderDir') || 'desc';
  const startDate = url.searchParams.get('startDate') || undefined;
  const endDate = url.searchParams.get('endDate') || undefined;

  const response = await BackupService.getBackups(
    page,
    limit,
    statuses,
    search,
    orderBy,
    orderDir,
    startDate,
    endDate
  );

  return {
    backups: response.backups || [],
    pagination: {
      page: response.page,
      limit: response.limit,
      total: response.total
    }
  };
};
