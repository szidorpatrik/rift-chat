<script setup lang="ts">
import router from '@router';
import { useAuthStore } from '@stores/authStore';
import type { DropdownMenuItem } from '@nuxt/ui';
import { computed, ref } from 'vue'
import { statusColorMap } from '@types';

const toast = useToast()
const authStore = useAuthStore()

const statusMessage = computed(() => {
    switch (authStore.user?.status) {
        case 'online':
            return 'Online';
        case 'idle':
            return 'Idle';
        case 'dnd':
            return 'Do Not Disturb';
        default:
            return 'Offline';
    }
})

const items = ref<DropdownMenuItem[][]>([
    [
        {
            label: 'Profile',
            icon: 'mdi-user'
        },
    ],
    [
        {
            label: 'Logout',
            icon: 'i-lucide-log-out',
            color: 'error',
            onClick: () => {
                onLogout();
            }
        }
    ]
])

const onLogout = async () => {
    const { error } = await authStore.logout()

    if (error) {
        toast.add({
            title: 'Logout failed',
            description: error.error || 'An unexpected error occurred. Please try again later.',
            color: 'error'
        })
        return
    }

    toast.add({
        title: 'Logged out',
        description: 'You have been logged out successfully.',
        color: 'success'
    })
    router.push({ name: 'login' })
}

</script>

<template>
    <div id="btn-profile"
         class="shrink-0 flex justify-between items-center min-h-16 rounded-lg border border-neutral-200 dark:border-neutral-800 bg-white dark:bg-neutral-900 shadow-sm p-2">
        <div class="flex-1 items-center min-w-0 mr-1">
            <UDropdownMenu :items="items"
                           :content="{ align: 'start', side: 'top', sideOffset: 16 }"
                           size="lg">
                <UButton variant="ghost" color="neutral" class="w-full rounded-lg text-left">
                    <UUser class="w-full overflow-x-hidden min-w-0 p-0"
                           :description="authStore.currentUser?.description ? authStore.currentUser.description : statusMessage"
                           :avatar="{ src: authStore.currentUser?.avatar_url, loading: 'lazy' }"
                           :chip="{ color: statusColorMap[authStore.currentUser?.status || 'offline'], position: 'bottom-right' }"
                           :ui="{ description: 'truncate', name: 'truncate' }">
                        <template #name>
                            <div class="flex flex-1 min-w-0 items-center gap-1">
                                <span class="truncate block min-w-0 max-w-26">
                                    {{ authStore.currentUser?.username }}
                                </span>
                            </div>
                        </template>
                    </UUser>
                </UButton>
            </UDropdownMenu>
        </div>

        <div class="select-none shrink-0">
            <div class="flex gap-1 h-full">
                <UButton variant="ghost" color="neutral" size="xs" class="p-1.5 rounded-lg">
                    <UIcon name="i-lucide-mic" size="20" />
                </UButton>
                <UButton variant="ghost" color="neutral" size="xs" class="p-1.5 rounded-lg">
                    <UIcon name="i-lucide-headset" size="20" />
                </UButton>
                <UButton variant="ghost" color="neutral" size="xs" class="p-1.5 rounded-lg">
                    <UIcon name="i-mdi-cog-outline" size="20" />
                </UButton>
            </div>
        </div>
    </div>
</template>
