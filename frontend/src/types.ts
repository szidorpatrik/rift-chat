export interface Server { id: string; name: string; avatar?: string }
export interface Channel { id: string; name: string; type: 0 | 1 | 2, position: number }
export interface Message { id: string; author: string; content: string }

export type LoginRequest = {
    username: string,
    password: string,
}

export type LoginResponse = {
    id: string,
    username: string,
    email: string,
    status: string,
    avatar_url?: string,
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

export type LogoutResponse = {
    // No fields needed for logout
}

export type ServerRequest = {
    name: string,
    is_public: boolean,
    description?: string,
}

export type ServerResponse = {
    id: string,
    name: string,
    owner_id: string,
    description?: string,
    is_public: boolean,
    icon_url?: string,
    background_url?: string,
    created_at: string,
    updated_at: string,
}

export type ErrorResponse = {
    error: string,
    status_code: number,
}

export type ChannelRequest = {
    name: string,
    type: 0 | 1 | 2, // 0 = text, 1 = voice, 2 = category
    position: number,
}

export type ChannelResponse = {
    id: string,
    server_id: string,
    name: string,
    type: 0 | 1 | 2, // 0 = text, 1 = voice, 2 = category
    position: number,
}
