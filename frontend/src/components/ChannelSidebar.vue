<script setup lang="ts">
import { useServerStore } from '@stores/serverStore';
import type { Channel } from '@types'
import type { ContextMenuItem } from '@nuxt/ui'
import CreateChannelView from '@views/CreateChannelView.vue'
import { ref, computed } from 'vue';

const props = defineProps<{
    channels: Channel[] | []
}>()


const sortedChannels = computed(() => {
    return [...props.channels].sort((a, b) => a.position - b.position)
})

const serverStore = useServerStore()
const openCreateChannel = ref(false)

// 1. New actions for the Server header
const serverActions = ref<ContextMenuItem[][]>([
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

// 2. Actions for empty space in the sidebar (Creating channels)
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

// 3. Changed from a static ref to a function to capture the specific channel clicked
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
                console.log('Deleting channel:', channel.name)
                // Add your channel delete logic here
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
    <div class="overflow-y-auto w-64 shrink-0 flex flex-col rounded-lg gap-4 border border-neutral-200 dark:border-neutral-800 bg-white dark:bg-neutral-900 shadow-sm [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">

        <!-- Server Context Menu -->
        <UContextMenu :items="serverActions">
            <div class="flex flex-0 mx-4 mt-4">
                <UButton class="p-2 max-w-fit truncate text-md font-semibold text-neutral-400 tracking-wide"
                         color="neutral"
                         variant="ghost"
                         size="md"
                         trailing-icon="i-lucide-chevron-down"
                         :label="serverStore.currentServer?.name || 'Select or Join a server'" />

                <UButton class="p-2 max-w-fit ml-auto text-neutral-400"
                         color="neutral"
                         variant="ghost"
                         size="md"
                         icon="i-lucide-user-plus" />
            </div>
        </UContextMenu>

        <!-- Empty Space Context Menu (Create Channels) -->
        <UContextMenu :items="channelCreateActions">
            <UScrollArea :ui="{ viewport: 'gap-2' }" class="h-full w-full px-4 pb-4 mt-0.5 hover-scrollbar">

                <template v-if="serverStore.areChannelsLoaded(serverStore.currentServerId!)">

                    <template v-if="!sortedChannels.length">
                        <div class="select-none flex flex-col items-center gap-2 text-sm text-center text-neutral-500">
                            <p>No channels yet.</p>
                            <p>Right-click <UIcon class="inline mb-0.5" name="i-lucide-arrow-right" /> Create Channel</p>
                        </div>
                    </template>

                    <template v-else>
                        <UContextMenu v-for="channel in sortedChannels"
                                      :key="channel.id"
                                      :items="getChannelEditActions(channel)">
                            <div class="h-10 w-full shrink-0 flex items-center px-2 rounded-lg hover:bg-neutral-100 dark:hover:bg-neutral-800 cursor-pointer">
                                <UIcon :name="getIcon(channel.type)" class="mr-2" />
                                {{ channel.name }}
                            </div>
                        </UContextMenu>
                    </template>

                </template>

                <template v-else>
                    <USkeleton v-for="_ in 4" class="h-8 w-full shrink-0" />
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