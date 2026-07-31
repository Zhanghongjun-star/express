<template>
  <main class="sf-user-app">
    <section class="sf-user-shell">
      <header v-if="active === 'home'" class="sf-user-topbar sf-home-topbar">
        <button class="sf-city-chip sf-home-city" @click="ElMessage.info('城市服务暂未开放')">{{ cityName }} <span>▾</span></button>
        <button class="sf-ai-search" @click="openTracking()">
          <span class="sf-ai-mark">AI</span>
          <span>寄件、查件、找功能</span>
        </button>
        <div class="sf-home-actions">
          <button class="sf-home-action" @click="openTracking()">⌕</button>
          <button class="sf-home-action" @click="active = 'mine'">◉</button>
        </div>
      </header>
      <header v-else class="sf-user-topbar">
        <div>
          <p class="sf-eyebrow">顺丰速运</p>
          <h1>{{ page.title }}</h1>
          <span>{{ page.subtitle }}</span>
        </div>
        <button class="sf-icon-pill" @click="handleShortcut">{{ page.shortcut }}</button>
      </header>

      <section v-if="active === 'home'" class="sf-page sf-home-page">
        <section class="sf-hero-banner">
          <div class="sf-hero-copy">
            <small>充值最高返 6.5%</small>
            <h2>多充多送，全国通用</h2>
            <button class="sf-promo-btn">马上参与</button>
          </div>
          <div class="sf-hero-card-art">
            <span>SF</span>
            <i></i>
          </div>
        </section>

        <section class="sf-tile-row">
          <article class="sf-tile sf-tile-accent">
            <strong>大众会员</strong>
            <span>领专属会员优惠</span>
            <button>去领取</button>
          </article>
          <article class="sf-tile">
            <strong>充值返利</strong>
            <span>多充多返，寄件享优惠</span>
            <button>去充值</button>
          </article>
        </section>

        <section class="sf-primary-actions">
          <button class="sf-primary-action" @click="active = 'ship'">
            <span class="sf-primary-icon">➤</span>
            <strong>寄快递</strong>
            <small>一小时取件</small>
          </button>
          <button class="sf-primary-action" @click="active = 'ship'">
            <span class="sf-primary-icon sf-primary-icon-orange">▰</span>
            <strong>发物流</strong>
            <small>大件 / 零担 / 整车</small>
          </button>
        </section>

        <section class="sf-grid-panel">
          <button v-for="item in homeServices" :key="item.label" class="sf-service-tile" @click="active = item.target">
            <b>{{ item.icon }}</b>
            <strong>{{ item.label }}</strong>
            <small>{{ item.desc }}</small>
          </button>
        </section>

        <section class="sf-coupon-strip">
          <article v-for="coupon in coupons" :key="coupon.title" class="sf-coupon-card">
            <span>{{ coupon.title }}</span>
            <b>{{ coupon.value }}</b>
            <small>{{ coupon.note }}</small>
          </article>
        </section>

        <section class="sf-home-list">
          <div class="sf-section-head">
            <strong>最近运单</strong>
            <a @click="active = 'track'">查看全部</a>
          </div>
          <article v-for="item in visibleOrders" :key="item.order_no" class="sf-track-card compact">
            <div class="sf-track-head">
              <strong>{{ item.order_no }}</strong>
              <span>{{ item.status }}</span>
            </div>
            <p>{{ item.sender_name }} → {{ item.receiver_name }}</p>
            <small>{{ item.create_time || "刚刚更新" }}</small>
          </article>
        </section>
      </section>

      <section v-else-if="active === 'ship'" class="sf-page sf-ship-page">
        <section class="sf-form-hero">
          <div>
            <p class="sf-eyebrow">寄件服务</p>
            <h2>填写信息，安排取件</h2>
            <span>提交后会写入订单数据库，生成待揽收运单。</span>
          </div>
          <button class="sf-back-button" @click="active = 'home'">返回首页</button>
        </section>

        <section class="sf-form-card">
          <div class="sf-form-section">
            <div class="sf-section-head"><strong>寄件人信息</strong><span>上门取件</span></div>
            <div class="sf-form-grid">
              <el-input v-model.trim="shipmentForm.sender_name" placeholder="姓名" />
              <el-input v-model.trim="shipmentForm.sender_phone" placeholder="手机号" />
              <el-input v-model.trim="shipmentForm.sender_province" placeholder="省" />
              <el-input v-model.trim="shipmentForm.sender_city" placeholder="市" />
              <el-input v-model.trim="shipmentForm.sender_district" placeholder="区" />
              <el-input v-model.trim="shipmentForm.sender_detail" class="sf-form-grid-wide" placeholder="详细地址" />
            </div>
          </div>
          <div class="sf-form-section">
            <div class="sf-section-head"><strong>收件人信息</strong><button class="sf-text-button" @click="swapShipmentPeople">交换地址</button></div>
            <div class="sf-form-grid">
              <el-input v-model.trim="shipmentForm.receiver_name" placeholder="姓名" />
              <el-input v-model.trim="shipmentForm.receiver_phone" placeholder="手机号" />
              <el-input v-model.trim="shipmentForm.receiver_province" placeholder="省" />
              <el-input v-model.trim="shipmentForm.receiver_city" placeholder="市" />
              <el-input v-model.trim="shipmentForm.receiver_district" placeholder="区" />
              <el-input v-model.trim="shipmentForm.receiver_detail" class="sf-form-grid-wide" placeholder="详细地址" />
            </div>
          </div>
          <div class="sf-form-section">
            <div class="sf-section-head"><strong>包裹与付款</strong><span>默认标快</span></div>
            <div class="sf-form-grid">
              <el-input v-model.number="shipmentForm.weight" type="number" placeholder="重量（kg）" />
              <el-select v-model="shipmentForm.payment_method" placeholder="付款方式">
                <el-option label="寄付现结" :value="1" />
                <el-option label="到付" :value="2" />
              </el-select>
            </div>
          </div>
          <div class="sf-ship-submit-row">
            <span>预计运费将在提交时由后端计算</span>
            <el-button type="primary" :loading="shipmentSaving" @click="submitShipment">提交寄件</el-button>
          </div>
        </section>
      </section>

      <section v-else-if="active === 'track'" class="sf-page sf-track-page">
        <section class="sf-search-row compact">
          <el-input v-model="searchKeyword" class="sf-search-input" placeholder="运单号查询/关键字索引" clearable @keyup.enter="searchTracking" />
          <button class="sf-tool-pill" @click="searchTracking">查询</button>
        </section>

        <div class="sf-track-tabs">
          <button v-for="tab in trackTabs" :key="tab.key" :class="{ on: tab.key === trackFilter }" @click="selectTrackFilter(tab.key)">
            {{ tab.label }} <span>{{ tab.count }}</span>
          </button>
        </div>

        <section class="sf-track-board" :class="{ empty: !filteredTrackRecords.length }">
          <template v-if="filteredTrackRecords.length">
            <article v-for="item in filteredTrackRecords" :key="item.order_no" class="sf-track-card">
              <div class="sf-track-head">
                <strong>{{ item.order_no }}</strong>
                <span>{{ item.status }}</span>
              </div>
              <p>{{ item.sender_name }} → {{ item.receiver_name }}</p>
              <small>{{ item.create_time || "刚刚更新" }}</small>
            </article>
          </template>
          <el-empty v-else description="点击寄快递，开始你的快递之旅吧" />
        </section>

        <aside class="sf-float-note">如需查询 3 个月以上运单请点击这里</aside>
      </section>

      <section v-else-if="active === 'member'" class="sf-page sf-member-page">
        <div class="sf-member-tabs">
          <button v-for="tab in memberTabs" :key="tab.key" :class="{ on: tab.key === memberFilter }" @click="memberFilter = tab.key">
            {{ tab.label }}
          </button>
        </div>

        <section class="sf-benefit-grid">
          <article v-for="item in visibleBenefitCards" :key="item.title" class="sf-benefit-card">
            <span>{{ item.label }}</span>
            <strong>{{ item.title }}</strong>
            <small>{{ item.note }}</small>
          </article>
        </section>

        <section class="sf-deal-grid">
          <article v-for="deal in visibleDeals" :key="deal.name" class="sf-deal-card">
            <div>
              <b>{{ deal.price }}</b>
              <span>{{ deal.name }}</span>
            </div>
            <small>{{ deal.note }}</small>
            <button>立即抢</button>
          </article>
        </section>
      </section>

      <section v-else-if="active === 'service'" class="sf-page sf-service-page">
        <section class="sf-service-hero">
          <strong>生活服务</strong>
          <p>把寄件、查件、地址、支付和会员服务都放在这里。</p>
        </section>

        <section class="sf-grid-panel service-large">
          <button v-for="item in lifeServices" :key="item.label" class="sf-service-tile" @click="handleService(item.target)">
            <b>{{ item.icon }}</b>
            <strong>{{ item.label }}</strong>
            <small>{{ item.desc }}</small>
          </button>
        </section>

        <section class="sf-process-card">
          <strong>业务实现逻辑</strong>
          <ol>
            <li>登录或注册后进入用户端首页。</li>
            <li>选择寄件、查件、地址簿或会员入口。</li>
            <li>有后端接口的功能直接请求接口，没有接口的先保留页面位。</li>
            <li>统一从底部导航切换各个页面。</li>
          </ol>
        </section>
      </section>

      <section v-else class="sf-page sf-mine-page">
        <header class="sf-profile-card">
          <button class="sf-avatar large" @click="openNicknameDialog">{{ profileLabel }}</button>
          <div>
            <p>{{ profile?.phone || "138 **** 8888" }}</p>
            <h2>{{ profile?.nick_name || "李先生" }}</h2>
            <span>{{ profile?.is_enterprise ? "企业账户" : "个人账户" }}</span>
          </div>
          <button class="sf-member-chip">会员中心</button>
        </header>

        <section class="sf-stats-row">
          <article>
            <b>{{ profile?.query_count_today ?? historyCount.query_count_today }}</b>
            <span>今日查询</span>
          </article>
          <article>
            <b>{{ addresses.length }}</b>
            <span>地址簿</span>
          </article>
          <article>
            <b>{{ historyRecords.length }}</b>
            <span>历史快递</span>
          </article>
          <article>
            <b>{{ historyCount.max_query_per_day }}</b>
            <span>每日上限</span>
          </article>
        </section>

        <section class="sf-action-grid">
          <button v-for="item in mineActions" :key="item.label" @click="item.action()">
            <b>{{ item.icon }}</b>
            <span>{{ item.label }}</span>
          </button>
        </section>

        <section class="sf-info-panel">
          <div class="sf-section-head">
            <strong>地址簿</strong>
            <button @click="openAddressDialog()">新增地址</button>
          </div>
          <article v-for="item in addresses" :key="item.id || `${item.receiver_name}-${item.detail_addr}`" class="sf-address-card">
            <div>
              <strong>{{ item.receiver_name }}</strong>
              <span>{{ item.receiver_phone }}</span>
            </div>
            <p>{{ item.province }} {{ item.city }} {{ item.district }} {{ item.detail_addr }}</p>
            <footer>
              <small>{{ item.is_default ? "默认地址" : "普通地址" }}</small>
              <span>
                <a @click="openAddressDialog(item)">编辑</a>
                <a @click="removeAddress(item)">删除</a>
              </span>
            </footer>
          </article>
          <el-empty v-if="!addresses.length" description="还没有地址，先新增一个试试" />
        </section>
      </section>

      <nav class="sf-tabbar">
        <button v-for="tab in tabs" :key="tab.key" :class="{ on: isTabActive(tab.key) }" @click="active = tab.key">
          <span>{{ tab.icon }}</span>
          <strong>{{ tab.label }}</strong>
        </button>
      </nav>
    </section>

    <el-dialog v-model="addressDialogVisible" :title="addressForm.id ? '编辑地址' : '新增地址'" width="420px">
      <el-form label-position="top">
        <el-form-item label="姓名"><el-input v-model.trim="addressForm.receiver_name" /></el-form-item>
        <el-form-item label="电话"><el-input v-model.trim="addressForm.receiver_phone" /></el-form-item>
        <el-form-item label="省 / 市 / 区">
          <div class="sf-dialog-grid">
            <el-input v-model.trim="addressForm.province" placeholder="省" />
            <el-input v-model.trim="addressForm.city" placeholder="市" />
            <el-input v-model.trim="addressForm.district" placeholder="区" />
          </div>
        </el-form-item>
        <el-form-item label="详细地址"><el-input v-model.trim="addressForm.detail_addr" /></el-form-item>
        <el-checkbox v-model="addressForm.is_default">设为默认地址</el-checkbox>
      </el-form>
      <template #footer>
        <el-button @click="addressDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="addressSaving" @click="saveAddress">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="nicknameDialogVisible" title="修改昵称" width="360px">
      <el-input v-model.trim="nicknameDraft" maxlength="20" show-word-limit />
      <template #footer>
        <el-button @click="nicknameDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="nicknameSaving" @click="saveNickname">保存</el-button>
      </template>
    </el-dialog>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { userApi, type ExpressOrderInput } from "@/api/user";
