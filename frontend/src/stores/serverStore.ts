import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { ServerResponse } from '@types'
import { apiFetch } from '@utils/api'

export const useServerStore = defineStore('server', () => {
    const toast = useToast()
    const servers = ref<ServerResponse[]>([])
    const publicServers = ref<ServerResponse[]>([])
    const currentServer = ref<ServerResponse | null>(null)

    const sessionStorageKey = 'lastServerId'

    const addServer = (server: ServerResponse) => {
        servers.value.push(server)
        setCurrentServer(server.id)
    }

    const removeServer = (serverId: string) => {
        servers.value = servers.value.filter(s => s.id !== serverId)
        if (currentServer.value?.id === serverId) {
            setCurrentServer(servers.value.length > 0 ? servers.value[0].id : null)
        }
    }

    const setServers = (serverArray: ServerResponse[]) => {
        servers.value = serverArray

        const lastServerId = sessionStorage.getItem(sessionStorageKey)
        if (lastServerId){
            setCurrentServer(lastServerId)
            return
        }

        if (servers.value.length > 0) {
            setCurrentServer(servers.value[0].id)
        }
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
            servers.value.push(data)
            setCurrentServer(data.id)
            return true
        }
        return false
    }

    const setCurrentServer = (serverId: string | null) => {
        if (!serverId) {
            currentServer.value = null
            sessionStorage.removeItem(sessionStorageKey)
            return
        }

        currentServer.value = servers.value.find(s => s.id === serverId) || null
        if (currentServer.value) {
            sessionStorage.setItem(sessionStorageKey, currentServer.value.id)
        } else {
            sessionStorage.removeItem(sessionStorageKey)
        }
    }

    return {
        servers,
        publicServers,
        currentServer,
        addServer,
        setServers,
        removeServer,
        setPublicServers,
        loadServers,
        loadPublicServers,
        addPublicServer,
        joinPublicServer,
        setCurrentServer,
    }
})
