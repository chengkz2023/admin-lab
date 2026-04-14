<template>
  <DefaultCell
      :value="value"
      :row="row"
      :column="enhancedColumn"
  />
</template>

<script setup>
import { computed } from 'vue';
import TypeCell from "@/components/ElTablePro/renderers/typeCell.vue";
import DefaultCell from "@/components/ElTablePro/renderers/defaultCell.vue";


const props = defineProps({
  value: { type: [ String, Number, Boolean, Object, Array ], default: null },
  row: { type: Object, default: () => ({}) },
  column: { type: Object, default: () => ({}) },
});

const formatter = (row,col,cell) => {
  const msgTypeMap = {
    0x0b: "ISMS流量结果上报策略",
    0x82: "IDC机房与IP地址段对应信息下发",
    0x85: "机房/ISMS管理策略绑定信息下发",
    0xc1: "ISMS模块心跳信息",
    0xc2: "ISMS模块策略同步请求",
    0xc3: "ISMS模块策略同步响应",
    0xcd: "策略下发Ack",
    0xd1: "EU设备通用信息下发",
    0xd2: "EU设备状态查询请求",
    0xd3: "EU设备状态查询回应(设备静态信息)",
    0xd4: "EU设备状态查询回应(设备动态信息)",
    0x10: "IDC/ISP信息安全管理",
    0xe6: "信息安全代码发布",
    0xe1: "流量采集管理指令",
    0xe2: "恶意报文监测指令",
    0xe5: "流量采集执行通知",
    0xe7: "网络安全设备状态管理指令",
    0xe8: "网络安全代码发布",
    0xe3: "恶意文件监测指令",
    0xe9: "网络安全设备状态管理指令（高级网安）",
    0xea: "网络安全代码发布（高级网安）",
    0xe4: "数据安全巡查指令",
    0xeb: "数据安全关键词指令",
    0xec: "数据安全代码指令发布",
  }

  return msgTypeMap[cell] ?? `未知类型（${cell}）`
}

// 将预设配置与用户传入的 config 进行合并
const enhancedColumn = computed(() => {
  return {
    ...props.column,
    formatter: formatter
  }
});
</script>
