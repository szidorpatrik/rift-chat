import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { ServerResponse, ChannelResponse, ServerRequest } from '@types'
import { serverService } from '@services/serverService'
import { usePublicServerStore } from '@stores/publicServerStore'
import { useChannelStore } from '@stores/channelStore'

export interface Server extends ServerResponse {
    channels: ChannelResponse[],
    channelsLoaded: boolean
}

export const useServerStore = defineStore('server', () => {
    const sessionStorageKey = 'lastServerId'
    const toast = useToast()

    // --- State ---
    const servers = ref<Server[]>([])
    const currentServerId = ref<string | null>(null)

    // --- Getters ---
    const currentServer = computed(() => servers.value.find(s => s.id === currentServerId.value) || null)
    const areChannelsLoaded = computed(() => (serverId: string) => {
        const server = servers.value.find(s => s.id === serverId)
        return server ? server.channelsLoaded : false
    })

    // --- Actions (Mutations) ---
    const setServerChannels = (serverId: string, channels: ChannelResponse[]) => {
        const server = servers.value.find(s => s.id === serverId)
        if (server) server.channels = channels
    }

    const addChannelToServer = (serverId: string, channel: ChannelResponse) => {
        const server = servers.value.find(s => s.id === serverId)
        if (server) server.channels.push(channel)
    }

    const setServers = async (serverArray: ServerResponse[]) => {
        servers.value = serverArray.map(server => ({ ...server, channels: [], channelsLoaded: false }))

        if (servers.value.length === 0) {
            setCurrentServer(null)
            return
        }

        const lastServerId = sessionStorage.getItem(sessionStorageKey) || servers.value[0].id
        setCurrentServer(lastServerId)

        const channelStore = useChannelStore()
        await channelStore.loadChannels(lastServerId)
    }

    const addServer = (server: ServerResponse) => {
        servers.value.push({ ...server, channels: [], channelsLoaded: false })
        setCurrentServer(server.id)
    }

    const removeServer = (serverId: string) => {
        servers.value = servers.value.filter(s => s.id !== serverId)
        if (currentServerId.value === serverId) {
            setCurrentServer(servers.value.length > 0 ? servers.value[0].id : null)
        }
    }

    const setCurrentServer = async (serverId: string | null) => {
        const channelStore = useChannelStore()

        currentServerId.value = serverId
        channelStore.setCurrentChannel(null)

        if (!serverId) {
            sessionStorage.removeItem(sessionStorageKey)
            return
        }

        sessionStorage.setItem(sessionStorageKey, serverId)

        const server = servers.value.find(s => s.id === serverId)
        if (!server) return

        if (server.channels.length === 0) {
            await channelStore.loadChannels(serverId)
        }

        channelStore.setCurrentChannel(server.channels?.[0]?.id ?? null)
    }

    // --- Actions (API Integrations) ---
    const loadServers = async () => {
        const { data, error } = await serverService.getServers()
        if (error) {
            toast.add({ title: 'Error', description: error.error, color: 'error' })
            return false
        }
        if (data) {
            await setServers(data)
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
        const publicServerStore = usePublicServerStore()
        const { data, error } = await serverService.createServer(payload)

        if (error) return { error }
        if (data) {
            addServer(data)
            if (data.is_public) {
                publicServerStore.addServer(data)
            }
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
        currentServerId,
        currentServer,
        areChannelsLoaded,
        setServerChannels,
        addChannelToServer,
        setServers,
        addServer,
        removeServer,
        setCurrentServer,
        loadServers,
        joinPublicServer,
        createServer,
        joinServerByInvite
    }
})
