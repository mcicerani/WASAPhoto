import { createRouter, createWebHashHistory } from 'vue-router'
import Login from '../views/Login.vue'
import UserStream from '../views/UserStream.vue'
import UserProfile from '../views/UserProfile.vue'
import EditProfile from '../views/EditProfile.vue'
import SearchUser from '../views/SearchUser.vue'

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', redirect: '/session' }, // Reindirizza alla pagina di login se l'utente non è autenticato
    { path: '/session', component: Login }, // Pagina di login
    { path: '/users/:userId/profile', component: UserProfile, meta: { requiresAuth: true } }, // Pagina del profilo utente
    { path: '/users/:userId/profile/edit', component: EditProfile, meta: { requiresAuth: true } }, // Pagina di modifica del profilo utente
    { path: '/users/:userId/stream', component: UserStream, meta: { requiresAuth: true } }, // Pagina dello stream foto del user loggato
    { path: '/users', component: SearchUser, meta: { requiresAuth: true } }, // Pagina di ricerca utenti
  ]
})

export default router
