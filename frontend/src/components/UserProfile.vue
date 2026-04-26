<script setup lang="ts">
defineProps<{
    username?: string
}>()
</script>

<template>
    <div class="shrink-0 flex min-h-16 rounded-2xl border border-neutral-200 dark:border-neutral-800 bg-white dark:bg-neutral-900 shadow-sm p-4 items-center">
        
        <template v-if="username">
            <div class="flex items-center gap-3 w-full">
                <UAvatar :alt="username" size="sm" />
                <span class="font-medium truncate">{{ username }}</span>
            </div>
        </template>

        <template v-else>
            <div class="flex items-center gap-3 w-full">
                <USkeleton class="size-8 rounded-full" />
                <USkeleton class="h-4 w-24" />
            </div>
        </template>
        
    </div>
</template>
