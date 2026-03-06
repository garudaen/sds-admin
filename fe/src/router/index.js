import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    component: () => import('../views/Layout.vue'),
    children: [
      {
        path: '',
        redirect: '/domains',
      },
      {
        path: 'domains',
        name: 'Domains',
        component: () => import('../views/Domains.vue'),
        meta: { title: '域名管理' },
      },
      {
        path: 'domains/:id/records',
        name: 'DomainRecords',
        component: () => import('../views/DomainRecords.vue'),
        meta: { title: '解析记录' },
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  document.title = to.meta?.title ? `${to.meta.title} - DNS管理` : 'DNS管理'
  next()
})

export default router
