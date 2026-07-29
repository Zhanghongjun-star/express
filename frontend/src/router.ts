import { createRouter, createWebHashHistory } from "vue-router";
import { roleHome, useAuthStore } from "@/stores/auth";
import type { Role } from "@/types";
import LoginView from "@/views/LoginView.vue";
import UserPortal from "@/views/UserPortal.vue";
import CourierPortal from "@/views/CourierPortal.vue";
import MerchantPortal from "@/views/MerchantPortal.vue";
import AdminPortal from "@/views/AdminPortal.vue";

export const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: "/", redirect: "/user/login" },
    { path: "/login", redirect: "/user/login" },
    { path: "/user/login", component: LoginView, meta: { public: true, loginRole: "user" } },
    { path: "/courier/login", component: LoginView, meta: { public: true, loginRole: "courier" } },
    { path: "/merchant/login", component: LoginView, meta: { public: true, loginRole: "merchant" } },
    { path: "/admin/login", component: LoginView, meta: { public: true, loginRole: "admin" } },
    { path: "/user", component: UserPortal, meta: { role: "user" } },
    { path: "/courier", component: CourierPortal, meta: { role: "courier" } },
    { path: "/merchant", component: MerchantPortal, meta: { role: "merchant" } },
    { path: "/admin", component: AdminPortal, meta: { role: "admin" } },
  ],
});

router.beforeEach((to) => {
  const auth = useAuthStore();
  if (to.meta.public) return true;
  const expectedRole = to.meta.role as Role | undefined;
  if (!auth.isLoggedIn) return { path: `/${expectedRole || "user"}/login`, query: { redirect: to.fullPath } };
  if (expectedRole && auth.session?.role !== expectedRole) return roleHome[auth.session?.role || "user"];
  return true;
});
