export type Role = "user" | "courier" | "merchant" | "admin";

export interface Session {
  userId: number;
  accessToken: string;
  refreshToken: string;
  role: Role;
  accountStatus: string;
  expiresIn: number;
}

export interface Profile {
  user_id: number;
  avatar_url: string;
  nick_name: string;
  phone: string;
  real_name_auth: number;
  account_status: number;
  is_enterprise: boolean;
  query_count_today: number;
}

export interface Address {
  id?: number;
  user_id: number;
  addr_type: number;
  receiver_name: string;
  receiver_phone: string;
  province: string;
  city: string;
  district: string;
  detail_addr: string;
  is_default?: boolean;
}

export interface HistoryRecord {
  id: number;
  order_no: string;
  express_no: string;
  sender_name: string;
  receiver_name: string;
  status: string;
  create_time?: string;
}
