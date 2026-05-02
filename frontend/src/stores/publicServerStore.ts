import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { ServerResponse } from '@types'
import { serverService } from '@services/serverService'

export const usePublicServerStore = defineStore('publicServer', () => {
    const toast = useToast()
    
    // --- State ---
    const servers = ref<ServerResponse[]>([])
    const areServersLoaded = ref(false)

    // --- Actions ---
    const loadServers = async () => {
        if (servers.value.length > 0){
            areServersLoaded.value = true
            return true
        }

        const { data, error } = await serverService.getPublicServers()
        if (error) {
            toast.add({ title: 'Error', description: error.error, color: 'error' })
            return false
        }
        if (data) {
            servers.value = data
            areServersLoaded.value = true
            return true
        }
        return false
    }

    const addServer = (server: ServerResponse) => {
        servers.value.push(server)
    }

    return {
        servers,
        areServersLoaded,
        loadServers,
        addServer
    }
})