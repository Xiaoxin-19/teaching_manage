import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useToast } from '../../composables/useToast';
import { GetDashboardSummary } from '../../api/dashboard';

export function useDashboard() {
  const router = useRouter();
  const toast = useToast();

  const loading = ref(true);
  const currentDate = computed(() => new Date().toLocaleDateString());

  // 核心数据模型
  const summaryData = ref({
    totalStudents: 0,
    monthlyHours: 0,
    totalRemainingHours: 0,
    newStudentsThisMonth: 0,
    monthOverMonth: '+0%',
  });

  const warningData = ref({
    balanceLow: 0,
    balanceNegative: 0
  });

  async function getSummary() {
    try {
      const data = await GetDashboardSummary();
      summaryData.value.totalStudents = data.total_students;
      summaryData.value.monthlyHours = data.monthly_hours;
      summaryData.value.totalRemainingHours = data.total_remaining_hours;
      summaryData.value.monthOverMonth = data.month_over_month;
      summaryData.value.newStudentsThisMonth = data.new_students_this_month;
      warningData.value.balanceLow = data.total_warning;
      warningData.value.balanceNegative = data.total_arrears;
    } catch {
      toast.error('获取概要数据失败');
    }
  }

  const loadDashboardData = async () => {
    loading.value = true;

    try {
      await getSummary();
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