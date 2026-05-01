<script setup lang="ts">
import { en } from '@nuxt/ui/locale'
</script>

<template>
    <UApp :locale="en">
        <div class="fixed inset-0 p-4 bg-neutral-50 dark:bg-neutral-950">
            <router-view />
        </div>
    </UApp>
</template>
