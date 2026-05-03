import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { MemberResponse } from '@types'
import { memberService } from '@services/memberService'
import { useServerStore } from '@stores/serverStore'
import type { UUID } from 'crypto'

export const useMemberStore = defineStore('member', () => {
    const toast = useToast()

    // --- State ---
    const loadingMembers = ref<Set<string>>(new Set())

    // --- Actions ---
    const loadMembers = async (serverId: string) => {
        if (loadingMembers.value.has(serverId)) return false

        loadingMembers.value.add(serverId)
        const serverStore = useServerStore()
        const server = serverStore.getServer(serverId)
        
        if (!server) {
            loadingMembers.value.delete(serverId)
            return false
        }

        if (server.membersLoaded) {
            loadingMembers.value.delete(serverId)
            return true
        }

        try {
            const { data, error } = await memberService.getServerMembers(serverId)
            if (error) {
                toast.add({ title: 'Error', description: error.error, color: 'error' })
                return false
            }
            if (data) {
                server.members = data
                server.membersLoaded = true
                return true
            }
        } finally {
            loadingMembers.value.delete(serverId)
        }
    }

    const addMember = (serverId: UUID, member: MemberResponse) => {
        const serverStore = useServerStore()
        const server = serverStore.getServer(serverId)
        if (server) {
            server.members.push(member)
        }
    }

    return {
        loadMembers,
        addMember,
    }
})