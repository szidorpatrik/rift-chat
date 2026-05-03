import { apiFetch } from '@utils/api'
import type { MemberResponse } from '@types'

export const memberService = {
    async getServerMembers(serverId: string) {
        return apiFetch<MemberResponse[]>(`/servers/${serverId}/members`)
    }
}
