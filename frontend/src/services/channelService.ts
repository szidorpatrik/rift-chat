import { apiFetch } from '@utils/api'
import type { ChannelResponse, ChannelRequest } from '@types'

export const channelService = {
    async getChannels(serverId: string) {
        return apiFetch<ChannelResponse[]>(`/servers/${serverId}/channels`)
    },

    async createChannel(serverId: string, payload: ChannelRequest, position: number) {
        return apiFetch<ChannelResponse>(`/servers/${serverId}/channels`, {
            method: 'POST',
            body: JSON.stringify({
                name: payload.name,
                type: payload.type,
                position,
            }),
        })
    },

    async deleteChannel(serverId: string, channelId: string) {
        return apiFetch(`/servers/${serverId}/channels/${channelId}`, {
            method: 'DELETE',
        })
    }
}

