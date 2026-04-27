export interface Server { id: string; name: string; avatar?: string }
export interface Channel { id: string; name: string }
export interface Message { id: string; author: string; content: string }

export type LoginRequest = {
    username: string,
    password: string,
}

export type LoginResponse = {
    id: string,
    username: string,
    email: string,
    avatar_url?: string,
    status?: string,
    description?: string,
}

export type RegisterRequest = {
    username: string,
    email: string,
    password: string,
}

export type RegisterResponse = {
    id: string,
    username: string,
    email: string,
    created_at: string,
    updated_at: string,
}

export type ErrorResponse = {
    error: string,
    status_code: number,
}