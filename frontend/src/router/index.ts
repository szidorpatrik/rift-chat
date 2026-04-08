import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '@views/LoginView.vue'
import AppView from '@views/AppView.vue'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'app',
            component: AppView
        },
        {
            path: '/login',
            name: 'login',
            component: LoginView
        },
    ]
})

export default router
