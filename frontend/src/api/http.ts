import axios from "axios";
import { useAuthStore } from "@/stores/auth";

export const http = axios.create({
  baseURL: "",
  timeout: 12000,
});

http.interceptors.request.use((config) => {
  const auth = useAuthStore();
  if (auth.session?.accessToken) {
    config.headers.Authorization = `Bearer ${auth.session.accessToken}`;
  }
  return config;
});

http.interceptors.response.use(
  (response) => {
    const data = response.data;
    if (data && typeof data === "object" && "code" in data && "data" in data) {
      return data.data;
    }
    return data;
  },
  (error) => {
    const message = friendlyApiMessage(error);
    return Promise.reject(new Error(message));
  },
);

function friendlyApiMessage(error: any) {
  if (!error.response) {
    return "无法连接后端服务，请确认服务已启动后再试。";
  }

  const data = error.response?.data || {};
  const rawMessage = String(data.message || error.message || "");
  const reason = String(data.reason || data.error || "").toUpperCase();
  const status = Number(error.response?.status || data.code || 0);
  const signal = `${reason} ${rawMessage}`.toLowerCase();

  if (signal.includes("invalid auth argument") || reason.includes("AUTH_INVALID_ARGUMENT")) {
    return "账号、密码或验证码不符合要求：注册密码需 8-32 位，并且要先获取正确的 6 位验证码。";
  }
  if (reason.includes("AUTH_RESOURCE_NOT_FOUND") || status === 404) {
    return "该手机号还没有注册，请先切换到注册并获取验证码。";
  }
  if (signal.includes("auth unauthenticated") || reason.includes("AUTH_UNAUTHENTICATED") || status === 401) {
    return "账号或密码错误；如果连续输错多次，账号可能会被临时锁定。";
  }
  if (status === 429 || signal.includes("verification")) {
    return "验证码已发送，请 60 秒后再获取。";
  }
  if (signal.includes("duplicated") || reason.includes("AUTH_DUPLICATE_REQUEST")) {
    return "该手机号或邮箱已经注册，请直接登录或换一个账号。";
  }
  if (status === 409) {
    return "当前操作与已有数据冲突，请检查后再试。";
  }
  if (signal.includes("account disabled") || reason.includes("AUTH_ACCOUNT_DISABLED") || status === 403) {
    return "当前账号已被禁用，请联系管理员处理。";
  }
  if (status >= 500) {
    return "后端服务暂时异常，请稍后再试。";
  }

  return rawMessage || "请求失败，请检查填写内容后再试。";
}