import { useAuthStore } from "@/stores/auth";
import { expressOrders } from "@/data/mock";
import type { Address, HistoryRecord, Profile } from "@/types";

type UserPage = "home" | "ship" | "track" | "member" | "service" | "mine";
type MemberFilter = "benefit" | "points" | "rights";

const auth = useAuthStore();
const active = ref<UserPage>("home");
const cityName = ref("宿迁");
const searchKeyword = ref("");
const trackFilter = ref("all");
const memberFilter = ref<MemberFilter>("benefit");
const profile = ref<Profile | null>(null);
const addresses = ref<Address[]>([]);
const historyRecords = ref<HistoryRecord[]>([]);
const historyCount = reactive({ query_count_today: 0, max_query_per_day: 3 });
const addressDialogVisible = ref(false);
const addressSaving = ref(false);
const nicknameDialogVisible = ref(false);
const nicknameSaving = ref(false);
const nicknameDraft = ref("");
const shipmentSaving = ref(false);
const isPreview = computed(() => auth.session?.accessToken === "local-preview-token");
const currentUserId = computed(() => auth.session?.userId || 1);
const addressForm = reactive<Address>(emptyAddress());
const shipmentForm = reactive<ExpressOrderInput>({
  user_id: currentUserId.value,
  sender_name: "李先生",
  sender_phone: "13800008888",
  sender_province: "江苏省",
  sender_city: "宿迁市",
  sender_district: "宿城区",
  sender_detail: "",
  receiver_name: "",
  receiver_phone: "",
  receiver_province: "",
  receiver_city: "",
  receiver_district: "",
  receiver_detail: "",
  weight: 1,
  channel_code: 1,
  pickup_date: new Date(Date.now() + 86400000).toISOString().slice(0, 10),
  payment_method: 1,
  privacy_protection: true,
  express_type: 1,
});

