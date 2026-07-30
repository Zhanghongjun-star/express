import { defineStore } from "pinia";
import { computed, ref } from "vue";
import type { Role, Session } from "@/types";
import { userApi } from "@/api/user";

const storageKey = "sf_session";

function normalizeSession(raw: any): Session {
  return {
    userId: Number(raw.userId ?? raw.user_id ?? 0),
    accessToken: raw.accessToken ?? raw.access_token ?? "",
    refreshToken: raw.refreshToken ?? raw.refresh_token ?? "",
    role: normalizeRole(raw.role),
    accountStatus: raw.accountStatus ?? raw.account_status ?? "active",
    expiresIn: Number(raw.expiresIn ?? raw.expires_in ?? 0),
  };
}

function normalizeRole(role: string): Role {
  const value = String(role || "user").toLowerCase();
  if (["courier", "merchant", "admin"].includes(value)) return value as Role;
  if (value.includes("快递")) return "courier";
  if (value.includes("商家") || value.includes("驿站")) return "merchant";
  if (value.includes("运营")) return "admin";
  return "user";
}

export const roleHome: Record<Role, string> = {
  user: "/user",
  courier: "/courier",
  merchant: "/merchant",
  admin: "/admin",
};

export const useAuthStore = defineStore("auth", () => {
  const saved = localStorage.getItem(storageKey);
  const session = ref<Session | null>(saved ? normalizeSession(JSON.parse(saved)) : null);
  const isLoggedIn = computed(() => Boolean(session.value?.accessToken));

  function setSession(next: Session) {
    session.value = normalizeSession(next);
    localStorage.setItem(storageKey, JSON.stringify(session.value));
  }

  async function login(account: string, password: string) {
    setSession(await userApi.login(account, password));
  }

  async function register(payload: { phone?: string; email?: string; password: string; verifyCode: string }) {
    setSession(await userApi.register(payload));
  }

  async function logout() {
    if (session.value) {
      await userApi.logout(session.value).catch(() => undefined);
    }
    session.value = null;
    localStorage.removeItem(storageKey);
  }

  function switchRole(role: Role) {
    const fallback: Session = {
      userId: 1,
      accessToken: "local-preview-token",
      refreshToken: "local-preview-refresh",
      role,
      accountStatus: "preview",
      expiresIn: 7200,
    };
    setSession({ ...(session.value || fallback), role });
  }

  return { session, isLoggedIn, setSession, login, register, logout, switchRole };
});
