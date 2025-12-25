/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
import type { database_BackupStats } from '../models/database_BackupStats';
import type { model_BackupListResponse } from '../models/model_BackupListResponse';
import type { model_BackupMetadata } from '../models/model_BackupMetadata';
import type { model_BackupRequest } from '../models/model_BackupRequest';
import type { model_BackupResponse } from '../models/model_BackupResponse';
export class BackupService {
	/**
	 * Trigger a database backup
	 * Queue a backup job for the specified database
	 * @param request Backup Request
	 * @returns model_BackupResponse Backup started
	 * @throws ApiError
	 */
	public static postBackup(request: model_BackupRequest): CancelablePromise<model_BackupResponse> {
		return __request(OpenAPI, {
			method: 'POST',
			url: '/backup',
			body: request,
			errors: {
				400: `error: Bad request`
			}
		});
	}
	/**
	 * List all backups
	 * Get a paginated list of all backups with optional filtering
	 * @param page Page number
	 * @param limit Items per page
	 * @param statuses Comma-separated status values (pending,generating,completed,failed)
	 * @param search Search keyword (searches in database, host, type)
	 * @param orderBy Field to order by
	 * @param orderDir Order direction (asc/desc)
	 * @param startDate Start date for filtering (RFC3339 format)
	 * @param endDate End date for filtering (RFC3339 format)
	 * @returns model_BackupListResponse List of backups
	 * @throws ApiError
	 */
	public static getBackups(
		page: number = 1,
		limit: number = 10,
		statuses?: string,
		search?: string,
		orderBy: string = 'createdAt',
		orderDir: string = 'desc',
		startDate?: string,
		endDate?: string
	): CancelablePromise<model_BackupListResponse> {
		return __request(OpenAPI, {
			method: 'GET',
			url: '/backups',
			query: {
				page: page,
				limit: limit,
				statuses: statuses,
				search: search,
				orderBy: orderBy,
				orderDir: orderDir,
				startDate: startDate,
				endDate: endDate
			},
			errors: {
				500: `error: Internal server error`
			}
		});
	}
	/**
	 * Get backup statistics
	 * Retrieve aggregated backup statistics by type and status
	 * @param startDate Start date for filtering (RFC3339 format)
	 * @param endDate End date for filtering (RFC3339 format)
	 * @returns database_BackupStats Backup statistics
	 * @throws ApiError
	 */
	public static getBackupsStats(
		startDate?: string,
		endDate?: string
	): CancelablePromise<database_BackupStats> {
		return __request(OpenAPI, {
			method: 'GET',
			url: '/backups/stats',
			query: {
				startDate: startDate,
				endDate: endDate
			},
			errors: {
				500: `error: Internal server error`
			}
		});
	}
	/**
	 * Get a single backup
	 * Retrieve a single backup by ID
	 * @param id Backup ID
	 * @returns model_BackupMetadata Backup details
	 * @throws ApiError
	 */
	public static getBackupById(id: string): CancelablePromise<model_BackupMetadata> {
		return __request(OpenAPI, {
			method: 'GET',
			url: '/backups/{id}',
			path: {
				id: id
			},
			errors: {
				400: `error: Bad request`,
				404: `error: Backup not found`,
				500: `error: Internal server error`
			}
		});
	}
	/**
	 * Delete a backup
	 * Delete a backup from both MongoDB and R2 storage
	 * @param id Backup ID
	 * @returns model_BackupResponse Backup deleted successfully
	 * @throws ApiError
	 */
	public static deleteBackups(id: string): CancelablePromise<model_BackupResponse> {
		return __request(OpenAPI, {
			method: 'DELETE',
			url: '/backups/{id}',
			path: {
				id: id
			},
			errors: {
				400: `error: Bad request`,
				404: `error: Backup not found`,
				500: `error: Internal server error`
			}
		});
	}
	/**
	 * Download a backup file
	 * Generate a presigned URL to download a backup file from R2 storage
	 * @param id Backup ID
	 * @returns string Download URL
	 * @throws ApiError
	 */
	public static getBackupsDownload(id: string): CancelablePromise<Record<string, string>> {
		return __request(OpenAPI, {
			method: 'GET',
			url: '/backups/{id}/download',
			path: {
				id: id
			},
			errors: {
				400: `error: Bad request`,
				404: `error: Backup not found`,
				500: `error: Internal server error`
			}
		});
	}
}
