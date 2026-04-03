<template>
  <div class="h-full w-full p-4 bg-gray-50 flex flex-col">
    <!-- 1. 顶部搜索栏 -->
    <div
      class="bg-white p-4 rounded-lg shadow-sm mb-4 flex justify-between items-center"
    >
      <div class="flex items-center gap-3">
        <span class="text-sm font-bold text-gray-600">关键词</span>
        <el-input
          v-model="searchForm.keyword"
          placeholder="请输入用户姓名/手机号"
          class="w-64"
          clearable
          @keyup.enter="search"
          @clear="search"
        />

        <span class="text-sm font-bold text-gray-600 ml-4">注册时间</span>
        <el-date-picker
          v-model="searchForm.dateRange"
          type="daterange"
          range-separator="-"
          start-placeholder="开始"
          end-placeholder="结束"
          value-format="YYYY-MM-DD"
          class="w-60"
        />

        <el-button type="primary" icon="Search" class="ml-2" @click="search"
          >查询</el-button
        >
        <el-button icon="RefreshLeft" @click="handleReset">重置</el-button>
      </div>

      <div class="flex items-center gap-2">
        <el-button type="primary" icon="Plus" @click="handleAdd"
          >新建用户</el-button
        >
        <el-button
          icon="Refresh"
          circle
          @click="loadTableData"
          title="刷新表格"
        />
      </div>
    </div>

    <!-- 2. 表格区域 -->
    <div class="flex-1 bg-white p-4 rounded-lg shadow-sm overflow-hidden">
      <!-- 
        ElTablePro 绑定说明：
        - :data, :total, :loading => 直接读取 Hook 返回的 tableState
        - v-model:current-page, v-model:page-size => 双向绑定分页状态
        - @page-change => 监听分页变化，触发 Hook 的 loadTableData
      -->
      <ElTablePro
        :columns="columns"
        :data="tableState.data"
        :total="tableState.total"
        :loading="tableState.loading"
        v-model:current-page="tableState.currentPage"
        v-model:page-size="tableState.pageSize"
        @page-change="loadTableData"
        stripe
        border
        height="100%"
      >
      </ElTablePro>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import ElTablePro from "@/components/ElTablePro";
import { useElTablePro } from "@/hooks/useElTablePro";
// 假设的 API 引用
import { getUserList as userListApi } from "@/api/user";

// --- 1. 搜索表单 ---
// 使用 ref 包裹对象，便于整体重置和传参时统一通过 .value 访问。
const searchForm = ref({
  keyword: "",
  dateRange: [],
});

// --- 2. 初始化 Hook (核心逻辑) ---
const { tableState, loadTableData, search } = useElTablePro(userListApi, {
  pageSize: 15, // 自定义默认分页大小

  // 【参数构造】每次请求前执行一次，把页面状态转换成后端需要的参数
  params: () => {
    return {
      key: searchForm.value.keyword, // 字段映射：前端叫 keyword，后端叫 key
      // 处理时间范围数组 => 拆解为 startTime/endTime
      startTime: searchForm.value.dateRange?.[0] || "",
      endTime: searchForm.value.dateRange?.[1] || "",
      // 可以在此附加其他固定参数，如 tenantId: 1001
    };
  },

  // 【数据转换】在渲染前预处理数据
  transform: (list) => {
    return list.map((item) => ({
      ...item,
      // 预计算 UI 展示逻辑，避免在 template 中写三元表达式
      _tagType: item.status === 1 ? "success" : "info",
      _statusText: item.status === 1 ? "启用" : "禁用",
    }));
  },
});

const handleReset = () => {
  searchForm.value = {
    keyword: "",
    dateRange: [],
  };
  search();
};

// --- 3. 表格列配置 ---
const columns = [
  { type: "index", label: "#", width: 60, align: "center" },
  { prop: "userName", label: "用户名称", minWidth: 120 },
  { prop: "mobile", label: "手机号码", width: 150 },
  {
    prop: "status",
    label: "状态",
    width: 100,
    slot: "status",
    align: "center",
  },
  { prop: "createTime", label: "注册时间", width: 180 },
  {
    label: "操作",
    width: 150,
    slot: "actions",
    fixed: "right",
    align: "center",
  },
];

// --- 4. 生命周期 ---
onMounted(() => {
  // 显式调用，控制权在开发者手中
  loadTableData();
});
</script>
