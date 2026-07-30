<template>
  <PortalShell
    title="用户端"
    subtitle="个人寄件与查件"
    eyebrow="User service"
    login-path="/user/login"
    :active="active"
    :menu="menu"
    @select="active = $event"
  >
    <section v-if="active === 'send'" class="mobile-stage">
      <div class="phone-frame">
        <div class="phone-header">
          <span>顺丰速运</span>
          <el-tag size="small" type="danger">极速下单</el-tag>
        </div>
        <el-tabs v-model="sendType">
          <el-tab-pane label="寄快递" name="express" />
          <el-tab-pane label="发物流" name="freight" />
        </el-tabs>
        <el-form label-position="top" class="compact-form">
          <el-form-item label="寄件人">
            <el-input v-model="order.sender" placeholder="选择或填写寄件人地址" />
          </el-form-item>
          <el-form-item label="收件人">
            <el-input v-model="order.receiver" placeholder="选择或填写收件人地址" />
          </el-form-item>
          <el-form-item label="寄件方式">
            <el-segmented v-model="order.method" :options="['上门取件', '服务点自寄']" />
          </el-form-item>
          <el-form-item v-if="order.method === '服务点自寄'" label="服务点">
            <el-select v-model="order.station">
              <el-option label="丰巢柜 · 小格口" value="丰巢柜" />
              <el-option label="顺丰网点 · 可称重" value="顺丰网点" />
              <el-option label="合作商家店 · 代收" value="合作商家店" />
            </el-select>
          </el-form-item>
          <div class="form-grid two">
            <el-form-item label="预约时间">
              <el-time-select v-model="order.time" start="09:00" step="01:00" end="20:00" />
            </el-form-item>
            <el-form-item label="付款方式">
              <el-select v-model="order.pay">
                <el-option label="寄付现结" value="寄付现结" />
                <el-option label="到付" value="到付" />
                <el-option label="寄付月结" value="寄付月结" />
              </el-select>
            </el-form-item>
          </div>
          <el-form-item label="物品信息">
            <el-input v-model="order.item" placeholder="物品名称、重量、件数" />
          </el-form-item>
          <div class="option-row">
            <el-check-tag checked>优惠券可用</el-check-tag>
            <el-check-tag>隐私号码</el-check-tag>
            <el-check-tag>保价</el-check-tag>
          </div>
          <el-button type="primary" class="full-button">提交订单</el-button>
        </el-form>
      </div>
      <div class="logic-panel">
        <h2>业务实现逻辑</h2>
        <el-timeline>
          <el-timeline-item timestamp="1">校验登录态，未登录跳转登录页并保留返回路径。</el-timeline-item>
          <el-timeline-item timestamp="2">拉取地址簿，用户选择寄件人和收件人。</el-timeline-item>
          <el-timeline-item timestamp="3">根据寄件方式判断上门时间、网点和格口是否可用。</el-timeline-item>
          <el-timeline-item timestamp="4">确认付款方式、优惠券和加密服务后生成订单。</el-timeline-item>
          <el-timeline-item timestamp="5">订单进入待揽收，同步查快递与运营统计。</el-timeline-item>
        </el-timeline>
      </div>
    </section>

    <section v-if="active === 'track'" class="content-grid">
      <el-card shadow="never" class="span-2">
        <template #header>
          <div class="card-header">
            <strong>查快递</strong>
            <el-input v-model="keyword" placeholder="运单号 / 关键字 / 标签" clearable />
          </div>
        </template>
        <el-tabs>
          <el-tab-pane label="寄件" />
          <el-tab-pane label="收件" />
          <el-tab-pane label="关注" />
          <el-tab-pane label="待支付" />
        </el-tabs>
        <el-table :data="filteredOrders">
          <el-table-column prop="no" label="订单号" />
          <el-table-column prop="receiver" label="收件人" />
          <el-table-column prop="sender" label="寄件人" />
          <el-table-column prop="route" label="线路" />
          <el-table-column prop="status" label="状态">
            <template #default="{ row }"><el-tag>{{ row.status }}</el-tag></template>
          </el-table-column>
          <el-table-column label="操作" width="220">
            <template #default>
              <el-button link type="primary">详情</el-button>
              <el-button link>分享</el-button>
              <el-button link>关注</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </section>

    <section v-if="active === 'mine'" class="content-grid">
      <el-card shadow="never">
        <template #header><strong>个人信息</strong></template>
        <div class="profile-card">
          <el-avatar :size="64">{{ profile?.nick_name?.slice(0, 1) || "顺" }}</el-avatar>
          <div>
            <h3>{{ profile?.nick_name || "本地预览用户" }}</h3>
            <p>{{ profile?.phone || "登录后展示手机号" }}</p>
            <el-tag>{{ profile?.is_enterprise ? "企业月结" : "个人用户" }}</el-tag>
          </div>
        </div>
      </el-card>
      <el-card shadow="never">
        <template #header><strong>历史快递</strong></template>
        <div class="metric-row">
          <div><b>{{ historyCount.query_count_today }}</b><span>今日查询</span></div>
          <div><b>{{ historyCount.max_query_per_day }}</b><span>每日上限</span></div>
        </div>
        <el-button class="full-button" @click="loadHistory">查询 90 天记录</el-button>
      </el-card>
      <el-card shadow="never" class="span-2">
        <template #header>
          <div class="card-header">
            <strong>地址簿</strong>
            <el-button type="primary" :icon="Plus" @click="addAddress">新增地址</el-button>
          </div>
        </template>
        <el-table :data="addresses">
          <el-table-column prop="receiver_name" label="姓名" />
          <el-table-column prop="receiver_phone" label="手机号" />
          <el-table-column label="地区">
            <template #default="{ row }">{{ row.province }} {{ row.city }} {{ row.district }}</template>
          </el-table-column>
          <el-table-column prop="detail_addr" label="详细地址" />
          <el-table-column label="操作" width="180">
            <template #default="{ row }">
              <el-button link type="primary" @click="editAddress(row)">编辑</el-button>
              <el-button link type="danger" @click="removeAddress(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </section>
  </PortalShell>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { ElMessage } from "element-plus";
