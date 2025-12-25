/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { model_BackupType } from './model_BackupType';
export type model_BackupRequest = {
	authSource?: string;
	connectionUri?: string;
	database?: string;
	host?: string;
	password?: string;
	port?: string;
	type?: model_BackupType;
	username?: string;
	webhookUrl?: string;
};