const previewProfile: Profile = {
  user_id: 1,
  avatar_url: "",
  nick_name: "李先生",
  phone: "13800008888",
  real_name_auth: 1,
  account_status: 1,
  is_enterprise: false,
  query_count_today: 2,
};

const previewAddresses: Address[] = [
  {
    id: 1,
    user_id: 1,
    addr_type: 1,
    receiver_name: "李先生",
    receiver_phone: "13800008888",
    province: "浙江省",
    city: "杭州市",
    district: "滨江区",
    detail_addr: "江南大道 88 号",
    is_default: true,
  },
];

const previewHistory: HistoryRecord[] = expressOrders.map((item, index) => ({
  id: index + 1,
  order_no: item.no,
  express_no: item.no,
  sender_name: item.sender,
  receiver_name: item.receiver,
  status: item.status,
  create_time: `2026-07-29 1${index}:00`,
}));

const tabs = [
  { key: "home", label: "寄快递", icon: "📮" },
  { key: "track", label: "查快递", icon: "🔍" },
  { key: "member", label: "会员福利", icon: "🎁" },
  { key: "service", label: "生活服务", icon: "🏠" },
  { key: "mine", label: "我的", icon: "👤" },
] as const;

const page = computed(() => {
  switch (active.value) {
    case "track":
      return { title: "查快递", subtitle: "运单查询、关注和待支付都放这里", shortcut: "筛选" };
    case "ship":
      return { title: "寄快递", subtitle: "填写寄收件信息，提交后生成待揽收订单", shortcut: "首页" };
    case "member":
      return { title: "会员福利", subtitle: "超值福利、积分商城、会员权益", shortcut: "权益" };
    case "service":
      return { title: "生活服务", subtitle: "寄件、查件、地址、支付一屏完成", shortcut: "寄件" };
    case "mine":
      return { title: "我的", subtitle: "个人资料、地址簿、历史快递和设置", shortcut: "编辑" };
    default:
      return { title: "寄快递", subtitle: "首页、寄件、查件、会员、我的都在这", shortcut: "我的" };
  }
});

