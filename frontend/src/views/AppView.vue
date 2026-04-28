<script setup lang="ts">
import { type LoginResponse, type Server, type Channel, type Message } from '@types'
import type { ServerResponse } from '@types'
import { onMounted, ref } from 'vue'
import { apiFetch } from '@api'
import { useAuthStore } from '@stores/authStore'
import router from '@router'
import ServerSidebar from '@components/ServerSidebar.vue'
import ChannelSidebar from '@components/ChannelSidebar.vue'
import ChatArea from '@components/ChatArea.vue'
import UserProfile from '@components/UserProfile.vue'
import { useServerStore } from '@stores/serverStore'


const authStore = useAuthStore()
const serverStore = useServerStore()
const toast = useToast()

const channelList = ref<Channel[]>([])
const messageList = ref<Message[]>([])

onMounted(async () => {
    try {
        await handleSessionCheck()
        await loadServers()
    } catch (error) {
        console.error(error)
    }
})

const handleSessionCheck = async () => {
    const { data, error } = await apiFetch<LoginResponse>('/session')
    if (data) {
        authStore.setUser(data)
        return
    }

    if (error?.status_code === 500) {
        toast.add({
            title: 'Server error',
            description: error.error || 'An unexpected error occurred. Please try again later.',
            color: 'error'
        })
        return
    }

    authStore.clearUser()
    toast.add({
        title: 'Session expired',
        description: 'Please log in again.',
        color: 'warning'
    })
    router.push({ name: 'login' })
}

const loadServers = async () => {
    const { data, error } = await apiFetch<ServerResponse[]>('/servers')
    if (error) {
        toast.add({ title: 'Error', description: error.error, color: 'error' })
        return
    }

    if (data) {
        serverStore.setServers(data)
    }
}
</script>

<template>
    <div class="flex flex-1 h-full gap-4">

        <div class="flex flex-col gap-4 h-full">

            <div class="flex flex-1 gap-4 min-h-0">
                <ServerSidebar />
                <ChannelSidebar :channels="channelList" />
            </div>

            <UserProfile :username="authStore.currentUser?.username"
                         :description="authStore.currentUser?.description"
                         :avatar_url="authStore.currentUser?.avatar_url"
                         :status="authStore.currentUser?.status" />
        </div>

        <ChatArea :messages="messageList" />

    </div>
</template>
