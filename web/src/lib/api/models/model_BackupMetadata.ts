/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { model_BackupStatus } from './model_BackupStatus';
export type model_BackupMetadata = {
	createdAt?: string;
	database?: string;
	error?: string;
	filePath?: string;
	fileSize?: number;
	host?: string;
	id?: string;
	objectKey?: string;
	/**
	 * pending, generating, completed, failed
	 */
	status?: model_BackupStatus;
	timestamp?: string;
	type?: string;
};
