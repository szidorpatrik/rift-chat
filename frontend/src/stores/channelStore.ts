import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { ChannelRequest } from '@types'
import { channelService } from '@services/channelService'
import { useServerStore } from '@stores/serverStore'

export const useChannelStore = defineStore('channel', () => {
    const toast = useToast()

    // --- State ---
    const currentChannelId = ref<string | null>(null)
    const loadingChannels = ref<Set<string>>(new Set())

    // --- Getters ---
    const currentChannel = computed(() => {
        const serverStore = useServerStore()
        return serverStore.currentServer?.channels.find(c => c.id === currentChannelId.value) || null
    })

    // --- Actions ---
    const setCurrentChannel = (channelId: string | null) => {
        currentChannelId.value = channelId
    }

    const loadChannels = async (serverId: string) => {
        if (loadingChannels.value.has(serverId)) return false

        loadingChannels.value.add(serverId)
        const serverStore = useServerStore()

        try {
            const { data, error } = await channelService.getChannels(serverId)
            if (error) {
                toast.add({ title: 'Error', description: error.error, color: 'error' })
                return false
            }
            if (data) {
                serverStore.setServerChannels(serverId, data)
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

    return {
        currentChannelId,
        currentChannel,
        setCurrentChannel,
        loadChannels,
        createChannel
    }
})