const homeServices = [
  { icon: "🚚", label: "寄快递", desc: "上门取件", target: "ship" },
  { icon: "🔍", label: "查快递", desc: "运单追踪", target: "track" },
  { icon: "📕", label: "地址簿", desc: "常用地址", target: "mine" },
  { icon: "🎁", label: "会员福利", desc: "优惠专区", target: "member" },
  { icon: "🏠", label: "生活服务", desc: "更多能力", target: "service" },
  { icon: "🧾", label: "电子面单", desc: "快速下单", target: "service" },
  { icon: "▣", label: "批量寄", desc: "多票一起寄", target: "ship" },
  { icon: "▤", label: "Excel寄", desc: "导入面单", target: "ship" },
  { icon: "◌", label: "服务点自寄", desc: "附近网点", target: "service" },
] as const;

const lifeServices = [
  { icon: "📮", label: "寄快递", desc: "上门取件 / 直营网点", target: "ship" },
  { icon: "🔍", label: "查快递", desc: "单号查询 / 关注", target: "track" },
  { icon: "📕", label: "地址簿", desc: "新增 / 编辑 / 删除", target: "mine" },
  { icon: "🎁", label: "会员福利", desc: "权益 / 券包 / 积分", target: "member" },
  { icon: "📞", label: "客服中心", desc: "在线咨询 / 进度查询", target: "mine" },
  { icon: "⚙️", label: "偏好设置", desc: "消息 / 地址 / 默认支付", target: "mine" },
] as const;

