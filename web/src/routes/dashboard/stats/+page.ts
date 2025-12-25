import { BackupService } from '$lib/api';
import type { PageLoad } from './$types';

export const ssr = false;

export const load: PageLoad = async ({ url }) => {
	const startDate = url.searchParams.get('startDate') || undefined;
	const endDate = url.searchParams.get('endDate') || undefined;

	const stats = await BackupService.getBackupsStats(startDate, endDate);

	return {
		stats
	};
};
