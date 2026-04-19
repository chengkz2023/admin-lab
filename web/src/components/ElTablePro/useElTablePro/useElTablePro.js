import { ref } from "vue";

/**
 * ElTablePro Hook
 *
 * @param {Function} service - API 请求函数
 *   - 必须返回 Promise
 *   - 响应结构需符合: { code: 0, data: { list: [], totalCount: 0 } }
 *
 * @param {Object} options - 配置对象
 * @param {number} [options.pageSize=20] - 默认每页显示条数
 *
 * @param {Function} [options.params] - 参数构造钩子
 *   - 作用：构造发送给后端的除分页外的其他参数（如搜索表单）。
 *   - 签名：() => Object
 *   - 注意：不需要入参，因为分页参数由 Hook 内部自动管理。
 *
 * @param {Function} [options.transform] - 数据转换钩子
 *   - 作用：在渲染前对后端返回的 list 数据进行格式化或增强。
 *   - 签名：(list: Array) => Array
 *   - 返回：处理后的新数组
 *
 * @returns {Object}
 *   - tableState: { data, total, loading, currentPage, pageSize } (Ref对象)
 *   - loadTableData: Function (触发加载)
 *   - search: Function (重置页码为1并加载)
 */
export function useElTablePro(service, options = {}) {
  // 解构配置项，设置默认值
  const { pageSize = 20, params, transform = (d) => d } = options;

  // 核心状态 (使用单一 Ref 对象管理，方便解构且不易丢失响应性)
  const tableState = ref({
    data: [], // 表格数据
    total: 0, // 总条数
    loading: false, // 加载 loading
    currentPage: 1, // 当前页码
    pageSize: pageSize, // 每页条数
  });

  /**
   * 内部辅助：重置表格数据
   * 用于请求失败或发生异常时兜底，防止展示脏数据
   */
  const resetData = () => {
    tableState.value.data = [];
    tableState.value.total = 0;
  };

  /**
   * 核心动作：加载数据
   * 触发场景：挂载时、翻页时、更改每页大小时、刷新时
   */
  const loadTableData = async () => {
    tableState.value.loading = true;

    try {
      // 1. 【项目约定】构造默认分页参数
      // 如果后端规范变更 (如 pageNo/pageSize)，请在此处修改，全项目生效
      const defaultPagination = {
        page: tableState.value.currentPage,
        limit: tableState.value.pageSize,
      };

      // 2. 获取用户自定义参数
      // 调用 params 钩子获取最新的表单数据
      const customParams =
        typeof params === "function" ? params() : params || {};

      // 3. 智能合并参数
      // 默认分页在前，用户参数在后。
      // 这允许用户在特殊情况下通过 params 返回同名 key 来覆盖默认分页逻辑。
      const finalPayload = { ...defaultPagination, ...customParams };

      // 4. 发起 API 请求
      const res = await service(finalPayload);

      // 5. 结果校验 (卫语句风格)
      // 根据项目实际响应结构调整判断逻辑
      if (res?.code !== 0 || !res?.data) {
        // 可以根据需要在此处添加 ElMessage.warning(res.msg)
        resetData();
        return;
      }

      // 6. 数据转换与赋值
      // 优先使用 transform 钩子处理数据，保持 template 纯净
      tableState.value.data = transform(res.data.list || []);
      tableState.value.total = res.data.totalCount || 0;
    } catch (error) {
      console.error("ElTablePro Load Error:", error);
      resetData();
    } finally {
      tableState.value.loading = false;
    }
  };

  /**
   * 辅助动作：搜索
   * 场景：点击查询按钮、重置按钮
   * 逻辑：将页码重置为 1，然后触发加载
   */
  const search = () => {
    tableState.value.currentPage = 1;
    loadTableData();
  };

  return {
    tableState, // 响应式状态对象
    loadTableData, // 核心加载方法
    search, // 搜索方法
  };
}
