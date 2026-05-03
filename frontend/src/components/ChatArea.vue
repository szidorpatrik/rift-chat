<script setup lang="ts">
import { ref } from 'vue'
import { useChannelStore } from '@stores/channelStore';
import MembersSidebar from '@components/MembersSidebar.vue';

const channelStore = useChannelStore()
const input = ref('')

const getIcon = (type: 0 | 1 | 2) => {
    switch (type) {
        case 0: return 'i-lucide-hash'
        case 1: return 'i-lucide-volume-2'
    }
}
</script>

<template>
    <div class="flex flex-1 w-full gap-4">
        <div class="flex flex-1 flex-col justify-between overflow-hidden rounded-lg">
            <!-- Header -->
            <div v-if="channelStore.currentChannel"
                 class="flex gap-1 pt-1.5 pb-2.5 items-center select-none font-bold text-muted">

                <UIcon :name="getIcon(channelStore.currentChannel?.type || 0)" class="size-5" />
                <span>{{ channelStore.currentChannel?.name }}</span>
            </div>

            <MessageList/>

            <!-- Message Input -->
            <UTextarea v-model="input" v-if="channelStore.currentChannel" autoresize :autofocus="true" :rows="2" :ui="{ base: 'min-h-[4rem] mt-4 rounded-lg max-h-128 focus:ring-0 border border-neutral-200 dark:border-neutral-800 bg-neutral-50 dark:bg-neutral-950' }" :placeholder="`Message #${channelStore.currentChannel?.name || '...'}`" variant="ghost" size="xl" />
        </div>
        <!-- Server Members -->
        <div v-if="channelStore.currentChannel" class="max-lg:hidden flex flex-col w-56 select-none">
            <div class="flex flex-1 max-w-56 max-h-9 mb-1 items-center gap-1">
                <UIcon name="i-lucide-users" class="size-5 text-muted" />
                <span class="font-bold text-muted">Members</span>
            </div>
            <div class="flex flex-1 items-center gap-2 p-2 border-b rounded-lg border border-neutral-200 dark:border-neutral-800 bg-white dark:bg-neutral-900 shadow-sm">
                <MembersSidebar />
            </div>
        </div>
    </div>
</template>
