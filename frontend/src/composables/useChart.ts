import { onMounted, onUnmounted, nextTick, watch, Ref } from 'vue';
import * as echarts from 'echarts';
import { useTheme } from 'vuetify';

/**
 * ECharts 通用 Hook
 * @param chartRef 图表容器的 DOM 引用
 * @param optionGenerator 生成图表配置的函数，接收 isDark 参数
 */
export function useChart(
  chartRef: Ref<HTMLElement | null>,
  optionGenerator: (isDark: boolean) => any
) {
  const theme = useTheme();
  let chart: echarts.ECharts | null = null;
  let resizeObserver: ResizeObserver | null = null;

  const init = () => {
    if (!chartRef.value) return;

    // 销毁旧实例（主要为了切换主题）
    if (chart) {
      chart.dispose();
    }

    const isDark = theme.global.current.value.dark;
    // 使用 ECharts 自带的 'dark' 主题或默认主题
    chart = echarts.init(chartRef.value, isDark ? 'dark' : undefined, {
      renderer: 'canvas'
    });

    const option = {
      ...optionGenerator(isDark),
      // 确保背景透明，适配 Vuetify 卡片背景
      backgroundColor: 'transparent'
    };
    chart.setOption(option);
  };

  const resize = () => {
    chart?.resize();
  };

  // 监听 Vuetify 主题变化，自动重绘
  watch(() => theme.global.name.value, () => {
    init();
  });

  onMounted(() => {
    nextTick(() => {
      init();

      if (chartRef.value) {
        resizeObserver = new ResizeObserver(() => {
          resize();
        });
        resizeObserver.observe(chartRef.value);
      }
    });
  });

  onUnmounted(() => {
    if (resizeObserver) {
      resizeObserver.disconnect();
    }
    chart?.dispose();
  });

  return {
    getInstance: () => chart,
    setOption: (opt: any) => chart?.setOption(opt),
    refresh: init
  };
}