import { http } from "./http";
import type { Address, HistoryRecord, Profile, Session } from "@/types";

const deviceId = () => localStorage.getItem("sf_device_id") || createDeviceId();

function createDeviceId() {
  const id = `web-${crypto.randomUUID?.() || Date.now()}`;
  localStorage.setItem("sf_device_id", id);
  return id;
}

export const userApi = {
  sendCode(target: string, targetType: "phone" | "email", scene = "register") {
    return http.post("/api/v1/auth/verification-codes", {
      target,
      target_type: targetType,
      scene,
    });
  },
  async register(payload: { phone?: string; email?: string; password: string; verifyCode: string }) {
    return normalizeSession(await http.post<unknown, any>("/api/v1/auth/register", {
      phone: payload.phone || "",
      email: payload.email || "",
      verify_code: payload.verifyCode,
      password: payload.password,
      device_id: deviceId(),
    }));
  },
  async login(account: string, password: string) {
    return normalizeSession(await http.post<unknown, any>("/api/v1/auth/login", {
      account,
      password,
      device_id: deviceId(),
    }));
  },
  logout(session: Session) {
    return http.post("/api/v1/auth/logout", {
      user_id: session.userId,
      access_token: session.accessToken,
      refresh_token: session.refreshToken,
      device_id: deviceId(),
      logout_all: false,
    });
  },
  getProfile(userId: number) {
    return http.get<unknown, Profile>("/api/user/profile", { params: { user_id: userId } });
  },
  updateNickname(userId: number, nickName: string) {
    return http.put<unknown, Profile>("/api/user/profile/nickname", {
      user_id: userId,
      nick_name: nickName,
    });
  },
  listAddresses(userId: number, addrType = 1) {
    return http.get<unknown, { addresses: Address[]; next_page_token: string }>("/v1/user/address/list", {
      params: { user_id: userId, addr_type: addrType, page_size: 20 },
    });
  },
  createAddress(address: Address) {
    return http.post<unknown, Address>("/v1/user/address/add", address);
  },
  updateAddress(address: Address) {
    return http.put<unknown, Address>(`/v1/user/address/edit/${address.id}`, address);
  },
  deleteAddress(userId: number, id: number) {
    return http.delete(`/v1/user/address/${id}`, { params: { user_id: userId } });
  },
  parseAddress(content: string) {
    return http.post<unknown, { address: Address }>("/v1/user/address/ai/parse", { content });
  },
  getHistoryCount(userId: number) {
    return http.get<unknown, { query_count_today: number; max_query_per_day: number }>("/api/order/history/count", {
      params: { user_id: userId },
    });
  },
  listHistory(userId: number) {
    return http.get<unknown, { records: HistoryRecord[]; total_count: number }>("/api/order/history/list", {
      params: { user_id: userId, page_size: 20 },
    });
  },
};

function normalizeSession(raw: any): Session {
  return {
    userId: Number(raw.userId ?? raw.user_id ?? 0),
    accessToken: raw.accessToken ?? raw.access_token ?? "",
    refreshToken: raw.refreshToken ?? raw.refresh_token ?? "",
    role: raw.role ?? "user",
    accountStatus: raw.accountStatus ?? raw.account_status ?? "active",
    expiresIn: Number(raw.expiresIn ?? raw.expires_in ?? 0),
  };
}
