import { createRouter, createWebHistory } from 'vue-router';

const routes = [
  {
    path: '/login',
    component: () => import('@/components/WebIndex.vue')
  },
];


const router = createRouter({
  history: createWebHistory(),  // history路由模式
  routes
});

export default router;
