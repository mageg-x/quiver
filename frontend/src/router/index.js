import { createRouter, createWebHistory } from "vue-router";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/login",
      name: "login",
      component: () => import("@/view/LoginView.vue"),
    },
    {
      path: "/",
      name: "home",
      component: () => import("@/view/HomeView.vue"),
      children: [
        {
          path: "",
          name: "main",
          component: () => import("@/view/MainView.vue"),
        },
        {
          path: 'apps',
          name: 'apps',
          component: () => import('@/view/AppView.vue')
        },
        {
          path: 'clusters',
          name: 'clusters',
          component: () => import('@/view/ClusterView.vue'),
          props: true
        },
        {
          path: 'namespaces',
          name: 'namespaces',
          component: () => import('@/view/NamespaceView.vue'),
          props: true
        },
        {
          path: 'users',
          name: 'users',
          component: () => import('@/view/UserView.vue'),
          props: true
        },
        {
          path: 'permissions',
          name: 'permissions',
          component: () => import('@/view/PermissionView.vue'),
          props: true
        },
      ],
    },
  ],
});

// 路由守卫

export default router;
