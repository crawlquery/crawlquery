import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'search',
      component: HomeView
    },
    {
      path: '/crawl',
      name: 'crawl',
      component: () => import('../views/CrawlView.vue')
    },
    {
      path: '/query',
      name: 'query',
      component: () => import('../views/QueryView.vue')
    },
  ]
})

export default router