const coupons = [
  { title: "超值券包", value: "3 元", note: "立减券" },
  { title: "折扣券", value: "9.2 折", note: "满额可用" },
  { title: "0.04 元", value: "去抢购", note: "限时专区" },
] as const;

const memberTabs = [
  { key: "benefit", label: "超值福利" },
  { key: "points", label: "积分商城" },
  { key: "rights", label: "会员权益" },
] as const;

const benefitCards = [
  { label: "新速运通卡", title: "首充赠 30 元", note: "新客专享" },
  { label: "顺运年卡", title: "全年 43 张优惠券", note: "持续可用" },
  { label: "家庭账户", title: "亲友互享专属折扣", note: "可共享" },
] as const;

const deals = [
  { name: "超值惊喜套餐", price: "3 元", note: "已抢 7%" },
  { name: "每月限量 88 折套餐", price: "8.8 折", note: "已抢 4%" },
  { name: "限量抢大件券", price: "6 折", note: "已抢 3%" },
] as const;

const memberContent = {
  benefit: {
    cards: benefitCards,
    deals,
  },
  points: {
    cards: [
      { label: "积分商城", title: "积分兑寄件券", note: "满 100 积分可兑换" },
      { label: "积分任务", title: "每日签到得积分", note: "连续签到更划算" },
      { label: "积分明细", title: "本月可用积分 860", note: "积分 12 个月内有效" },
    ],
    deals: [
      { name: "100 积分抵 3 元", price: "3 元", note: "立即兑换" },
      { name: "500 积分抵 10 元", price: "10 元", note: "可兑换 2 次" },
    ],
  },
  rights: {
    cards: [
      { label: "会员权益", title: "专属客服优先响应", note: "服务时间 8:00-22:00" },
      { label: "会员权益", title: "寄件价格更优惠", note: "多种寄件方式可选" },
      { label: "会员权益", title: "生日月专属礼遇", note: "记得领取权益" },
    ],
    deals: [
      { name: "开通年卡享全年优惠", price: "43 张", note: "优惠券礼包" },
      { name: "家庭账户共享", price: "免费", note: "邀请亲友加入" },
    ],
  },
} as const;

