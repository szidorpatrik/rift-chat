import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { ServerResponse, ChannelResponse } from '@types'
import { apiFetch } from '@utils/api'

export interface Server extends ServerResponse {
    channels: ChannelResponse[]
}

export const useServerStore = defineStore('server', () => {
    const sessionStorageKey = 'lastServerId'
    const toast = useToast()

    // --- State ---
    const servers = ref<Server[]>([])
    const publicServers = ref<ServerResponse[]>([])

    const currentServerId = ref<string | null>(sessionStorage.getItem(sessionStorageKey))
    const currentChannelId = ref<string | null>(null)

    // --- Getters ---
    const currentServer = computed(() => servers.value.find(s => s.id === currentServerId.value) || null)
    const currentChannel = computed(() => currentServer.value?.channels.find(c => c.id === currentChannelId.value) || null)

    // --- Actions ---
    const setServers = (serverArray: ServerResponse[]) => {
        servers.value = serverArray.map(server => ({
            ...server,
            channels: []
        }))

        if (!currentServerId.value && servers.value.length > 0) {
            setCurrentServer(servers.value[0].id)
        }
    }

    const addServer = (server: ServerResponse) => {
        servers.value.push({ ...server, channels: [] })
        setCurrentServer(server.id)
    }

    const removeServer = (serverId: string) => {
        servers.value = servers.value.filter(s => s.id !== serverId)
        if (currentServerId.value === serverId) {
            setCurrentServer(servers.value.length > 0 ? servers.value[0].id : null)
        }
    }

    const setCurrentServer = (serverId: string | null) => {
        currentServerId.value = serverId
        if (serverId) {
            sessionStorage.setItem(sessionStorageKey, serverId)
        } else {
            sessionStorage.removeItem(sessionStorageKey)
            currentChannelId.value = null
        }
    }

    const setCurrentChannel = (channelId: string | null) => {
        currentChannelId.value = channelId
    }

    const loadServers = async () => {
        const { data, error } = await apiFetch<ServerResponse[]>('/servers')
        if (error) {
            toast.add({ title: 'Error', description: error.error, color: 'error' })
            return false
        }

        if (data) {
            setServers(data)
            return true
        }
        return false
    }

    const loadChannels = async (serverId: string) => {
        const { data, error } = await apiFetch<ChannelResponse[]>(`/servers/${serverId}/channels`)
        if (error) {
            toast.add({ title: 'Error', description: error.error, color: 'error' })
            return false
        }

        if (data) {
            const serverIndex = servers.value.findIndex(s => s.id === serverId)
            if (serverIndex !== -1) {
                servers.value[serverIndex].channels = data
            }
            return true
        }
        return false
    }

    const loadPublicServers = async () => {
        const { data, error } = await apiFetch<ServerResponse[]>('/servers/public')
        if (error) {
            toast.add({ title: 'Error', description: error.error, color: 'error' })
            return false
        }

        if (data) {
            publicServers.value = data
            return true
        }
        return false
    }

    const addPublicServer = (server: ServerResponse) => {
        publicServers.value.push(server)
    }

    const joinPublicServer = async (serverId: string) => {
        const { data, error } = await apiFetch<ServerResponse>('/servers/join', {
            method: 'POST',
            body: JSON.stringify({ server_id: serverId }),
        })

        if (error) {
            toast.add({ title: 'Failed to join server', description: error.error, color: 'error' })
            return false
        }

        if (data) {
            addServer(data)
            return true
        }
        return false
    }

    const createServer = async (payload: { name: string; is_public?: boolean; description?: string }) => {
        const { data, error } = await apiFetch<ServerResponse>('/servers', {
            method: 'POST',
            body: JSON.stringify({
                name: payload.name,
                is_public: payload.is_public,
                description: payload.description,
            }),
        })

        if (error) return { error }

        if (data) {
            addServer(data)
            if (data.is_public) {
                addPublicServer(data)
            }
            return { data }
        }
        return  { error } 
    }

    const joinServerByInvite = async (inviteCode: string) => {
        const { data, error } = await apiFetch<ServerResponse>('/servers/join/invite', {
            method: 'POST',
            body: JSON.stringify({ invite_code: inviteCode }),
        })

        if (error) return { error }

        if (data) {
            addServer(data)
            return { data }
        }
        return { error }
    }

    return {
        // State
        servers,
        publicServers,
        currentServerId,
        currentChannelId,
        // Getters
        currentServer,
        currentChannel,
        // Actions
        addServer,
        setServers,
        removeServer,
        setCurrentServer,
        setCurrentChannel,
        loadServers,
        loadChannels,
        loadPublicServers,
        addPublicServer,
        joinPublicServer,
        createServer,
        joinServerByInvite,
    }
})
