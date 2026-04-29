import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { ServerResponse } from '@types'
import { apiFetch } from '@/utils/api'

export const useServerStore = defineStore('server', () => {
    const toast = useToast()
    const servers = ref<ServerResponse[]>([])
    const publicServers = ref<ServerResponse[]>([])

    const addServer = (server: ServerResponse) => {
        servers.value.push(server)
    }

    const removeServer = (serverId: string) => {
        servers.value = servers.value.filter(s => s.id !== serverId)
    }

    const setServers = (serverArray: ServerResponse[]) => {
        servers.value = serverArray
    }

    const setPublicServers = (serverArray: ServerResponse[]) => {
        publicServers.value = serverArray
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

    const loadPublicServers = async () => {
        const { data, error } = await apiFetch<ServerResponse[]>('/servers/public')
        if (error) {
            toast.add({ title: 'Error', description: error.error, color: 'error' })
            return false
        }

        if (data) {
            setPublicServers(data)
            return true
        }
        return false
    }

    return {
        servers,
        publicServers,
        addServer,
        setServers,
        removeServer,
        setPublicServers,
        loadServers,
        loadPublicServers,
    }
})
