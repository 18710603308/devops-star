import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";
import { useUserStore } from "@/stores/user";

const routes: RouteRecordRaw[] = [
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/login/index.vue"),
    meta: { title: "登录", requiresAuth: false },
  },
  {
    path: "/",
    redirect: "/dashboard",
  },
  {
    path: "/dashboard",
    name: "Dashboard",
    component: () => import("@/views/dashboard/index.vue"),
    meta: { title: "仪表盘", requiresAuth: true },
  },
  {
    path: "/projects",
    name: "Projects",
    component: () => import("@/views/projects/index.vue"),
    meta: { title: "项目管理", requiresAuth: true },
  },
  {
    path: "/projects/:id",
    name: "ProjectDetail",
    component: () => import("@/views/projects/detail.vue"),
    meta: { title: "项目详情", requiresAuth: true },
  },
  {
    path: "/pipeline",
    name: "Pipeline",
    component: () => import("@/views/pipeline/index.vue"),
    meta: { title: "流水线", requiresAuth: true },
  },
  {
    path: "/pipeline/editor/:id?",
    name: "PipelineEditor",
    component: () => import("@/views/pipeline/editor.vue"),
    meta: { title: "流水线编排", requiresAuth: true },
  },
  {
    path: "/registry",
    name: "Registry",
    component: () => import("@/views/registry/index.vue"),
    meta: { title: "制品管理", requiresAuth: true },
  },
  {
    path: "/deploy",
    name: "Deploy",
    component: () => import("@/views/deploy/index.vue"),
    meta: { title: "部署管理", requiresAuth: true },
  },
  {
    path: "/monitor",
    name: "Monitor",
    component: () => import("@/views/monitor/index.vue"),
    meta: { title: "监控大屏", requiresAuth: true },
  },
  {
    path: "/settings",
    name: "Settings",
    component: () => import("@/views/settings/index.vue"),
    meta: { title: "系统设置", requiresAuth: true },
  },
  {
    path: "/:pathMatch(.*)*",
    redirect: "/dashboard",
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// 路由守卫
router.beforeEach((to, _from, next) => {
  const userStore = useUserStore();
  const requiresAuth = to.meta.requiresAuth !== false;

  if (requiresAuth && !userStore.token) {
    next("/login");
  } else if (to.path === "/login" && userStore.token) {
    next("/dashboard");
  } else {
    next();
  }
});

export default router;
