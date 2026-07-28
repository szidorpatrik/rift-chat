import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { ServerResponse, ChannelResponse, ServerRequest, MemberResponse } from '@types'
import { serverService } from '@services/serverService'
import { usePublicServerStore } from '@stores/publicServerStore'
import { useChannelStore } from '@stores/channelStore'
import { useMemberStore } from '@stores/memberStore'

export interface Server extends ServerResponse {
    channels: ChannelResponse[],
    members: MemberResponse[],
    channelsLoaded: boolean,
    membersLoaded: boolean,
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
        const server = getServer(serverId)
        return server ? server.channelsLoaded : false
    })
    const areMembersLoaded = computed(() => (serverId: string) => {
        const server = getServer(serverId)
        return server ? server.membersLoaded : false
    })

    // --- Actions (Mutations) ---
    const setServerChannels = (serverId: string, channels: ChannelResponse[]) => {
        const server = getServer(serverId)
        if (server) server.channels = channels
    }

    const addChannelToServer = (serverId: string, channel: ChannelResponse) => {
        const server = getServer(serverId)
        if (server) server.channels.push(channel)
    }

    const setServers = async (serverArray: ServerResponse[]) => {
        servers.value = serverArray.map(server => ({
            ...server,
            channels: [],
            channelsLoaded: false,
            members: [],
            membersLoaded: false
        }))

        if (servers.value.length === 0) {
            await setCurrentServer(null)
            return
        }

        const lastServerId = sessionStorage.getItem(sessionStorageKey) || servers.value[0].id
        const lastServer = getServer(lastServerId)

        if (lastServer) {
            await setCurrentServer(lastServer.id)
        } else {
            await setCurrentServer(servers.value[0].id)
            console.log(servers.value[0].id)
        }
    }

    const addServer = (server: ServerResponse) => {
        servers.value.push({ ...server, channels: [], channelsLoaded: false, members: [], membersLoaded: false })
        setCurrentServer(server.id)
    }

    const removeServer = (serverId: string) => {
        servers.value = servers.value.filter(s => s.id !== serverId)
        if (currentServerId.value === serverId) {
            setCurrentServer(servers.value.length > 0 ? servers.value[0].id : null)
        }
    }

    const removeChannelFromServer = (serverId: string, channelId: string) => {
        const server = getServer(serverId)
        if (server) {
            server.channels = server.channels.filter(c => c.id !== channelId)
        }
    }

    const getServer = (serverId: string) => {
        return servers.value.find(s => s.id === serverId) || null
    }

    const setCurrentServer = async (serverId: string | null) => {
        const channelStore = useChannelStore()
        const memberStore = useMemberStore()

        currentServerId.value = serverId

        if (!serverId) {
            sessionStorage.removeItem('lastServerId')
            return
        }

        sessionStorage.setItem('lastServerId', serverId)
        const server = getServer(serverId)

        if (!server) return

        if (!server.channelsLoaded) {
            await channelStore.loadChannels(serverId)
        }

        if (!server.membersLoaded) {
            await memberStore.loadMembers(serverId)
        }

        // Attempt to restore the last TEXT channel
        const savedId = channelStore.getPersistedChannelId(serverId)
        const textChannels = server.channels.filter(c => c.type === 0)

        const isValidSavedChannel = server.channels.some(c => c.id === savedId && c.type === 0)

        if (savedId && isValidSavedChannel) {
            channelStore.setCurrentChannel(serverId, savedId)
        } else if (textChannels.length > 0) {
            // Fallback to the first available text channel
            channelStore.setCurrentChannel(serverId, textChannels[0].id)
        }
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
        areMembersLoaded,
        setServerChannels,
        addChannelToServer,
        setServers,
        addServer,
        getServer,
        removeServer,
        setCurrentServer,
        loadServers,
        joinPublicServer,
        createServer,
        joinServerByInvite,
        removeChannelFromServer,
    }
})
