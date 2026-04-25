export type LoginRequest = {
    username: string,
    password: string,
}

export type LoginResponse = {
    id: string,
    username: string,
    email: string,
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