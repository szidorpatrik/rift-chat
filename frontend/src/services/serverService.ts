import { apiFetch } from '@utils/api'
import type { ServerResponse, ServerRequest } from '@types'

export const serverService = {
    async getServers() {
        return apiFetch<ServerResponse[]>('/servers')
    },

    async getPublicServers() {
        return apiFetch<ServerResponse[]>('/servers/public')
    },

    async joinPublicServer(serverId: string) {
        return apiFetch<ServerResponse>('/servers/join', {
            method: 'POST',
            body: JSON.stringify({ server_id: serverId }),
        })
    },

    async createServer(payload: ServerRequest) {
        return apiFetch<ServerResponse>('/servers', {
            method: 'POST',
            body: JSON.stringify(payload),
        })
    },

    async joinByInvite(inviteCode: string) {
        return apiFetch<ServerResponse>('/servers/join/invite', {
            method: 'POST',
            body: JSON.stringify({ invite_code: inviteCode }),
        })
    }
}
