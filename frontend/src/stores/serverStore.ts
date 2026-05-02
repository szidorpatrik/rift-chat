import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { ServerResponse, ChannelResponse, ChannelRequest, ServerRequest } from '@types'
import { serverService } from '@services/serverService'
import { channelService } from '@services/channelService'

export interface Server extends ServerResponse {
    channels: ChannelResponse[]
}

export const useServerStore = defineStore('server', () => {
    const sessionStorageKey = 'lastServerId'
    const toast = useToast()

    // --- State ---
    const servers = ref<Server[]>([])
    const publicServers = ref<ServerResponse[]>([])

    const currentServerId = ref<string | null>(null)
    const currentChannelId = ref<string | null>(null)
    const loadingChannels = ref<Set<string>>(new Set())

    // --- Getters ---
    const currentServer = computed(() => servers.value.find(s => s.id === currentServerId.value) || null)
    const currentChannel = computed(() => currentServer.value?.channels.find(c => c.id === currentChannelId.value) || null)

    // --- Actions (State Mutations) ---
    const setServers = async (serverArray: ServerResponse[]) => {
        servers.value = serverArray.map(server => ({ ...server, channels: [] }))

        if (servers.value.length === 0) {
            setCurrentServer(null)
            return
        }
        const lastServerId = sessionStorage.getItem(sessionStorageKey) || servers.value[0].id
        setCurrentServer(lastServerId)
        await loadChannels(lastServerId)
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

    const addPublicServer = (server: ServerResponse) => {
        publicServers.value.push(server)
    }

    const setCurrentServer = async (serverId: string | null) => {
        currentServerId.value = serverId
        currentChannelId.value = null

        if (!serverId) {
            sessionStorage.removeItem(sessionStorageKey)
            return
        }

        sessionStorage.setItem(sessionStorageKey, serverId)

        const server = servers.value.find(s => s.id === serverId)
        if (!server) return

        if (server.channels.length === 0) {
            await loadChannels(serverId)
        }

        currentChannelId.value = server.channels?.[0]?.id ?? null
    }

    const setCurrentChannel = (channelId: string | null) => {
        currentChannelId.value = channelId
    }

    // --- Actions (API Integrations) ---
    const loadServers = async () => {
        const { data, error } = await serverService.getServers()
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
        if (loadingChannels.value.has(serverId)) return false
        loadingChannels.value.add(serverId)

        try {
            const { data, error } = await channelService.getChannels(serverId)
            if (error) {
                toast.add({ title: 'Error', description: error.error, color: 'error' })
                return false
            }
            if (data) {
                const server = servers.value.find(s => s.id === serverId)
                if (server) server.channels = data
                return true
            }
        } finally {
            loadingChannels.value.delete(serverId)
        }
        return false
    }

    const createChannel = async (payload: ChannelRequest) => {
        if (!currentServer.value) {
            toast.add({ title: 'Error', description: 'No server selected', color: 'error' })
            return false
        }

        const position = currentServer.value.channels.length || 0
        const { data, error } = await channelService.createChannel(currentServer.value.id, payload, position)

        if (error) {
            toast.add({ title: 'Error', description: error.error, color: 'error' })
            return false
        }

        if (data) {
            const server = servers.value.find(s => s.id === data.server_id)
            if (server) {
                server.channels.push(data)
                toast.add({ title: 'Channel created', description: `Channel "${payload.name}" has been created.`, color: 'success' })
                return true
            }
        }
        return false
    }

    const loadPublicServers = async () => {
        const { data, error } = await serverService.getPublicServers()
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

    const joinPublicServer = async (serverId: string) => {
        const { data, error } = await serverService.joinPublicServer(serverId)
        if (error) {
            toast.add({ title: 'Failed to join', description: error.error, color: 'error' })
            return false
        }
        if (data) {
            addServer(data)
            return true
        }
        return false
    }

    const createServer = async (payload: ServerRequest) => {
        const { data, error } = await serverService.createServer(payload)
        if (error) return { error }
        if (data) {
            addServer(data)
            if (data.is_public) addPublicServer(data)
            return { data }
        }
        return { error }
    }

    const joinServerByInvite = async (inviteCode: string) => {
        const { data, error } = await serverService.joinByInvite(inviteCode)
        if (error) return { error }
        if (data) {
            addServer(data)
            return { data }
        }
        return { error }
    }

    return {
        servers,
        publicServers,
        currentServerId,
        currentChannelId,
        currentServer,
        currentChannel,
        addServer,
        setServers,
        removeServer,
        setCurrentServer,
        setCurrentChannel,
        loadServers,
        loadChannels,
        createChannel,
        loadPublicServers,
        addPublicServer,
        joinPublicServer,
        createServer,
        joinServerByInvite,
    }
})