const visibleBenefitCards = computed(() => memberContent[memberFilter.value].cards);
const visibleDeals = computed(() => memberContent[memberFilter.value].deals);

const mineActions = [
  { icon: "📦", label: "我的快递", action: () => (active.value = "track") },
  { icon: "📕", label: "地址簿", action: () => (active.value = "mine") },
  { icon: "⚙️", label: "偏好设置", action: () => (active.value = "mine") },
  { icon: "👮", label: "专属快递员", action: () => (active.value = "service") },
  { icon: "🔒", label: "签收码", action: () => ElMessage.info("后续可继续对接签收码功能") },
  { icon: "🎧", label: "客服中心", action: () => ElMessage.info("后续可继续对接客服中心页面") },
  { icon: "🔍", label: "服务查询", action: () => (active.value = "track") },
  { icon: "🧭", label: "地址迁移", action: () => ElMessage.info("后续可继续对接地址迁移功能") },
] as const;

const visibleOrders = computed(() => filteredTrackRecords.value.slice(0, 3));
const profileLabel = computed(() => profile.value?.nick_name?.slice(0, 1) || "李");

const trackTabs = computed(() => {
  const all = historyRecords.value.length;
  const receive = historyRecords.value.filter((item) => item.status.includes("签收")).length;
  const focus = historyRecords.value.filter((item) => item.status.includes("运输") || item.status.includes("关注")).length;
  const pending = historyRecords.value.filter((item) => item.status.includes("待")).length;
  return [
    { key: "all", label: "寄件", count: all },
    { key: "receive", label: "收件", count: receive },
    { key: "focus", label: "关注", count: focus },
    { key: "pending", label: "待支付", count: pending },
  ];
});

const filteredTrackRecords = computed(() => {
  const keyword = searchKeyword.value.trim();
  return historyRecords.value.filter((item) => {
    const matchesKeyword =
      !keyword ||
      [item.order_no, item.express_no, item.sender_name, item.receiver_name, item.status].some((value) =>
        String(value).includes(keyword),
      );
    const matchesTab =
      trackFilter.value === "all" ||
      (trackFilter.value === "receive" && item.status.includes("签收")) ||
      (trackFilter.value === "focus" && (item.status.includes("运输") || item.status.includes("关注"))) ||
      (trackFilter.value === "pending" && item.status.includes("待"));
    return matchesKeyword && matchesTab;
  });
});

onMounted(async () => {
  await Promise.all([loadProfile(), loadAddresses(), loadHistory()]);
});

function emptyAddress(): Address {
  return {
    user_id: currentUserId.value,
    addr_type: 1,
    receiver_name: "",
    receiver_phone: "",
    province: "",
    city: "",
    district: "",
    detail_addr: "",
    is_default: false,
  };
}

function openTracking() {
  active.value = "track";
}

function handleService(target: string) {
  active.value = target as UserPage;
}

function handleShortcut() {
  if (active.value === "ship") active.value = "home";
  else if (active.value === "service") active.value = "ship";
  else if (active.value === "track") searchTracking();
  else if (active.value === "member") memberFilter.value = "rights";
  else active.value = "mine";
}

function isTabActive(key: string) {
  return active.value === key || (key === "home" && active.value === "ship");
}

async function selectTrackFilter(key: string) {
  trackFilter.value = key;
  if (key === "all" || isPreview.value) return;
  const category = key === "receive" ? "received" : key === "focus" ? "followed" : "unpaid";
  try {
    const data = await userApi.listExpressOrdersByCategory(currentUserId.value, category);
    historyRecords.value = (data.orders || []).map(toHistoryRecord);
  } catch {
    ElMessage.warning("该分类暂时无法加载，请稍后再试");
  }
}

