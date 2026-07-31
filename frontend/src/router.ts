import { createRouter, createWebHashHistory } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import LoginView from "@/views/LoginView.vue";
import UserPortal from "@/views/UserPortal.vue";

export const router = createRouter({
  history: createWebHashHistory(),
  scrollBehavior: () => ({ top: 0 }),
  routes: [
    { path: "/", redirect: "/user/login" },
    { path: "/login", redirect: "/user/login" },
    { path: "/user/login", component: LoginView, meta: { public: true, loginRole: "user" } },
    { path: "/user", component: UserPortal, meta: { role: "user" } },
    { path: "/:pathMatch(.*)*", redirect: "/user/login" },
  ],
});

router.beforeEach((to) => {
  const auth = useAuthStore();
  if (to.meta.public) return true;
  if (!auth.isLoggedIn) return { path: "/user/login", query: { redirect: to.fullPath } };
  return true;
});
