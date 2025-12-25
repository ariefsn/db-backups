/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { model_BackupResponse } from '../models/model_BackupResponse';
import type { model_CreateDatabaseRequest } from '../models/model_CreateDatabaseRequest';
import type { model_Database } from '../models/model_Database';
import type { model_DatabaseListResponse } from '../models/model_DatabaseListResponse';
import type { model_UpdateDatabaseRequest } from '../models/model_UpdateDatabaseRequest';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class DatabaseService {
    /**
     * List all databases
     * List all saved database configurations
     * @param page Page number
     * @param limit Items per page
     * @returns model_DatabaseListResponse OK
     * @throws ApiError
     */
    public static getDatabases(
        page?: number,
        limit?: number,
    ): CancelablePromise<model_DatabaseListResponse> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/databases',
            query: {
                'page': page,
                'limit': limit,
            },
            errors: {
                500: `error: Internal server error`,
            },
        });
    }
    /**
     * Create a new database
     * Save a new database configuration
     * @param request Database Configuration
     * @returns model_Database Created
     * @throws ApiError
     */
    public static postDatabases(
        request: model_CreateDatabaseRequest,
    ): CancelablePromise<model_Database> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/databases',
            body: request,
            errors: {
                400: `error: Bad request`,
                500: `error: Internal server error`,
            },
        });
    }
    /**
     * Get a database
     * Retrieve a database configuration by ID
     * @param id Database ID
     * @returns model_Database OK
     * @throws ApiError
     */
    public static getDatabases1(
        id: string,
    ): CancelablePromise<model_Database> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/databases/{id}',
            path: {
                'id': id,
            },
            errors: {
                400: `error: Bad request`,
                404: `error: Database not found`,
                500: `error: Internal server error`,
            },
        });
    }
    /**
     * Update a database
     * Update an existing database configuration
     * @param id Database ID
     * @param request Database Configuration
     * @returns model_Database OK
     * @throws ApiError
     */
    public static putDatabases(
        id: string,
        request: model_UpdateDatabaseRequest,
    ): CancelablePromise<model_Database> {
        return __request(OpenAPI, {
            method: 'PUT',
            url: '/databases/{id}',
            path: {
                'id': id,
            },
            body: request,
            errors: {
                400: `error: Bad request`,
                404: `error: Database not found`,
                500: `error: Internal server error`,
            },
        });
    }
    /**
     * Delete a database
     * Delete a database configuration
     * @param id Database ID
     * @returns model_BackupResponse Database deleted successfully
     * @throws ApiError
     */
    public static deleteDatabases(
        id: string,
    ): CancelablePromise<model_BackupResponse> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/databases/{id}',
            path: {
                'id': id,
            },
            errors: {
                400: `error: Bad request`,
                404: `error: Database not found`,
                500: `error: Internal server error`,
            },
        });
    }
    /**
     * Trigger a backup for a database
     * Manually trigger a backup for a saved database configuration
     * @param id Database ID
     * @returns model_BackupResponse Backup job submitted successfully
     * @throws ApiError
     */
    public static postDatabasesBackup(
        id: string,
    ): CancelablePromise<model_BackupResponse> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/databases/{id}/backup',
            path: {
                'id': id,
            },
            errors: {
                400: `error: Bad request`,
                404: `error: Database not found`,
                500: `error: Internal server error`,
            },
        });
    }
}
