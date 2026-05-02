import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { ChannelRequest } from '@types'
import { channelService } from '@services/channelService'
import { useServerStore } from '@stores/serverStore'

export const useChannelStore = defineStore('channel', () => {
    const sessionStorageKey = 'lastChannelId'
    const toast = useToast()

    // --- State ---
    const currentChannelId = ref<string | null>(null)
    const activeVoiceChannelId = ref<string | null>(null)
    const loadingChannels = ref<Set<string>>(new Set())

    // --- Getters ---
    const currentChannel = computed(() => {
        const serverStore = useServerStore()
        return serverStore.currentServer?.channels.find(c => c.id === currentChannelId.value) || null
    })
    const getChannel = computed(() => (channelId: string | null) => {
        if (!channelId) return null
        const serverStore = useServerStore()
        return serverStore.currentServer?.channels.find(c => c.id === channelId) || null
    })

    const getSessionKey = (serverId: string) => `${sessionStorageKey}_${serverId}`

    // --- Actions ---
    const setCurrentChannel = (serverId: string | null, channelId: string | null) => {
        const channel = getChannel.value(channelId)
        
        if (channel?.type === 1) {
            // VOICE CHANNEL LOGIC
            activeVoiceChannelId.value = channelId
            console.log(`Connecting to voice: ${channelId}`)
            return
        }
        currentChannelId.value = channelId
        // TEXT CHANNEL LOGIC
        if (serverId && channelId && channel?.type === 0) {
            sessionStorage.setItem(getSessionKey(serverId), channelId)
        }
    }

    const getPersistedChannelId = (serverId: string) => {
        return sessionStorage.getItem(getSessionKey(serverId))
    }

    const loadChannels = async (serverId: string) => {
        if (loadingChannels.value.has(serverId)) return false

        loadingChannels.value.add(serverId)
        const serverStore = useServerStore()

        if (serverStore.areChannelsLoaded(serverId)) {
            loadingChannels.value.delete(serverId)
            return true
        }

        try {
            const { data, error } = await channelService.getChannels(serverId)
            if (error) {
                toast.add({ title: 'Error', description: error.error, color: 'error' })
                return false
            }
            if (data) {
                serverStore.setServerChannels(serverId, data)
                serverStore.servers.find(s => s.id === serverId)!.channelsLoaded = true
                return true
            }
        } finally {
            loadingChannels.value.delete(serverId)
        }
        return false
    }

    const createChannel = async (payload: ChannelRequest) => {
        const serverStore = useServerStore()

        if (!serverStore.currentServer) {
            toast.add({ title: 'Error', description: 'No server selected', color: 'error' })
            return false
        }

        const position = serverStore.currentServer.channels.length || 0
        const { data, error } = await channelService.createChannel(serverStore.currentServer.id, payload, position)

        if (error) {
            toast.add({ title: 'Error', description: error.error, color: 'error' })
            return false
        }

        if (data) {
            serverStore.addChannelToServer(data.server_id, data)
            toast.add({ title: 'Channel created', description: `Channel "${payload.name}" has been created.`, color: 'success' })
            return true
        }
        return false
    }

    const deleteChannel = async (serverId: string, channelId: string) => {
        const serverStore = useServerStore()
        
        const { error } = await channelService.deleteChannel(serverId, channelId)

        if (error) {
            toast.add({ title: 'Error', description: error.error, color: 'error' })
            return false
        }

        serverStore.removeChannelFromServer(serverId, channelId)
        toast.add({ title: 'Channel deleted', description: `Channel has been deleted.`, color: 'success' })
        return true
    }

    return {
        currentChannelId,
        currentChannel,
        activeVoiceChannelId,
        getPersistedChannelId,
        setCurrentChannel,
        loadChannels,
        createChannel,
        deleteChannel,
    }
})