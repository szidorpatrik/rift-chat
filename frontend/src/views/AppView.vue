<script setup lang="ts">
import type { LoginResponse, Server, Channel, Message } from '@types'
import { onMounted, ref } from 'vue'
import { apiFetch } from '@api'
import { useAuthStore } from '@stores/authStore'
import router from '@router'
import ServerSidebar from '@components/ServerSidebar.vue'
import ChannelSidebar from '@components/ChannelSidebar.vue'
import ChatArea from '@components/ChatArea.vue'
import UserProfile from '@components/UserProfile.vue'


const authStore = useAuthStore()
const toast = useToast()

const serverList = ref<Server[]>([])
const channelList = ref<Channel[]>([])
const messageList = ref<Message[]>([])

onMounted(async () => {
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
})
</script>

<template>
    <div class="flex flex-1 h-full gap-4">

        <div class="flex flex-col gap-4 h-full">

            <div class="flex flex-1 gap-4 min-h-0">
                <ServerSidebar :servers="serverList" />
                <ChannelSidebar :channels="channelList" />
            </div>

            <UserProfile :username="authStore.currentUser?.username" />
        </div>

        <ChatArea :messages="messageList" />

    </div>
</template>
