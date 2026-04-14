import * as echarts from 'echarts'


/**
 * 生成ECharts图表的自动渐变色样式
 * @param {number} count - 需要生成的颜色样式数量
 * @param {Object} options - 配置选项
 * @param {number} options.gradientTopAlpha - 渐变顶部透明度 (0-1)
 * @param {number} options.gradientBottomAlpha - 渐变底部透明度 (0-1)
 * @param {number} options.shadowAlpha - 阴影透明度 (0-1)
 * @param {number} options.shadowBlur - 阴影模糊程度
 * @param {number} options.lineAlpha - 线条透明度 (0-1)
 * @param {boolean} options.useECharts - 是否使用ECharts生成渐变 (默认true)
 * @param {string} options.palette - 调色板类型: 'default' | 'modern' | 'warm' | 'cool' | 'business' | 'vibrant' | 'natural' | 'pastel' (默认 'default')
 * @param {boolean} options.optimizeForMixed - 是否针对混合图表类型优化 (默认true)
 * @returns {Array<{
 *   lineColor: string,
 *   itemColor: string,
 *   lineStyle: { color: string },
 *   itemStyle: { color: string },
 *   areaStyle: {
 *     color: echarts.graphic.LinearGradient | string,
 *     shadowColor: string,
 *     shadowBlur: number
 *   },
 *   barStyle: {
 *     color: echarts.graphic.LinearGradient | string,
 *     borderColor: string,
 *     shadowColor: string,
 *     shadowBlur: number
 *   }
 * }>} 样式对象数组, 包含lineColor, itemColor, lineStyle, itemStyle, areaStyle, barStyle;使用示例:
 * ```
 *   // 生成颜色配置
 *   const colors = generateColorStyles(chartData.length, {palette: 'vibrant'});
 *   // 生成series配置
 *   const series = chartData?.map((item, index) => ({
 *         name: item.name,
 *         type: 'line',
 *         smooth: true,
 *         data: item.data.map(d => d.value),
 *         lineStyle: colors[index].lineStyle,  // 线条样式
 *         itemStyle: colors[index].itemStyle,  // 图表项样式（标记点等）
 *         areaStyle: colors[index].areaStyle,  // 区域填充样式
 *       })
 *   )
 * ```
 *
 */
