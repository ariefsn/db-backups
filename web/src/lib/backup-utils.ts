import type { model_BackupMetadata } from '$lib/api';

/**
 * Label for a backup, falling back through the fields most likely to be set.
 * The backend fills `name` for new backups; the fallbacks keep older records
 * that were saved before that readable.
 */
export function primaryLabel(backup: model_BackupMetadata): string {
	return backup.name || backup.database || backup.host || '-';
}
