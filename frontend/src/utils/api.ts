import type { ErrorResponse } from '@types'

export type ApiResult<T> = {
    data: T | null,
    error: ErrorResponse | null, 
}

export async function apiFetch<TResponse>(
    endpoint: string, 
    options: RequestInit = {}
): Promise<ApiResult<TResponse>> {
    try {
        const url = `/api/v1${endpoint}`
        
        const res = await fetch(url, {
            ...options,
            headers: {
                'Content-Type': 'application/json',
                ...options.headers,
            },
            credentials: 'include',
        })

        const result = await res.json()

        if (res.ok) {
            return { data: result as TResponse, error: null }
        } else {
            return { 
                data: null, 
                error: {
                    error: result.error || 'An unexpected error occurred',
                    status_code: result.status_code || res.status
                }
            }
        }
    } catch (err: any) {
        return { 
            data: null, 
            error: {
                error: err.message || 'Network Error',
                status_code: 0,
            }
        }
    }
}
