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
        
        const fetchConfig: RequestInit = {
            ...options,
            headers: {
                'Content-Type': 'application/json',
                ...options.headers,
            },
            credentials: 'include',
        }

        // Initial request
        let res = await fetch(url, fetchConfig)

        // If it fails with 401, try to refresh
        if (res.status === 401 && endpoint !== '/login' && endpoint !== '/refresh') {
            const refreshRes = await fetch('/api/v1/refresh', { 
                method: 'POST', 
                credentials: 'include' 
            })

            if (refreshRes.ok) {
                res = await fetch(url, fetchConfig)
            } else {
                // The refresh token is expired or missing. 
                // Let it fall through to return the 401 to the component.
            }
        }

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
