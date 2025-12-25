/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { model_BackupType } from './model_BackupType';
export type model_Database = {
    connectionUri?: string;
    createdAt?: string;
    cronExpression?: string;
    database?: string;
    host?: string;
    id?: string;
    isActive?: boolean;
    name?: string;
    password?: string;
    port?: string;
    type?: model_BackupType;
    updatedAt?: string;
    username?: string;
    webhookUrl?: string;
};