import { Box, HomeFilled, Search, User, Plus } from "@element-plus/icons-vue";
import PortalShell from "@/components/PortalShell.vue";
import { userApi } from "@/api/user";
import { useAuthStore } from "@/stores/auth";
import { expressOrders } from "@/data/mock";
import type { Address, Profile } from "@/types";

const auth = useAuthStore();
const active = ref("send");
const sendType = ref("express");
const keyword = ref("");
const profile = ref<Profile | null>(null);
const addresses = ref<Address[]>([]);
const historyCount = reactive({ query_count_today: 0, max_query_per_day: 3 });
const menu = [
  { key: "send", label: "寄快递", icon: Box },
  { key: "track", label: "查快递", icon: Search },
  { key: "mine", label: "我的", icon: User },
  { key: "service", label: "服务查询", icon: HomeFilled },
];
const order = reactive({
  sender: "",
  receiver: "",
  method: "上门取件",
  station: "丰巢柜",
  time: "10:00",
  pay: "寄付现结",
  item: "",
});
const filteredOrders = computed(() => {
  if (!keyword.value) return expressOrders;
  return expressOrders.filter((item) => Object.values(item).some((value) => String(value).includes(keyword.value)));
});

onMounted(async () => {
  await Promise.all([loadProfile(), loadAddresses(), loadCount()]);
});

async function loadProfile() {
  if (!auth.session?.userId || auth.session.accessToken === "local-preview-token") return;
  profile.value = await userApi.getProfile(auth.session.userId).catch(() => null);
}

async function loadAddresses() {
  if (!auth.session?.userId || auth.session.accessToken === "local-preview-token") {
    addresses.value = [
      { id: 1, user_id: 1, addr_type: 1, receiver_name: "罗明月", receiver_phone: "13800000000", province: "浙江省", city: "杭州市", district: "滨江区", detail_addr: "江南大道 88 号" },
    ];
    return;
  }
  const data = await userApi.listAddresses(auth.session.userId).catch(() => ({ addresses: [] }));
  addresses.value = data.addresses || [];
}

async function loadCount() {
  if (!auth.session?.userId || auth.session.accessToken === "local-preview-token") return;
  Object.assign(historyCount, await userApi.getHistoryCount(auth.session.userId).catch(() => historyCount));
}

async function loadHistory() {
  if (!auth.session?.userId || auth.session.accessToken === "local-preview-token") return ElMessage.info("预览模式展示静态记录");
  const data = await userApi.listHistory(auth.session.userId);
  ElMessage.success(`查到 ${data.total_count} 条历史快递`);
}

function addAddress() {
  addresses.value.unshift({
    id: Date.now(),
    user_id: auth.session?.userId || 1,
    addr_type: 1,
    receiver_name: "新收件人",
    receiver_phone: "13900000000",
    province: "浙江省",
    city: "杭州市",
    district: "上城区",
    detail_addr: "请输入详细地址",
  });
}

async function editAddress(row: Address) {
  row.detail_addr = `${row.detail_addr} 已编辑`;
  if (auth.session?.accessToken !== "local-preview-token") await userApi.updateAddress(row);
  ElMessage.success("地址已更新");
}

async function removeAddress(row: Address) {
  addresses.value = addresses.value.filter((item) => item.id !== row.id);
  if (auth.session?.userId && row.id && auth.session.accessToken !== "local-preview-token") {
    await userApi.deleteAddress(auth.session.userId, row.id);
  }
  ElMessage.success("地址已删除");
}
</script>
