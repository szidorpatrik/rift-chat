<script setup lang="ts">
import { type Channel, type Message } from '@types'
import { onMounted, ref } from 'vue'
import { useAuthStore } from '@stores/authStore'
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
        const isAuthenticated = await authStore.handleSessionCheck()
        if (!isAuthenticated) return

        const serversLoaded = await serverStore.loadServers()
        // await loadChannels()
        // await loadMessages()
        const publicServersLoaded = await serverStore.loadPublicServers()
    } catch (error) {
        console.error('Initialization failed:', error)
    }
})

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
