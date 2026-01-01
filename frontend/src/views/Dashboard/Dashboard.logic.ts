import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useToast } from '../../composables/useToast';
import { GetDashboardSummaryResponse } from '../../types/response';
import { ResponseWrapper } from '../../types/appModels';
import { Dispatch } from '../../../wailsjs/go/main/App';
// import { Dispatch } from '../../../wailsjs/go/main/App'; // TODO: 对接后端时启用

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

  function getSummary() {
    Dispatch('dashboard_manager:get_summary', '').then((res: string) => {
      const resp = JSON.parse(res) as ResponseWrapper<GetDashboardSummaryResponse>;
      console.log('Dashboard Summary Response:', resp);
      if (resp.code == 200) {
        summaryData.value.totalStudents = resp.data.total_students;
        summaryData.value.monthlyHours = resp.data.monthly_hours;
        summaryData.value.totalRemainingHours = resp.data.total_remaining_hours;
        summaryData.value.monthOverMonth = resp.data.month_over_month;
        summaryData.value.newStudentsThisMonth = resp.data.new_students_this_month;
        warningData.value.balanceLow = resp.data.total_warning;
        warningData.value.balanceNegative = resp.data.total_arrears;
      }
    }).catch(() => {
      toast.error('获取概要数据失败');
    })
  }

  const loadDashboardData = async () => {
    loading.value = true;

    try {
      // 模拟请求
      await new Promise(resolve => setTimeout(resolve, 800));
      getSummary();
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