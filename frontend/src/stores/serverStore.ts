import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { ServerResponse } from '@types'

export const useServerStore = defineStore('server', () => {
    const servers = ref<ServerResponse[]>([])

    const addServer = (server: ServerResponse) => {
        servers.value.push(server)
    }

    const removeServer = (serverId: string) => {
        servers.value = servers.value.filter(s => s.id !== serverId)
    }

    const setServers = (serverArray: ServerResponse[]) => {
        servers.value = serverArray
    }

    return {
        servers,
        addServer,
        setServers,
        removeServer
    }
})
