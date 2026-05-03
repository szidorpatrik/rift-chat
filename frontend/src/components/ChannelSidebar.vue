<script setup lang="ts">
import { useServerStore } from '@stores/serverStore';
import type { Channel } from '@types'
import type { ContextMenuItem } from '@nuxt/ui'
import CreateChannelView from '@views/CreateChannelView.vue'
import { ref } from 'vue';
import { useChannelStore } from '@stores/channelStore';

const serverStore = useServerStore()
const channelStore = useChannelStore()
const openCreateChannel = ref(false)

// Server header actions
const serverHeaderActions = ref<ContextMenuItem[][]>([
    [
        {
            label: 'Server Settings',
            icon: 'i-lucide-settings',
        },
        {
            label: 'Invite People',
            icon: 'i-lucide-user-plus',
        }
    ]
])

// Sidebar actions (Creating channels)
const channelCreateActions = ref<ContextMenuItem[][]>([
    [
        {
            label: 'Create Channel',
            icon: 'i-lucide-plus',
            onClick: () => {
                openCreateChannel.value = !openCreateChannel.value
            },
        },
        {
            label: 'Create Category',
            icon: 'i-lucide-folder-plus',
            disabled: true,
        },
    ],
])

const getChannelEditActions = (channel: Channel): ContextMenuItem[][] => [
    [
        {
            label: 'Edit Channel',
            icon: 'i-lucide-pencil',
            onClick: () => {
                console.log('Editing channel:', channel.name)
                // Add your channel edit logic here
            },
        },
        {
            label: 'Delete Channel',
            icon: 'i-lucide-trash-2',
            color: 'error',
            onClick: () => {
                channelStore.deleteChannel(serverStore.currentServerId!, channel.id)
            }
        },
    ],
]

const getIcon = (type: 0 | 1 | 2) => {
    switch (type) {
        case 0: return 'i-lucide-hash'
        case 1: return 'i-lucide-volume-2'
        case 2: return 'i-lucide-folder'
    }
}
</script>

<template>
    <div class="select-none flex flex-1 flex-col max-w-56 justify-between overflow-hidden rounded-lg">
        <template v-if="serverStore.currentServer">
            <!-- Server Context Menu -->
            <UContextMenu :items="serverHeaderActions">
                <div class="flex mb-1">
                    <UButton class=" max-w-fit truncate text-md font-semibold text-neutral-400 tracking-wide"
                             color="neutral"
                             variant="ghost"
                             size="xs"
                             trailing-icon="i-lucide-chevron-down"
                             :label="serverStore.currentServer?.name" />

                    <UButton class=" max-w-fit ml-auto text-neutral-400"
                             color="neutral"
                             variant="ghost"
                             size="lg"
                             icon="i-lucide-user-plus" />
                </div>
            </UContextMenu>
            <div v-if="serverStore.currentServer"
                 class="overflow-y-auto shrink-0 flex flex-col flex-1 rounded-lg border border-neutral-200 dark:border-neutral-800 bg-white dark:bg-neutral-900 shadow-sm [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">
                <!-- Empty Space Context Menu (Create Channels) -->
                <UContextMenu :items="channelCreateActions">
                    <UScrollArea class="h-full w-full p-2 mt-0.5 hover-scrollbar">

                        <template v-if="serverStore.currentServer?.channelsLoaded">

                            <template v-if="!serverStore.currentServer?.channels.length">
                                <div class="mt-2 select-none flex flex-col items-center gap-2 text-sm text-center text-neutral-500">
                                    <p>No channels yet.</p>
                                    <div class="flex items-center gap-2">
                                        <span>Right-click</span>
                                        <UIcon class="inline mb-0.5" name="i-lucide-arrow-right" />
                                        <span>Create Channel</span>
                                    </div>
                                </div>
                            </template>

                            <template v-else>
                                <UContextMenu v-for="channel in serverStore.currentServer?.channels"
                                              :key="channel.id"
                                              :items="getChannelEditActions(channel)">

                                    <UButton class="p-2 w-full truncate text-sm tracking-wide"
                                             color="neutral"
                                             variant="ghost"
                                             size="md"
                                             :icon="getIcon(channel.type)"
                                             :label="channel.name"
                                             @click="channelStore.setCurrentChannel(serverStore.currentServerId, channel.id)" />
                                </UContextMenu>
                            </template>

                        </template>
                    </UScrollArea>
                </UContextMenu>
                <UModal v-model:open="openCreateChannel" title="Create Channel">
                    <template #body>
                        <CreateChannelView v-model:open="openCreateChannel" />
                    </template>
                </UModal>
            </div>
        </template>

        <template v-else>
            <div class="flex flex-1 gap-1 items-start pt-4 px-2 justify-center text-sm text-neutral-500 border border-neutral-200 dark:border-neutral-800 bg-white dark:bg-neutral-900 shadow-sm">
                <UIcon name="i-lucide-arrow-left" class="size-5 text-neutral-400" /> 
                <span>Create or join a server</span>
            </div>
        </template>
    </div>
</template>