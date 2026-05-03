import type { UUID } from "crypto";

export interface Server { id: UUID; name: string; avatar?: string }
export interface Channel { id: UUID; name: string; type: 0 | 1 | 2, position: number }
export interface Message { id: UUID; author: string; content: string }

export type Status = 'online' | 'offline' | 'idle' | 'dnd'
export type ChipColor = 'primary' | 'neutral' | 'warning' | 'error'

export const statusColorMap: Record<Status, ChipColor> = {
    online: 'primary',
    offline: 'neutral',
    idle: 'warning',
    dnd: 'error',
}


export type LoginRequest = {
    username: string,
    password: string,
}

export type LoginResponse = {
    id: UUID,
    username: string,
    email: string,
    status: Status,
    avatar_url?: string,
    description?: string,
}

export type RegisterRequest = {
    username: string,
    email: string,
    password: string,
}

export type RegisterResponse = {
    id: UUID,
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
    id: UUID,
    name: string,
    owner_id: UUID,
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
    id: UUID,
    server_id: UUID,
    name: string,
    type: 0 | 1 | 2, // 0 = text, 1 = voice, 2 = category
    position: number,
}

export type MemberResponse = {
    id: UUID,
    username: string,
    status: Status,
    avatar_url?: string,
    description?: string,
}