async function searchTracking() {
  const keyword = searchKeyword.value.trim();
  if (!keyword || isPreview.value) return;
  try {
    const data = await userApi.searchExpressOrders(currentUserId.value, keyword);
    historyRecords.value = (data.orders || []).map(toHistoryRecord);
    trackFilter.value = "all";
  } catch {
    ElMessage.error("查件失败，请检查后端服务");
  }
}

function toHistoryRecord(item: any): HistoryRecord {
  return {
    id: Number(item.id || 0),
    order_no: item.order_no || item.orderNo || "",
    express_no: item.express_no || item.expressNo || "",
    sender_name: item.sender_name || item.senderName || "",
    receiver_name: item.receiver_name || item.receiverName || "",
    status: normalizeOrderStatus(item.status),
    create_time: normalizeCreatedAt(item.created_at || item.createdAt),
  };
}

function normalizeOrderStatus(status: unknown) {
  const value = String(status || "");
  const map: Record<string, string> = {
    "1": "待揽收",
    "2": "已接单",
    "3": "待上门",
    "4": "已取件",
    "5": "运输中",
    "6": "派送中",
    "7": "已签收",
    "8": "已取消",
    PENDING_PICKUP: "待揽收",
    ACCEPTED: "已接单",
    AWAITING_PICKUP: "待上门",
    PICKED_UP: "已取件",
    IN_TRANSIT: "运输中",
    DELIVERING: "派送中",
    SIGNED: "已签收",
    CANCELLED: "已取消",
  };
  return map[value] || value || "待揽收";
}

