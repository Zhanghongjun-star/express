export const expressOrders = [
  { no: "SF2026073001", sender: "王明", receiver: "李娜", route: "杭州 → 上海", status: "待揽收", pay: "寄付现结" },
  { no: "SF2026073002", sender: "赵强", receiver: "陈晨", route: "深圳 → 北京", status: "运输中", pay: "到付" },
  { no: "SF2026073003", sender: "周婷", receiver: "吴越", route: "广州 → 成都", status: "已签收", pay: "月结" },
  { no: "SF2026073004", sender: "宋佳", receiver: "韩冬", route: "南京 → 苏州", status: "待支付", pay: "寄付现结" },
];

export const courierTasks = [
  { no: "SF8302", name: "沈女士", phone: "138****9021", address: "滨江区江南大道 88 号", tags: ["需签收", "易碎品"], time: "10:30 前" },
  { no: "SF8303", name: "何先生", phone: "156****3110", address: "西湖区文三路 19 号", tags: ["货到付款"], time: "12:00 前" },
  { no: "SF8304", name: "林女士", phone: "139****6608", address: "上城区延安路 302 号", tags: ["大件"], time: "14:00 前" },
];

export const merchantJobs = [
  { title: "批量入库扫码", count: 126, note: "到港货车 3 批次待确认" },
  { title: "客户取件核验", count: 38, note: "含 6 件滞留快件" },
  { title: "上门取件调度", count: 14, note: "待分配快递员" },
  { title: "售后异常工单", count: 9, note: "破损、拒收、改地址" },
];

export const adminRows = [
  { no: "SF2026073001", time: "2026-07-29 09:12", status: "已揽收", pay: "寄付", from: "杭州", to: "上海" },
  { no: "SF2026073002", time: "2026-07-29 10:05", status: "运输中", pay: "到付", from: "深圳", to: "北京" },
  { no: "SF2026073003", time: "2026-07-29 11:48", status: "派送中", pay: "月结", from: "广州", to: "成都" },
  { no: "SF2026073004", time: "2026-07-29 13:20", status: "转寄退回", pay: "寄付", from: "南京", to: "苏州" },
];
