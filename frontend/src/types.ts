export type LoginRequest = {
    username: string,
    password: string,
}

export type LoginResponse = {
    id: string,
    username: string,
    email: string,
}

export type ErrorResponse = {
    error: string,
    status_code: number,
}