import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useToast } from '../../composables/useToast';
// import { Dispatch } from '../../../wailsjs/go/main/App'; // TODO: 对接后端时启用

export function useDashboard() {
  const router = useRouter();
  const toast = useToast();

  const loading = ref(true);
  const currentDate = computed(() => new Date().toLocaleDateString());

  // 核心数据模型
  const summaryData = ref({
    totalStudents: 0,
    totalTeachers: 0,
    monthlyHours: 0,
    totalRemainingHours: 0,
    monthOverMonth: '+0%',
    todayLessons: 0
  });

  const warningData = ref({
    balanceLow: 0,
    balanceNegative: 0
  });

  const loadDashboardData = async () => {
    loading.value = true;

    try {
      // 模拟请求
      await new Promise(resolve => setTimeout(resolve, 800));

      // TODO: 替换为实际后端请求
      // const res = await Dispatch('dashboard:get_summary', '');

      summaryData.value = {
        totalStudents: 128,
        totalTeachers: 12,
        totalRemainingHours: 1850,
        monthlyHours: 345,
        monthOverMonth: '+12%',
        todayLessons: 8
      };

      warningData.value = {
        balanceLow: 15,
        balanceNegative: 3
      };
    } catch (e) {
      toast.error('加载数据失败');
    } finally {
      loading.value = false;
    }
  };

  const navigateTo = (routeName: string) => {
    router.push({ name: routeName });
  };

  onMounted(() => {
    loadDashboardData();
  });

  return {
    loading,
    summaryData,
    warningData,
    currentDate,
    loadDashboardData,
    navigateTo
  };
}