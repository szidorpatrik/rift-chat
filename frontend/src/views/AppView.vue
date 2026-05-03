<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '@stores/authStore'
import ServerSidebar from '@components/ServerSidebar.vue'
import ChannelSidebar from '@components/ChannelSidebar.vue'
import ChatArea from '@components/ChatArea.vue'
import UserProfile from '@components/UserProfile.vue'
import { useServerStore } from '@stores/serverStore'

const authStore = useAuthStore()
const serverStore = useServerStore()

onMounted(async () => {
    try {
        const isAuthenticated = await authStore.handleSessionCheck()
        if (!isAuthenticated) return

        await serverStore.loadServers()
    } catch (error) {
        console.error('Initialization failed:', error)
    }
})
</script>

<template>
    <div class="flex flex-1 h-full gap-4 min-w-0">

        <div class="flex flex-col gap-4 h-full w-full min-w-0 max-w-72">

            <div class="flex flex-1 gap-4 min-h-0 min-w-0">
                <ServerSidebar class="shrink-0" />
                <ChannelSidebar class="shrink-0" />
            </div>

            <UserProfile :username="authStore.currentUser?.username"
                         :description="authStore.currentUser?.description"
                         :avatar_url="authStore.currentUser?.avatar_url"
                         :status="authStore.currentUser?.status" />
        </div>

        <ChatArea class="min-w-0 flex-1" />

    </div>
</template>