function normalizeCreatedAt(value: unknown) {
  if (!value) return "";
  if (typeof value === "string") return value;
  if (typeof value === "object" && "seconds" in value) {
    const seconds = Number((value as { seconds?: number }).seconds || 0);
    if (seconds > 0) {
      const date = new Date(seconds * 1000);
      const pad = (part: number) => String(part).padStart(2, "0");
      return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}`;
    }
  }
  return "";
}

function swapShipmentPeople() {
  const sender = {
    name: shipmentForm.sender_name,
    phone: shipmentForm.sender_phone,
    province: shipmentForm.sender_province,
    city: shipmentForm.sender_city,
    district: shipmentForm.sender_district,
    detail: shipmentForm.sender_detail,
  };
  shipmentForm.sender_name = shipmentForm.receiver_name;
  shipmentForm.sender_phone = shipmentForm.receiver_phone;
  shipmentForm.sender_province = shipmentForm.receiver_province;
  shipmentForm.sender_city = shipmentForm.receiver_city;
  shipmentForm.sender_district = shipmentForm.receiver_district;
  shipmentForm.sender_detail = shipmentForm.receiver_detail;
  shipmentForm.receiver_name = sender.name;
  shipmentForm.receiver_phone = sender.phone;
  shipmentForm.receiver_province = sender.province;
  shipmentForm.receiver_city = sender.city;
  shipmentForm.receiver_district = sender.district;
  shipmentForm.receiver_detail = sender.detail;
}

async function submitShipment() {
  const required = [
    shipmentForm.sender_name,
    shipmentForm.sender_phone,
    shipmentForm.sender_detail,
    shipmentForm.receiver_name,
    shipmentForm.receiver_phone,
    shipmentForm.receiver_detail,
  ];
  if (required.some((value) => !String(value || "").trim()) || shipmentForm.weight <= 0) {
    ElMessage.warning("请填写完整的寄收件信息和包裹重量");
    return;
  }
  shipmentSaving.value = true;
  try {
    if (isPreview.value) {
      ElMessage.success("预览模式：寄件页面已完成，登录后会写入订单");
    } else {
      const result = await userApi.createExpressOrder({
        ...shipmentForm,
        user_id: currentUserId.value,
      });
      ElMessage.success(`寄件提交成功，运单号 ${result.order_no || "待生成"}`);
      historyRecords.value.unshift(toHistoryRecord(result));
    }
    active.value = "track";
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "寄件提交失败");
  } finally {
    shipmentSaving.value = false;
  }
}

function assignAddressForm(address?: Address) {
  Object.assign(addressForm, address ? { ...address } : emptyAddress());
}

async function loadProfile() {
  if (isPreview.value) {
    profile.value = previewProfile;
    historyCount.query_count_today = previewProfile.query_count_today;
    return;
  }
  try {
    profile.value = await userApi.getProfile(currentUserId.value);
    historyCount.query_count_today = profile.value?.query_count_today || 0;
  } catch {
    profile.value = previewProfile;
    historyCount.query_count_today = previewProfile.query_count_today;
  }
}

async function loadAddresses() {
  if (isPreview.value) {
    addresses.value = previewAddresses.map((item) => ({ ...item }));
    return;
  }
  try {
    const data = await userApi.listAddresses(currentUserId.value).catch(() => ({ addresses: [] as Address[] }));
    addresses.value = data.addresses?.length ? data.addresses : previewAddresses.map((item) => ({ ...item }));
  } catch {
    addresses.value = previewAddresses.map((item) => ({ ...item }));
  }
}

async function loadHistory() {
  if (isPreview.value) {
    historyRecords.value = previewHistory;
    historyCount.max_query_per_day = 3;
    return;
  }
  try {
    const countResult = await userApi.getHistoryCount(currentUserId.value).catch(() => historyCount);
    historyCount.query_count_today = countResult.query_count_today ?? historyCount.query_count_today;
    historyCount.max_query_per_day = countResult.max_query_per_day ?? historyCount.max_query_per_day;
    const expressData = await userApi.searchExpressOrders(currentUserId.value, "").catch(() => ({ orders: [], total_count: 0 }));
    if (expressData.orders?.length) {
      historyRecords.value = expressData.orders.map(toHistoryRecord);
      return;
    }
    const data = await userApi.listHistory(currentUserId.value).catch(() => ({ records: [] as HistoryRecord[], total_count: 0 }));
    historyRecords.value = data.records?.length ? data.records : previewHistory;
  } catch {
    historyRecords.value = previewHistory;
  }
}

function openAddressDialog(address?: Address) {
  assignAddressForm(address);
  addressDialogVisible.value = true;
}

function openNicknameDialog() {
  nicknameDraft.value = profile.value?.nick_name || "";
  nicknameDialogVisible.value = true;
}

async function saveNickname() {
  const nextName = nicknameDraft.value.trim();
  if (!nextName) {
    ElMessage.warning("请输入昵称");
    return;
  }
  nicknameSaving.value = true;
  try {
    if (!isPreview.value) {
      await userApi.updateNickname(currentUserId.value, nextName);
      await loadProfile();
    } else if (profile.value) {
      profile.value = { ...profile.value, nick_name: nextName };
    }
    ElMessage.success("昵称已更新");
    nicknameDialogVisible.value = false;
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "昵称更新失败");
  } finally {
    nicknameSaving.value = false;
  }
}

async function saveAddress() {
  if (!addressForm.receiver_name || !addressForm.receiver_phone || !addressForm.detail_addr) {
    ElMessage.warning("请填写姓名、电话和详细地址");
    return;
  }
  addressSaving.value = true;
  try {
    const payload = { ...addressForm, user_id: currentUserId.value };
    if (addressForm.id) {
      if (!isPreview.value) {
        await userApi.updateAddress(payload);
        await loadAddresses();
      } else {
        const index = addresses.value.findIndex((item) => item.id === payload.id);
        if (index >= 0) addresses.value[index] = payload;
      }
      ElMessage.success("地址已更新");
    } else {
      if (!isPreview.value) {
        await userApi.createAddress(payload);
        await loadAddresses();
      } else {
        addresses.value.unshift({ ...payload, id: Date.now() });
      }
      ElMessage.success("地址已新增");
    }
    addressDialogVisible.value = false;
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : "地址保存失败");
  } finally {
    addressSaving.value = false;
  }
}

async function removeAddress(row: Address) {
  if (row.id) {
    await ElMessageBox.confirm("确定删除这个地址吗？", "删除地址", { type: "warning" });
  }
  addresses.value = addresses.value.filter((item) => item.id !== row.id);
  if (!isPreview.value && row.id) {
    await userApi.deleteAddress(currentUserId.value, row.id);
  }
  ElMessage.success("地址已删除");
}
</script>