export const generateColorStyles = (count = 1, options = {}) => {
    // 默认配置
    const config = {
        gradientTopAlpha: 0.3,  // 渐变顶部透明度 (0-1)
        gradientBottomAlpha: 0,   // 渐变底部透明度 (0-1)
        shadowAlpha: 0.8,   // 阴影透明度 (0-1)
        shadowBlur: 16,  // 阴影模糊程度
        lineAlpha: 1,   // 线条透明度 (0-1)
        useECharts: true, // 是否使用ECharts生成渐变 (默认true)
        palette: 'default',
        optimizeForMixed: true, // 是否针对混合图表类型优化
        ...options
    };

    // 预定义调色板集合
    const PALETTES = {
        // === 基础主题 ===
        // 默认主题 (源于您的图表)
        default: [
            '#3b82f6', // 商务蓝
            '#10b981', // 翡翠绿
            '#8b5cf6', // 活力紫
            '#f97316', // 警告橙
            '#1f2937', // 石墨黑
            '#9ca3af', // 中性灰
        ],

        // ECharts经典主题
        classic: ['#5470c6', '#91cc75', '#fac858', '#ee6666', '#73c0de', '#3ba272', '#fc8452', '#9a60b4', '#ea7ccc'],

        // === 现代主题系列 ===
        // 活力彩虹
        vibrant: ['#ff6b6b', '#f7931e', '#ffd166', '#06d6a0', '#118ab2', '#073b4c', '#ef476f', '#fca311', '#9d4edd'],

        // 现代石墨
        // graphite: ['#262626', '#525252', '#737373', '#a3a3a3', '#d4d4d4', '#f5f5f5'], // 已移除

        // === 专业商务系列 ===
        // 企业级灰蓝
        corporate: ['#1e40af', '#374151', '#6b7280', '#9ca3af', '#d1d5db', '#f3f4f6', '#1f2937', '#111827'],

        // === 科技主题系列 ===
        // 科技蓝调
        techBlue: ['#0a3d91', '#0062d1', '#008cff', '#4dafff', '#8cd1ff', '#c2e5ff'],

        // === 自然主题系列 ===
        // 森林绿调
        forest: ['#064e3b', '#065f46', '#047857', '#059669', '#10b981', '#34d399', '#6ee7b7', '#a7f3d0'],

        // === 特殊效果系列 ===
        // 黑白灰度
        monochrome: ['#000000', '#1f2937', '#374151', '#6b7280', '#9ca3af', '#d1d5db', '#f3f4f6', '#ffffff'],

        // 高对比度
        highContrast: ['#000000', '#ffffff', '#dc2626', '#16a34a', '#2563eb', '#ca8a04', '#9333ea', '#c2410c'],

        // 暗黑模式
        darkMode: ['#f8fafc', '#e2e8f0', '#cbd5e1', '#94a3b8', '#64748b', '#475569', '#334155', '#1e293b'],
    };

    // 工具函数：十六进制转RGB
    const hexToRgb = (hex) => {
        const match = hex.match(/^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i);
        return match ? match.slice(1).map(x => parseInt(x, 16)) : [0, 0, 0];
    };

    // 获取当前调色板
    const currentPalette = PALETTES[config.palette] || PALETTES.default;

    // 获取指定索引的RGB颜色
    const getRgbColor = (index) => {
        if (index < currentPalette.length) {
            return hexToRgb(currentPalette[index]);
        }

        // 超出调色板范围时，循环使用并调整亮度
        const baseIndex = index % currentPalette.length;
        const variation = Math.floor(index / currentPalette.length);
        const factor = 1 + (variation * 0.2 * (index % 2 === 0 ? 1 : -1));

        return hexToRgb(currentPalette[baseIndex])
            .map(c => Math.max(0, Math.min(255, Math.round(c * factor))));
    };

    // 创建柱状图渐变色
    const createBarGradient = (r, g, b) => {
        const { useECharts } = config;
        const barTopAlpha = config.optimizeForMixed ? 0.9 : config.gradientTopAlpha;
        const barBottomAlpha = config.optimizeForMixed ? 0.7 : config.gradientBottomAlpha;

        return useECharts
            ? new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                { offset: 0, color: `rgba(${r}, ${g}, ${b}, ${barTopAlpha})` },
                { offset: 1, color: `rgba(${r}, ${g}, ${b}, ${barBottomAlpha})` }
            ], false)
            : `linear-gradient(to bottom, rgba(${r}, ${g}, ${b}, ${barTopAlpha}) 0%, rgba(${r}, ${g}, ${b}, ${barBottomAlpha}) 100%)`;
    };

    // 创建线图渐变色
    const createLineGradient = (r, g, b) => {
        const { gradientTopAlpha, gradientBottomAlpha, useECharts } = config;

        return useECharts
            ? new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                { offset: 0, color: `rgba(${r}, ${g}, ${b}, ${gradientTopAlpha})` },
                { offset: 0.9, color: `rgba(${r}, ${g}, ${b}, ${gradientBottomAlpha})` }
            ], false)
            : `linear-gradient(to bottom, rgba(${r}, ${g}, ${b}, ${gradientTopAlpha}) 0%, rgba(${r}, ${g}, ${b}, ${gradientBottomAlpha}) 90%)`;
    };

    // 创建单个颜色样式
    const createColorStyle = (rgb) => {
        const [r, g, b] = rgb;
        const { lineAlpha, shadowAlpha, shadowBlur } = config;
        const baseColor = `rgba(${r}, ${g}, ${b}, ${lineAlpha})`;
        const enhancedShadowBlur = config.optimizeForMixed ? shadowBlur * 0.5 : shadowBlur;

        return {
            lineColor: baseColor,
            itemColor: baseColor,
            lineStyle: {
                color: baseColor
            },
            itemStyle: {
                color: baseColor
            },
            areaStyle: {
                color: createLineGradient(r, g, b),
                shadowColor: `rgba(${r}, ${g}, ${b}, ${shadowAlpha})`,
                shadowBlur: enhancedShadowBlur
            },
            // 新增：专门的柱状图样式
            barStyle: {
                color: createBarGradient(r, g, b),
                borderColor: `rgba(${r}, ${g}, ${b}, 0.8)`,
                shadowColor: `rgba(${r}, ${g}, ${b}, 0.3)`,
                shadowBlur: 8
            }
        };
    };

    // 生成所有颜色样式
    return Array.from({ length: count }, (_, index) =>
        createColorStyle(getRgbColor(index))
    );
};
