<script setup lang="ts">
import { ref } from 'vue'
import type { Message } from '@types'

defineProps<{
    messages: Message[]
}>()

const input = ref('')
</script>

<template>
    <div class="flex-1 rounded-2xl border border-neutral-200 dark:border-neutral-800 bg-white dark:bg-neutral-900 shadow-sm flex flex-col overflow-hidden">

        <div class="flex-1 p-4 overflow-y-auto flex flex-col gap-4">
            <template v-if="messages.length">
                <div v-for="msg in messages" :key="msg.id" class="flex flex-col gap-1">
                    <span class="font-bold text-sm">{{ msg.author }}</span>
                    <p class="text-neutral-600 dark:text-neutral-300">{{ msg.content }}</p>
                </div>
            </template>
            <template v-else>
                <USkeleton v-for="_ in 60" class="h-8 w-full shrink-0" />
            </template>
        </div>

        <UTextarea v-model="input" 
            autoresize 
            placeholder="Message @placeholder..." 
            :autofocus="true" 
            :rows="2" 
            variant="none" 
            size="xl" 
            class="bg-neutral-50 dark:bg-neutral-950" 
            :ui="{ base: 'min-h-[5rem] max-h-128 focus:ring-0' }"/>
    </div>
</template>
