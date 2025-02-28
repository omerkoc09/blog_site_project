import type { App } from 'vue'


import { setupLayouts } from 'virtual:generated-layouts'
import type { RouteRecordRaw } from 'vue-router/auto'

import { createRouter, createWebHistory } from 'vue-router/auto'
import JwtService from '@/services/JwtService'
import { useUserStore } from '@/store/user'

function recursiveLayouts(route: RouteRecordRaw): RouteRecordRaw {
  if (route.children) {
    for (let i = 0; i < route.children.length; i++)
      route.children[i] = recursiveLayouts(route.children[i])

    return route
  }

  return setupLayouts([route])[0]
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  scrollBehavior(to) {
    if (to.hash)
      return { el: to.hash, behavior: 'smooth', top: 60 }

    return { top: 0 }
  },
  extendRoutes: pages => [
    ...[...pages].map(route => recursiveLayouts(route)),
  ],
  
})


router.beforeEach((to, from, next) => {

  // First, define public routes array
  const publicRoutes = [
    '/auth/login',
    '/auth/register',
    '/auth/forgot-password',
    '/reset-password',

    

  ]

  // Check if current route path is public
  if (publicRoutes.includes(to.path)) {
    return next()
  }

  const isLoggedIn = !!JwtService.getAccessToken()

  if (!isLoggedIn) {
    if (to.meta.redirectIfLoggedIn)
      return next()
    return next({ name: 'auth-login', query: { to: to.name !== 'index' ? to.fullPath : undefined } })
  }

  const roleToNumber = {
    'admin': 10,
    'normal': 1
  }

  if (to.meta.redirectIfLoggedIn)
    return next('/')

  if (typeof to.meta.role === 'string' && !useUserStore().hasRole(roleToNumber[to.meta.role])) {
    return next({ name: 'not-authorized' })
  }

  return next()
})

export { router }

export default function (app: App) {
  app.use(router)
}
